package entity

import (
	"database/sql"
	"vpn-tg-bot/internal/storage"
)

type NullInt64 struct {
	sql.NullInt64
}

type NullString struct {
	sql.NullString
}

type NullTime struct {
	sql.NullTime
}

type NullSubscription struct {
	UserID                NullInt64  `json:"user_id"`
	ServerID              NullInt64  `json:"server_id"`
	SubscriptionStatus    NullString `json:"subscription_status"`
	SubscriptionExpiredAt NullTime   `json:"subscription_expired_at"`
}

type UserWithNullSubscription struct {
	storage.User
	NullSubscription
}
