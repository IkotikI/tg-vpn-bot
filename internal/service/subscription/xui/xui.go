package xui_service

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"vpn-tg-bot/internal/storage"
	x_ui "vpn-tg-bot/pkg/clients/x-ui"
	"vpn-tg-bot/pkg/clients/x-ui/model"
	"vpn-tg-bot/pkg/e"

	"github.com/google/uuid"
)

var ErrZeroUserID = errors.New("user id is 0")
var ErrZeroServerID = errors.New("server id is 0")
var ErrInvalidSubscriptionStatus = errors.New("invalid subscription status")
var ErrClientNotFound = errors.New("client not found")

type ClientID = x_ui.ClientID

type XUIService struct {
	inboundID int
	retries   int
	TokenKey  string

	storage   storage.Storage
	authStore storage.ServerAuthorizations

	clients map[storage.ServerID]struct {
		updated_at time.Time
		xui        *x_ui.XUIClient
	}

	mutex sync.Mutex
}

func NewXUIService(tokenKey string, store storage.Storage, authStore storage.ServerAuthorizations) *XUIService {

	xui := &XUIService{
		inboundID: 1,
		retries:   3,
		TokenKey:  tokenKey,
		storage:   store,
		authStore: authStore,
		clients: make(map[storage.ServerID]struct {
			updated_at time.Time
			xui        *x_ui.XUIClient
		}, 1),
	}
	go xui.watchAndClearUnusedClients(30*time.Second, 30*time.Second)
	return xui
}

// func (s *XUIService) Onlines(ctx context.Context, serverID storage.ServerID) (emails *[]string, err error) {

// }

func (s *XUIService) SubscriptionLink(ctx context.Context, serverID storage.ServerID, userID storage.UserID) (link string, err error) {
	if userID == 0 {
		return "", ErrZeroUserID
	}
	if serverID == 0 {
		return "", ErrZeroServerID
	}

	xui, err := s.xuiClientInstance(ctx, serverID)
	if err != nil {
		return "", err
	}

	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	subID := s.ClientSubIDByUser(user)

	return xui.GetSubBySubID(ctx, subID)
}

// Update Subscription with XUI API and in the Storage.
func (s *XUIService) UpdateSubscription(ctx context.Context, sub *storage.Subscription) (err error) {
	defer func() { e.WrapIfErr("can't update subscription", err) }()

	if sub.UserID == 0 {
		return ErrZeroUserID
	}

	if sub.ServerID == 0 {
		return ErrZeroServerID
	}

	oldSub, err := s.storage.GetSubscriptionByIDs(ctx, sub.UserID, sub.ServerID)
	if err != nil && err != storage.ErrNoSuchSubscription {
		return err
	}

	if sub.SubscriptionExpiredAt.IsZero() {
		sub.SubscriptionExpiredAt = oldSub.SubscriptionExpiredAt
	}

	var enable bool = true
	// If disabled, then disabled. Otherwise, if not provided, set status
	// based of expired_at time.
	switch {
	case oldSub.SubscriptionStatus == storage.SubscriptionStatusDisabled:
		sub.SubscriptionStatus = storage.SubscriptionStatusDisabled
		enable = false
	case sub.SubscriptionStatus == "":
		{
			if sub.SubscriptionExpiredAt.Unix() > time.Now().Unix() {
				sub.SubscriptionStatus = storage.SubscriptionStatusActive
			} else {
				sub.SubscriptionStatus = storage.SubscriptionStatusExpired
			}
		}
	}

	if !sub.IsCorrectStatus() {
		return ErrInvalidSubscriptionStatus
	}

	user, err := s.storage.GetUserByID(ctx, sub.UserID)
	if err != nil {
		return err
	}

	xui, err := s.xuiClientInstance(ctx, sub.ServerID)
	if err != nil {
		return e.Wrap("can't get xui client instance", err)
	}

	client := &model.Client{
		ID:         s.ClientIDByUser(user).String(),
		Email:      s.ClientEmailByUser(user),
		ExpiryTime: sub.SubscriptionExpiredAt.UnixMilli(),
		Enable:     enable,
		SubID:      s.ClientSubIDByUser(user),
	}

	for i := 0; i < s.retries; i++ {
		ctx2, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
		err = xui.UpdateClient(ctx2, s.inboundID, client)
		cancel()
		if err != nil {
			if err == x_ui.ErrRecordNotFound {
				log.Printf("[INFO] can't find client record on remote, try to add a new one")
				ctx2, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
				err = xui.AddClient(ctx2, s.inboundID, client)
				cancel()
				if err == nil {
					break
				} else {
					log.Printf("[ERR] can't add client, retry %d, error: %s", i+1, err)
					break
				}
			}
			log.Printf("[ERR] can't update client, retry %d, error: %s", i+1, err)
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

	if err != nil {
		return e.Wrap(fmt.Sprintf("can't update client after %d retries", s.retries), err)
	}

	// Save subscription only after it have processed by x-ui API
	return s.storage.SaveSubscription(ctx, sub)
}

func (s *XUIService) DeleteUserSubscription(ctx context.Context, serverID storage.ServerID, userID storage.UserID) (err error) {
	defer func() { e.WrapIfErr("can't delete subscription", err) }()

	if userID == 0 {
		return ErrZeroUserID
	}

	if serverID == 0 {
		return ErrZeroServerID
	}

	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	xui, err := s.xuiClientInstance(ctx, serverID)
	if err != nil {
		return e.Wrap("can't get xui client instance", err)
	}

	uuid := s.uuidFromTelegramID(user.TelegramID)
	err = xui.DeleteClient(ctx, s.inboundID, uuid)
	if err != nil {
		return e.Wrap("can't delete client on remote", err)
	}

	return s.storage.RemoveSubscriptionByIDs(ctx, userID, serverID)
}

func (s *XUIService) GetClientByIDs(ctx context.Context, serverID storage.ServerID, userID storage.UserID) (clientPtr *model.Client, err error) {
	defer func() { e.WrapIfErr("can't get client", err) }()

	xui, err := s.xuiClientInstance(ctx, serverID)
	if err != nil {
		return nil, e.Wrap("can't get xui client instance", err)
	}

	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	clientID := s.ClientIDByUser(user)

	inbound, err := xui.GetInbound(ctx, s.inboundID)
	if err != nil {
		return nil, err
	}

	client, err := s.GetClient(inbound, clientID)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// Returns Client from Inbound by ClientID.
func (s *XUIService) GetClient(inbound *model.Inbound, id ClientID) (model.Client, error) {
	clients, err := s.GetClients(inbound)
	if err != nil {
		return model.Client{}, err
	}

	for _, client := range clients {
		if client.ID == id.String() {
			return client, nil
		}
	}

	return model.Client{}, ErrClientNotFound

}

// @see vendor/github.com/MHSanaei/3x-ui/web/service/inbound.go
// @link github.com/MHSanaei/3x-ui/blob/main/web/service/inbound.go
// Return slice of Clients from Inbound.
func (s *XUIService) GetClients(inbound *model.Inbound) ([]model.Client, error) {
	settings := map[string][]model.Client{}
	json.Unmarshal([]byte(inbound.Settings), &settings)
	if settings == nil {
		return nil, fmt.Errorf("setting is null")
	}

	clients := settings["clients"]
	if clients == nil {
		return nil, nil
	}
	return clients, nil
}

func (s *XUIService) ClientIDByUser(user *storage.User) ClientID {
	return s.uuidFromTelegramID(user.TelegramID)
}

func (s *XUIService) ClientEmailByUser(user *storage.User) string {
	return s.ClientIDByUser(user).String()
}

func (s *XUIService) ClientSubIDByUser(user *storage.User) string {
	return s.ClientIDByUser(user).String()
}

// func (s *XUIService) validateSubscriptions(ctx context.Context, sub *storage.Subscription) (err error) {

// }

func (s *XUIService) uuidFromTelegramID(id storage.TelegramID) uuid.UUID {
	return s.uuidFromInt64(int64(id))
}

// Provide uuid.UUID from hashing int64.
// Hashing method is important for identify user in new server.
func (s *XUIService) uuidFromInt64(id int64) uuid.UUID {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))

	hash := sha256.Sum256(buf)

	u, err := uuid.FromBytes(hash[:16])
	if err != nil {
		panic(err)
	}

	return u
}

func (s *XUIService) xuiClientInstance(ctx context.Context, serverID storage.ServerID) (cli *x_ui.XUIClient, err error) {
	defer func() { e.WrapIfErr("can't get xui client instance", err) }()
	client, ok := s.clients[serverID]
	if !ok {
		// Mb defer would now work, because err variable don't make specifically?
		return s.createXUIClient(ctx, serverID)
	}
	return client.xui, nil
}

func (s *XUIService) createXUIClient(ctx context.Context, serverID storage.ServerID) (*x_ui.XUIClient, error) {
	server, err := s.storage.GetServerByID(ctx, serverID)
	if err != nil {
		return nil, e.Wrap("can't create xui client", err)
	}
	return x_ui.New(s.TokenKey, server, s.authStore), nil
}

func (s *XUIService) watchAndClearUnusedClients(alive_time, watch_delay time.Duration) {
	for {
		for serverID, v := range s.clients {
			if time.Since(v.updated_at) > alive_time {
				delete(s.clients, serverID)
			}
		}
		time.Sleep(watch_delay)
	}
}
