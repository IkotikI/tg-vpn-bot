package sqlite_service

// import (
// 	"context"
// 	"log"
// 	"vpn-tg-bot/internal/entity"
// 	"vpn-tg-bot/internal/storage"
// 	"vpn-tg-bot/pkg/e"
// 	"vpn-tg-bot/pkg/sqlbuilder"
// 	"vpn-tg-bot/pkg/sqlbuilder/builder"
// )

// type SQLiteStorageService struct {
// 	storage.SQLStorage
// 	db      *sqlx.DB
// 	builder *builder.SQLBuilder
// }

// func New(storage storage.SQLStorage) *SQLiteStorageService {
// 	sqlInstance, err := storage.SQLStorageInstance()
// 	if err != nil {
// 		log.Fatalf("[ERR] Can't get SQL instance: %v", err)
// 	}

// 	sqlxDB := sqlx.NewDb(sqlInstance, "sqlite3")

// 	builder, err := sqlbuilder.NewSQLBuilder("sqlite3")
// 	if err != nil {
// 		log.Fatalf("[ERR] Can't create sqlite builder: %v", err)
// 	}

// 	return &SQLiteStorageService{
// 		SQLStorage: storage,
// 		db:         sqlxDB,
// 		builder:    builder,
// 	}
// }

// func (s *SQLiteStorageService) ChangeSubscriptionStatus(ctx context.Context, args builder.Arguments) (subscription *storage.Subscription, err error) {
// 	defer func() { e.WrapIfErr("can't change subscription status", err) }()

// 	return nil, nil
// }
