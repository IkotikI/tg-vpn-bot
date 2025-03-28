package auth

import (
	"context"
	"vpn-tg-bot/internal/storage"
)

type AuthorizationService struct {
	storage storage.Storage
}

func New(storage storage.Storage) *AuthorizationService {
	return &AuthorizationService{storage: storage}
}

func RegisterUser(ctx context.Context, telegramID storage.TelegramID) (userID storage.UserID, err error) {
	return 0, nil
}
