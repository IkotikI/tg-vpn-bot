package subscription

import (
	"context"
	"vpn-tg-bot/internal/storage"
)

type VPN_API interface {
	SubscriptionLink(ctx context.Context, serverID storage.ServerID, userID storage.UserID) (link string, err error)
	UpdateSubscription(ctx context.Context, sub *storage.Subscription) (err error)
	DeleteUserSubscription(ctx context.Context, serverID storage.ServerID, userID storage.UserID) (err error)
}
