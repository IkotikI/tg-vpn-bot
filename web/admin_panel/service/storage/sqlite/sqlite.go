package sqlite_service

import (
	"context"
	"log"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"
	"vpn-tg-bot/pkg/sqlbuilder"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
	"vpn-tg-bot/web/admin_panel/entity"

	"github.com/jmoiron/sqlx"
)

type SQLiteStorageService struct {
	storage.SQLStorage
	db      *sqlx.DB
	builder *builder.SQLBuilder
}

func New(storage storage.SQLStorage) *SQLiteStorageService {
	sqlInstance, err := storage.SQLStorageInstance()
	if err != nil {
		log.Fatalf("[ERR] Can't get SQL instance: %v", err)
	}

	sqlxDB := sqlx.NewDb(sqlInstance, "sqlite3")

	builder, err := sqlbuilder.NewSQLBuilder("sqlite3")
	if err != nil {
		log.Fatalf("[ERR] Can't create sqlite builder: %v", err)
	}

	return &SQLiteStorageService{
		SQLStorage: storage,
		db:         sqlxDB,
		builder:    builder,
	}
}

func (s *SQLiteStorageService) GetUsersWithSubscription(ctx context.Context, args builder.Arguments) (users *[]entity.UserWithSubscription, err error) {
	defer func() { e.WrapIfErr("can't get users with subscription", err) }()

	q := `
	SELECT * FROM users AS u
	LEFT OUTER JOIN subscriptions AS us
	ON u.id = us.user_id
	`
	q += s.builder.BuildParts([]string{"where", "order_by"}, args)

	users = &[]entity.UserWithSubscription{}
	err = s.db.SelectContext(ctx, users, q)

	return users, err
}

func (s *SQLiteStorageService) GetServersWithAuthorization(ctx context.Context, args builder.Arguments) (servers *[]entity.ServerWithAuthorization, err error) {
	defer func() { e.WrapIfErr("can't get users with subscription", err) }()

	q := `
	SELECT * FROM servers AS s
	LEFT OUTER JOIN servers_authorizations AS sa
	ON s.id = sa.server_id
	`

	servers = &[]entity.ServerWithAuthorization{}
	err = s.db.SelectContext(ctx, servers, q)

	return servers, err
}
