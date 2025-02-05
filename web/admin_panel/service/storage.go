package service

import (
	"context"
	"errors"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/internal/storage/sqlite"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
	"vpn-tg-bot/web/admin_panel/entity"
	sqlite_service "vpn-tg-bot/web/admin_panel/service/storage/sqlite"
)

var ErrNeedSQLCompatibleInterface = errors.New("SQL storage must be SQLCompatible")
var ErrDatabaseDriverIsNotSupported = errors.New("This database driver is not supported")

type StorageServiceQueries interface {
	GetUsersWithSubscription(ctx context.Context, args builder.Arguments) (*[]entity.UserWithSubscription, error)
	GetServersWithAuthorization(ctx context.Context, args builder.Arguments) (*[]entity.ServerWithAuthorization, error)
}

type StorageService interface {
	storage.Storage
	StorageServiceQueries
}

func NewStorageService(storage storage.Storage) (s StorageService, err error) {

	// Type asserting storage.Storage interface to concrete storage.
	// Then creating new Service, which shall implement both StorageServiceQueries
	// storage.Storage interfaces.
	// This is method to extend functional of storage.Storage interface.
	// If interface is not supported, then return error.
	switch storage.(type) {
	case *sqlite.SQLStorage:
		sqlStorage, ok := storage.(*sqlite.SQLStorage)
		if !ok {
			return nil, errors.New("can't provide type assertion to *sqlite.SQLStorage")
		}
		sqlService := sqlite_service.New(sqlStorage)

		s = sqlService
	default:
		return nil, ErrDatabaseDriverIsNotSupported
	}

	return s, nil
}
