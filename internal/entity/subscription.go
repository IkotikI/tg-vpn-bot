package entity

import "time"

type Subscription struct {
	UserID                int
	ServerID              int
	SubscriptionStatus    string
	SubscriptionExpiredAt time.Time
}
