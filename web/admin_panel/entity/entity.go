package entity

import "vpn-tg-bot/internal/storage"

type UserWithSubscription struct {
	storage.User
	storage.Subscription
}

type ServerWithAuthorization struct {
	storage.VPNServer
	storage.VPNServerAuthorization
	storage.Country
}
