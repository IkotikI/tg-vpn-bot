package sqlite_service

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *SQLiteStorageService) GetEntityUsers(ctx context.Context, args builder.Arguments) (users *[]entity.User, err error) {
	q := `SELECT * FROM users`

	queryEnd, queryArgs := s.builder.BuildParts([]string{"where", "order_by", "limit"}, args)
	q += queryEnd
	log.Printf("query: `%s` args: %+v", q, queryArgs)

	users = &[]entity.User{}
	err = SelectContextWithNullFallback(ctx, s.db, users, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (s *SQLiteStorageService) GetEntityServers(ctx context.Context, args builder.Arguments) (servers *[]entity.Server, err error) {
	q := `
		SELECT * FROM servers AS s
		LEFT JOIN countries AS c
		ON s.country_id = c.country_id
	`

	queryEnd, queryArgs := s.builder.BuildParts([]string{"where", "order_by", "limit"}, args)
	q += queryEnd
	log.Printf("query: `%s` args: %+v", q, queryArgs)

	servers = &[]entity.Server{}
	err = SelectContextWithNullFallback(ctx, s.db, servers, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *SQLiteStorageService) GetEntityServerByID(ctx context.Context, id storage.ServerID) (server *entity.Server, err error) {
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

// TODO: Move to storage & make functions to get Users and Servers by ID slice
func (s *SQLiteStorageService) GetSubscriptions(ctx context.Context, args builder.Arguments) (subs *[]entity.Subscription, err error) {
	q := `
		SELECT * FROM subscriptions
	`
	queryEnd, queryArgs := s.builder.BuildParts([]string{"where", "order_by", "limit"}, args)
	q += queryEnd

	subs = &[]entity.Subscription{}
	err = SelectContextWithNullFallback(ctx, s.db, subs, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SQLiteStorageService) GetSubscriptionWithUserAndServerByIDs(ctx context.Context, userID storage.UserID, serverID storage.ServerID) (sub *entity.SubscriptionWithUserAndServer, err error) {
	q := `
		SELECT * FROM subscriptions AS sub
		LEFT JOIN users AS u
		ON sub.user_id = u.id
		LEFT JOIN servers AS serv 
		ON sub.server_id = serv.id
		LEFT JOIN countries AS c
		ON serv.country_id = c.country_id
		WHERE user_id = ? AND server_id = ?
		LIMIT 1
	`

	sub = &entity.SubscriptionWithUserAndServer{}
	err = GetContextWithNullFallback(ctx, s.db, sub, q, userID, serverID)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSuchSubscription
	} else if err != nil {
		return nil, err
	}

	return sub, err
}

func (s *SQLiteStorageService) GetSubscriptionsWithUsersAndServers(ctx context.Context, args builder.Arguments) (subs *[]entity.SubscriptionWithUserAndServer, err error) {
	q := `
		SELECT * FROM subscriptions AS sub
		LEFT JOIN users AS u
		ON sub.user_id = u.id
		LEFT JOIN servers AS serv 
		ON sub.server_id = serv.id
		LEFT JOIN countries AS c
		ON serv.country_id = c.country_id
	`

	queryEnd, queryArgs := s.builder.BuildParts([]string{"where", "order_by", "limit"}, args)
	q += queryEnd

	subs = &[]entity.SubscriptionWithUserAndServer{}
	err = SelectContextWithNullFallback(ctx, s.db, subs, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return subs, err
}

// func (s *SQLiteStorageService) GetSubscriptionsWithUsersAndServers(ctx context.Context, args builder.Arguments) (subs *[]entity.SubscriptionWithUserAndServer, errs []error) {
// 	errs = make([]error, 3)
// 	q := `
// 		SELECT * FROM subscriptions
// 	`
// 	queryEnd, queryArgs := s.builder.BuildParts([]string{"where", "order_by", "limit"}, args)
// 	q += queryEnd

// 	subscriptionsPtr := &[]entity.Subscription{}
// 	errs[0] = SelectContextWithNullFallback(ctx, s.db, subscriptionsPtr, q, queryArgs...)
// 	if errs[0] != nil {
// 		return nil, errs
// 	}

// 	subscriptions := *subscriptionsPtr
// 	length := len(subscriptions)
// 	usersIDs := make([]storage.UserID, length)
// 	serversIDs := make([]storage.ServerID, length)

// 	for i, sub := range subscriptions {
// 		usersIDs[i] = sub.UserID
// 		serversIDs[i] = sub.ServerID
// 	}

// 	questions := strings.Repeat(",?", length)[1:]

// 	var orderBy builder.OrderBy
// 	var order string = "ASC"
// 	orderBy,ok := args["order_by"].(builder.OrderBy)
// 	if ok {
// 		order = orderBy.Order
// 		if order != "ASC" || order != "DESC" {
// 			order = "ASC"
// 		}
// 	}

// 	q = fmt.Sprintf(`
// 		SELECT * FROM servers
// 		WHERE id IN (%s)
// 	`, questions)
// 	queryArgs = []interface{}{serversIDs}
// 	queryEnd, queryArgsAdd := s.builder.BuildParts([]string{"order_by"}, args)
// 	q += queryEnd
// 	queryArgs

// 	usersPtr := &[]entity.User{}
// 	errs[1] = SelectContextWithNullFallback(ctx, s.db, usersPtr, q, []interface{}{usersIDs})

// 	q = fmt.Sprintf(`
// 		SELECT * FROM servers
// 		WHERE id IN (%s)
// 	`, questions)
// 	serversPtr := &[]entity.Server{}
// 	errs[2] = SelectContextWithNullFallback(ctx, s.db, serversPtr, q, []interface{}{serversIDs})

// 	usersOk := errs[1] == nil
// 	serverOk := errs[2] == nil
// 	for i, sub := range subscriptions {
// 		if usersOk {

// 		}
// 	}

// 	for i, sub := range subscriptions {
// 		subs[i] = entity.SubscriptionWithUserAndServer{
// 			Subscription: sub,
// 		}
// 		if usersOk {
// 			subs[i].User =
// 		}
// 	}

// 	return subs, nil
// }

func (s *SQLiteStorageService) GetSubscriptionsWithServersByUserID(ctx context.Context, user_id storage.UserID, args builder.Arguments) (subs *[]entity.SubscriptionWithServer, err error) {
	q := `
		SELECT * FROM subscriptions AS sub
		JOIN servers AS serv 
		ON sub.server_id = serv.id
		JOIN countries AS c
		ON serv.country_id = c.country_id
		WHERE sub.user_id = ?
	`
	queryArgs := []interface{}{user_id}
	queryEnd, queryArgsAdd := s.builder.BuildParts([]string{"order_by", "limit"}, args)
	q += queryEnd
	queryArgs = append(queryArgs, queryArgsAdd...)
	log.Printf("query: `%s` args: %+v", q, queryArgs)

	subs = &[]entity.SubscriptionWithServer{}
	err = SelectContextWithNullFallback(ctx, s.db, subs, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SQLiteStorageService) GetSubscriptionsWithUsersByServerID(ctx context.Context, server_id storage.ServerID, args builder.Arguments) (subs *[]entity.SubscriptionWithUser, err error) {
	q := `
		SELECT * FROM subscriptions AS s
		JOIN users AS u
		ON s.user_id = u.id
		WHERE s.server_id = ?
	`

	queryArgs := []interface{}{server_id}
	queryEnd, queryArgsAdd := s.builder.BuildParts([]string{"order_by", "limit"}, args)
	q += queryEnd
	queryArgs = append(queryArgs, queryArgsAdd...)
	log.Printf("query: `%s` args: %+v", q, queryArgs)

	subs = &[]entity.SubscriptionWithUser{}
	err = SelectContextWithNullFallback(ctx, s.db, subs, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SQLiteStorageService) Count(ctx context.Context, args builder.Arguments) (n int64, err error) {
	q := `
		SELECT count(*) AS count 
	`

	queryEnd, queryArgs := s.builder.BuildParts([]string{"from", "where"}, args)
	q += queryEnd

	log.Printf("query: `%s` args: %+v", q, queryArgs)
	count := &[]int64{}
	err = s.db.SelectContext(ctx, count, q, queryArgs...)
	if err != nil {
		return -1, err
	}

	fmt.Printf("got count %+v", count)

	return (*count)[0], nil
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
