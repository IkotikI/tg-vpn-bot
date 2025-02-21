package entity

import "vpn-tg-bot/internal/storage"

var TimeLayout = "2006-01-02 15:04:05"

type UserWithSubscription struct {
	storage.User
	storage.Subscription
}

type ServerWithAuthorization struct {
	storage.VPNServer
	storage.VPNServerAuthorization
	storage.Country
}

type User struct {
	storage.User
}

type Server struct {
	storage.VPNServer
	storage.Country
}

type SubscriptionWithServer struct {
	Server
	storage.Subscription
}

type SubscriptionWithUser struct {
	User
	storage.Subscription
}
