package sqlite_service

import (
	"context"
	"log"
	"strings"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"
	"vpn-tg-bot/pkg/sqlbuilder"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
	"vpn-tg-bot/pkg/structconv"
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

	qEnd, qArgs := s.builder.BuildParts([]string{"where", "order_by"}, args)
	q += qEnd

	// Make slice of UsersWithSubscription with values, changed to sql.Null types.
	_users := []entity.UserWithSubscription{}
	rows, err := s.db.QueryxContext(ctx, q, qArgs...)
	for rows.Next() {
		user := entity.UserWithSubscription{}
		err := structScanRowsWithNullFallback(rows, &user)
		if err != nil {
			log.Printf("[ERR] Can't scan user: %v", err)
		}

		_users = append(_users, user)
	}
	users = &_users

	return users, err
}

func (s *SQLiteStorageService) GetServersWithAuthorization(ctx context.Context, args builder.Arguments) (servers *[]entity.ServerWithAuthorization, err error) {
	defer func() { e.WrapIfErr("can't get users with subscription", err) }()

	q := `
		SELECT * FROM servers AS s
		LEFT OUTER JOIN servers_authorizations AS sa
		ON s.id = sa.server_id
	`

	qEnd, qArgs := s.builder.BuildParts([]string{"where", "order_by"}, args)
	q += qEnd

	// servers = &[]entity.ServerWithAuthorization{}
	// err = s.db.SelectContext(ctx, servers, q, qArgs...)
	_servers := []entity.ServerWithAuthorization{}
	rows, err := s.db.QueryxContext(ctx, q, qArgs...)
	for rows.Next() {
		server := entity.ServerWithAuthorization{}
		err := structScanRowsWithNullFallback(rows, &server)
		if err != nil {
			log.Printf("[ERR] Can't scan user: %v", err)
		}

		_servers = append(_servers, server)
	}
	servers = &_servers

	return servers, err
}

// This function make support queries for Null fields in sql database, providing
// fallback to specific transformation of original struct to struct with basic fields
// replaces with sql.Null_. This behavior can affect performance and using for fallback,
// not just for main functional of app. [WARINING] will be logged, when fallback
// is used.
// v - pointer to struct.
func structScanRowsWithNullFallback(rows *sqlx.Rows, v interface{}) (err error) {
	err = rows.StructScan(v)
	if err != nil {
		log.Printf("[WARNING] Can't scan value: %v. Try to scan to sql.Null struct.", err)
		if strings.Contains(err.Error(), "converting NULL") {
			userNull := structconv.CreateSQLNullStruct(v)
			// fmt.Printf("userNull %+v\n", userNull)
			// structconv.PrintStructFields(userNull)
			err = rows.StructScan(userNull)
			if err != nil {
				return e.Wrap("Can't scan user to struct with sql.Null_ values", err)
			} else {
				structconv.ConvertSQLNullStructToBasic(userNull, v)
			}
		} else {
			return e.Wrap("can't scan user", err)
		}
	}
	return nil
}
