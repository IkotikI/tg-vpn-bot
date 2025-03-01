package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"
)

/* ---- Queries Interface implementation ---- */

func (s *SQLStorage) GetUserByServerID(ctx context.Context, serverID storage.ServerID) (users []*storage.User, err error) {
	q := `
		SELECT * FROM users AS u
		JOIN users_servers ON users.id = users_servers.user_id AS us
		WHERE us.server_id = ?
	`

	err = s.db.SelectContext(ctx, users, q, serverID)
	if err != nil {
		return nil, e.Wrap("can't execute query", err)
	}

	return users, nil
}

/* ---- Reader Interface implementation ---- */
// TODO: Method query is bad for selection users, because it breaks encapsulation of storage
// interface. Better approach don't provide general query method as well as don't query full table
// of users. Select queries must be always limited by category, or, at least, total number.
// Uncontrolled select query can overfill memory. In this implementation I shall operate just
// defined in storage.go abstractions or basic Golang types.
//
// where, order_by, limit - are 3 pillars of making SQL selection by given table (join).
// Where can accept: column name, operator, value(s)
// Order_by can accept: column name, order (ASC DESC)
// Limit can accept: offset, limit

func (s *SQLStorage) GetUsers(ctx context.Context, args *storage.QueryArgs) (users *[]storage.User, err error) {

	selectArgs := s.parseQueryArgs(args)

	q := `SELECT * FROM users`

	queryEnd, queryArgs := s.builder.BuildParts([]string{"where", "order_by", "limit"}, selectArgs)
	q += queryEnd

	users = &[]storage.User{}
	err = s.db.SelectContext(ctx, users, q, queryArgs...)

	return users, err
}

func (s *SQLStorage) GetUserByID(ctx context.Context, id storage.UserID) (user *storage.User, err error) {
	q := `SELECT * FROM users WHERE id = ? LIMIT 1`

	user = &storage.User{}
	err = s.db.GetContext(ctx, user, q, id)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSuchUser
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

/* ---- Writer Interface implementation ---- */

func (s *SQLStorage) SaveUser(ctx context.Context, user *storage.User) (userID storage.UserID, err error) {
	var id int64

	q := `SELECT * FROM users WHERE id = ? OR telegram_id = ? LIMIT 1`

	var result sql.Result
	oldUser := &storage.User{}

	err = s.db.GetContext(ctx, oldUser, q, user.ID, user.TelegramID)
	if err == sql.ErrNoRows {
		if user.TelegramID == 0 {
			return 0, errors.New("user with this id or telegram_id doesn't exist; creating new user require non-empty telegram_id")
		}
		q = `INSERT INTO users (telegram_id, telegram_name, created_at, updated_at) VALUES (?, ?, ?, ?)`

		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt

		result, err = s.db.ExecContext(ctx, q, user.TelegramID, user.TelegramName, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			return 0, e.Wrap("can't execute query", err)
		}

		id, err = result.LastInsertId()
		if err != nil {
			return 0, e.Wrap("can't get last inserted id", err)
		}
		userID = storage.UserID(id)
		fmt.Println("INSERT query executed, id:", id)
	} else if err != nil {
		return 0, e.Wrap("can't scan row", err)
	} else {
		q = `UPDATE users SET telegram_name = ?, updated_at = ? WHERE telegram_id = ?`

		if user.TelegramName == "" {
			user.TelegramName = oldUser.TelegramName
		}
		if user.UpdatedAt.IsZero() {
			user.UpdatedAt = time.Now()
		}

		result, err = s.db.ExecContext(ctx, q, user.TelegramName, user.UpdatedAt, oldUser.TelegramID)
		if err != nil {
			return 0, e.Wrap("can't execute query", err)
		}

		userID = oldUser.ID
		fmt.Println("UPDATE query executed, id:", id)
	}

	return userID, nil
}

func (s *SQLStorage) RemoveUserByID(ctx context.Context, id storage.UserID) (err error) {
	q := "DELETE FROM users WHERE id = ?"

	var result sql.Result
	result, err = s.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	var n int64
	n, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("zero rows affected")
	}

	return nil
}

/* ---- Helper functions ---- */

// func (s *SQLStorage) scanUsersRows(rows *sqlx.Rows) (users []*storage.User, err error) {
// 	for rows.Next() {
// 		user := &storage.User{}
// 		err = rows.Scan(user)
// 		if err != nil {
// 			return nil, e.Wrap("can't scan user", err)
// 		}
// 		users = append(users, user)
// 	}
// 	return users, nil
// }
