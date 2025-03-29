package service

import (
	"context"
	"errors"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/internal/storage/sqlite"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
	sqlite_service "vpn-tg-bot/web/admin_panel/service/storage/sqlite"
)

type StorageServiceQueries interface {
	GetUsers(ctx context.Context, args builder.Arguments) (*[]storage.User, error)
	GetServers(ctx context.Context, args builder.Arguments) (*[]storage.VPNServerWithCountry, error)
	GetServerByID(ctx context.Context, serverID storage.ServerID) (*storage.VPNServerWithCountry, error)
	GetSubscriptionsWithServersByUserID(ctx context.Context, user_id storage.UserID, args builder.Arguments) (*[]storage.SubscriptionWithServer, error)
	GetSubscriptionsWithUsersByServerID(ctx context.Context, user_id storage.ServerID, args builder.Arguments) (*[]storage.SubscriptionWithUser, error)
	GetSubscriptionsWithUsersAndServers(ctx context.Context, args builder.Arguments) (*[]storage.SubscriptionWithUserAndServer, error)
	GetSubscriptionWithUserAndServerByIDs(ctx context.Context, userID storage.UserID, serverID storage.ServerID) (*storage.SubscriptionWithUserAndServer, error)
	CountWithBuilder(ctx context.Context, args builder.Arguments) (n int64, err error)
}

type StorageService interface {
	// storage.Storage
	StorageServiceQueries
}

func NewStorageService(store storage.Storage) (service StorageService, err error) {

	// Type asserting storage.Storage interface to concrete storage.
	// Then creating new Service, which shall implement both StorageServiceQueries
	// storage.Storage interfaces.
	// This is method to extend functional of storage.Storage interface.
	// If interface is not supported, then return error.
	switch store.(type) {
	case *sqlite.SQLStorage:
		sqlStorage, ok := store.(*sqlite.SQLStorage)
		if !ok {
			return nil, errors.New("can't provide type assertion to *sqlite.SQLStorage")
		}
		sqlService := sqlite_service.New(sqlStorage)

		service = sqlService
	default:
		return nil, storage.ErrDatabaseDriverIsNotSupported
	}

	return service, nil
}
