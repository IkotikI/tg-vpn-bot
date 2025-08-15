package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"
	"vpn-tg-bot/pkg/structconv"

	"github.com/jmoiron/sqlx"
)

func (s *SQLStorage) GetServerWithCountryByID(ctx context.Context, serverID storage.ServerID) (server *storage.VPNServerWithCountry, err error) {
	q := `
		SELECT * FROM servers AS s
		JOIN countries AS c
		ON s.country_id = c.country_id
		WHERE id = ? LIMIT 1
	`

	server = &storage.VPNServerWithCountry{}

	err = GetContextWithNullFallback(ctx, s.db, server, q, serverID)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *SQLStorage) GetServersWithCountries(ctx context.Context, args *storage.QueryArgs) (servers *[]storage.VPNServerWithCountry, err error) {
	q := `
		SELECT * FROM servers AS s
		LEFT JOIN countries AS c
		ON s.country_id = c.country_id
	`

	queryEnd, queryArgs := s.buildParts([]string{"where", "order_by", "limit"}, args)
	q += queryEnd
	log.Printf("query: `%s` args: %+v\n", q, queryArgs)

	servers = &[]storage.VPNServerWithCountry{}
	err = SelectContextWithNullFallback(ctx, s.db, servers, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *SQLStorage) GetSubscriptionWithUserAndServerByIDs(ctx context.Context, userID storage.UserID, serverID storage.ServerID) (sub *storage.SubscriptionWithUserAndServer, err error) {
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

	sub = &storage.SubscriptionWithUserAndServer{}
	err = GetContextWithNullFallback(ctx, s.db, sub, q, userID, serverID)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSuchSubscription
	} else if err != nil {
		return nil, err
	}

	return sub, err
}

func (s *SQLStorage) GetSubscriptionsWithUsersAndServers(ctx context.Context, args *storage.QueryArgs) (subs *[]storage.SubscriptionWithUserAndServer, err error) {
	q := `
		SELECT * FROM subscriptions AS sub
		LEFT JOIN users AS u
		ON sub.user_id = u.id
		LEFT JOIN servers AS serv 
		ON sub.server_id = serv.id
		LEFT JOIN countries AS c
		ON serv.country_id = c.country_id
	`

	queryEnd, queryArgs := s.buildParts([]string{"where", "order_by", "limit"}, args)
	q += queryEnd

	subs = &[]storage.SubscriptionWithUserAndServer{}
	err = SelectContextWithNullFallback(ctx, s.db, subs, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return subs, err
}

// func (s *SQLStorage) GetSubscriptionsWithUsersAndServers(ctx context.Context, args *storage.QueryArgs) (subs *[]storage.SubscriptionWithUserAndServer, errs []error) {
// 	errs = make([]error, 3)
// 	q := `
// 		SELECT * FROM subscriptions
// 	`
// 	queryEnd, queryArgs := s.buildParts([]string{"where", "order_by", "limit"}, args)
// 	q += queryEnd

// 	subscriptionsPtr := &[]storage.Subscription{}
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
// 	queryEnd, queryArgsAdd := s.buildParts([]string{"order_by"}, args)
// 	q += queryEnd
// 	queryArgs

// 	usersPtr := &[]storage.User{}
// 	errs[1] = SelectContextWithNullFallback(ctx, s.db, usersPtr, q, []interface{}{usersIDs})

// 	q = fmt.Sprintf(`
// 		SELECT * FROM servers
// 		WHERE id IN (%s)
// 	`, questions)
// 	serversPtr := &[]storage.VPNServerWithCountry{}
// 	errs[2] = SelectContextWithNullFallback(ctx, s.db, serversPtr, q, []interface{}{serversIDs})

// 	usersOk := errs[1] == nil
// 	serverOk := errs[2] == nil
// 	for i, sub := range subscriptions {
// 		if usersOk {

// 		}
// 	}

// 	for i, sub := range subscriptions {
// 		subs[i] = storage.SubscriptionWithUserAndServer{
// 			Subscription: sub,
// 		}
// 		if usersOk {
// 			subs[i].User =
// 		}
// 	}

// 	return subs, nil
// }

func (s *SQLStorage) GetSubscriptionsWithServersByUserID(ctx context.Context, user_id storage.UserID, args *storage.QueryArgs) (subs *[]storage.SubscriptionWithServer, err error) {
	q := `
		SELECT * FROM subscriptions AS sub
		JOIN servers AS serv 
		ON sub.server_id = serv.id
		JOIN countries AS c
		ON serv.country_id = c.country_id
		WHERE sub.user_id = ?
	`
	queryArgs := []interface{}{user_id}
	queryEnd, queryArgsAdd := s.buildParts([]string{"order_by", "limit"}, args)
	q += queryEnd
	queryArgs = append(queryArgs, queryArgsAdd...)
	log.Printf("query: `%s` args: %+v\n", q, queryArgs)

	subs = &[]storage.SubscriptionWithServer{}
	err = SelectContextWithNullFallback(ctx, s.db, subs, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SQLStorage) GetSubscriptionsWithUsersByServerID(ctx context.Context, server_id storage.ServerID, args *storage.QueryArgs) (subs *[]storage.SubscriptionWithUser, err error) {
	q := `
		SELECT * FROM subscriptions AS s
		JOIN users AS u
		ON s.user_id = u.id
		WHERE s.server_id = ?
	`

	queryArgs := []interface{}{server_id}
	queryEnd, queryArgsAdd := s.buildParts([]string{"order_by", "limit"}, args)
	q += queryEnd
	queryArgs = append(queryArgs, queryArgsAdd...)
	log.Printf("query: `%s` args: %+v\n", q, queryArgs)

	subs = &[]storage.SubscriptionWithUser{}
	err = SelectContextWithNullFallback(ctx, s.db, subs, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SQLStorage) CountWithBuilder(ctx context.Context, args *storage.QueryArgs) (n int64, err error) {
	q := `
		SELECT count(*) AS count
	`

	queryEnd, queryArgs := s.buildParts([]string{"from", "where"}, args)
	q += queryEnd

	log.Printf("query: `%s` args: %+v\n", q, queryArgs)
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
