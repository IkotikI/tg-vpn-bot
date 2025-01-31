package entity

import "time"

type UserID int64
type TelegramID int64

type User struct {
	ID           UserID
	TelegramID   TelegramID
	TelegramName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
