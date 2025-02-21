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

func (s *SQLiteStorageService) GetUsers(ctx context.Context, args builder.Arguments) (users *[]entity.User, err error) {
	defer func() { e.WrapIfErr("can't get users", err) }()

	q := `SELECT * FROM users`

	qEnd, qArgs := s.builder.BuildParts([]string{"where", "order_by", "limit"}, args)
	q += qEnd
	log.Printf("query: `%s` args: %+v", q, qArgs)

	users = &[]entity.User{}
	err = SelectContextWithNullFallback(ctx, s.db, users, q, qArgs...)
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (s *SQLiteStorageService) GetServers(ctx context.Context, args builder.Arguments) (servers *[]entity.Server, err error) {
	defer func() { e.WrapIfErr("can't get users", err) }()

	q := `
		SELECT * FROM servers AS s
		JOIN countries AS c
		ON s.country_id = c.country_id
	`

	qEnd, qArgs := s.builder.BuildParts([]string{"where", "order_by", "limit"}, args)
	q += qEnd
	log.Printf("query: `%s` args: %+v", q, qArgs)

	servers = &[]entity.Server{}
	err = SelectContextWithNullFallback(ctx, s.db, servers, q, qArgs...)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *SQLiteStorageService) GetEntityServerByID(ctx context.Context, id storage.ServerID) (server *entity.Server, err error) {
	defer func() { e.WrapIfErr("can't get user by id", err) }()

	q := `
		SELECT * FROM servers AS s
		JOIN countries AS c
		ON s.country_id = c.country_id
		WHERE id = ? LIMIT 1
	`

	server = &entity.Server{}

	err = GetContextWithNullFallback(ctx, s.db, server, q, id)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *SQLiteStorageService) GetSubscriptionsWithServersByUserID(ctx context.Context, user_id storage.UserID, args builder.Arguments) (subs *[]entity.SubscriptionWithServer, err error) {
	defer func() { e.WrapIfErr("can't get users", err) }()

	q := `
		SELECT * FROM subscriptions AS sub
		JOIN servers AS serv 
		ON sub.server_id = serv.id
		JOIN countries AS c
		ON serv.country_id = c.country_id
		WHERE sub.user_id = ?
	`
	qArgs := []interface{}{user_id}
	qEnd, qArgsAdd := s.builder.BuildParts([]string{"order_by", "limit"}, args)
	q += qEnd
	qArgs = append(qArgs, qArgsAdd...)
	log.Printf("query: `%s` args: %+v", q, qArgs)

	subs = &[]entity.SubscriptionWithServer{}
	err = SelectContextWithNullFallback(ctx, s.db, subs, q, qArgs...)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SQLiteStorageService) GetSubscriptionsWithUsersByServerID(ctx context.Context, server_id storage.ServerID, args builder.Arguments) (subs *[]entity.SubscriptionWithUser, err error) {
	defer func() { e.WrapIfErr("can't get users", err) }()

	q := `
		SELECT * FROM subscriptions AS s
		JOIN users AS u
		ON s.user_id = u.id
		WHERE s.server_id = ?
	`

	qArgs := []interface{}{server_id}
	qEnd, qArgsAdd := s.builder.BuildParts([]string{"order_by", "limit"}, args)
	q += qEnd
	qArgs = append(qArgs, qArgsAdd...)
	log.Printf("query: `%s` args: %+v", q, qArgs)

	subs = &[]entity.SubscriptionWithUser{}
	err = SelectContextWithNullFallback(ctx, s.db, subs, q, qArgs...)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func GetContextWithNullFallback[T any](ctx context.Context, db *sqlx.DB, dest *T, query string, args ...interface{}) (err error) {
	row := db.QueryRowxContext(ctx, query, args...)
	return structScanWithNullFallback(row, dest)
}

func SelectContextWithNullFallback[T any](ctx context.Context, db *sqlx.DB, dest *[]T, query string, args ...interface{}) (err error) {
	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}
	_values := *dest
	for rows.Next() {
		value := *new(T)
		err := structScanWithNullFallback(rows, &value)
		if err != nil {
			log.Printf("[ERR] Can't scan value: %v", err)
		}

		_values = append(_values, value)
	}
	*dest = _values
	// log.Printf("Type of selected values: %T, _values %+v", _values, _values)
	return nil
}

type StructScannable interface {
	StructScan(dest interface{}) error
}

// This function make support queries for Null fields in sql database, providing
// fallback to specific transformation of original struct to struct with basic fields
// replaces with sql.Null_. This behavior can affect performance and using for fallback,
// not just for main functional of app. [WARINING] will be logged, when fallback
// is used.
// rows - pointer to sqlx.Row or sqlx.Rows
// v - pointer to struct.
func structScanWithNullFallback(rows StructScannable, v interface{}) (err error) {
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
