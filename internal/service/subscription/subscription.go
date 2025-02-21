package subscription

import (
	"context"
	"errors"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"
)

type VPN_API interface {
	UpdateSubscription(ctx context.Context, sub *storage.Subscription) (err error)
}

type SubscriptionServiceQueries interface {
}

type SubscriptionService struct {
	client  VPN_API
	storage storage.Subscriptions
}

func NewStorageService(store storage.Subscriptions, client VPN_API) (service *SubscriptionService) {

	// Type asserting storage.Storage interface to concrete storage.
	// Then creating new Service, which shall implement both StorageServiceQueries
	// storage.Storage interfaces.
	// This is method to extend functional of storage.Storage interface.
	// If interface is not supported, then return error.
	// switch store.(type) {
	// case *sqlite.SQLStorage:
	// 	sqlStorage, ok := store.(*sqlite.SQLStorage)
	// 	if !ok {
	// 		return nil, errors.New("can't provide type assertion to *sqlite.SQLStorage")
	// 	}
	// 	subsStorage :=

	// 	service.storage, err = sqlite.New(sqlStorage)
	// default:
	// 	return nil, storage.ErrDatabaseDriverIsNotSupported
	// }

	// return service, nil

	return &SubscriptionService{
		storage: store,
		client:  client,
	}
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, sub *storage.Subscription) (err error) {
	defer func() { e.WrapIfErr("can't update subscription", err) }()

	if sub.UserID == 0 {
		return errors.New("user id is 0")
	}

	if sub.ServerID == 0 {
		return errors.New("server id is 0")
	}

	return nil

}
