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

func (s *SQLStorage) GetServerByUserID(ctx context.Context, userID storage.UserID) (server []*storage.VPNServer, err error) {
	defer func() { e.WrapIfErr("can't get server by user id", err) }()

	q := `
		SELECT * FROM server AS s
		JOIN servers_servers ON servers.id = servers_servers.user_id AS us
		WHERE us.server_id = ?
	`

	err = s.db.GetContext(ctx, server, q, userID)

	return server, err
}

/* ---- Reader Interface implementation ---- */

func (s *SQLStorage) GetAllServers(ctx context.Context) (servers []*storage.VPNServer, err error) {
	defer func() { e.WrapIfErr("can't get all servers", err) }()

	q := `SELECT * FROM servers`

	err = s.db.SelectContext(ctx, servers, q)

	return servers, err
}

func (s *SQLStorage) GetServerByID(ctx context.Context, id storage.ServerID) (server *storage.VPNServer, err error) {
	defer func() { e.WrapIfErr("can't get user by id", err) }()

	q := `SELECT * FROM servers WHERE id = ? LIMIT 1`

	server = &storage.VPNServer{}
	err = s.db.GetContext(ctx, server, q, id)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSuchServer
	} else if err != nil {
		return nil, err
	}

	return server, nil
}

/* ---- Writer Interface implementation ---- */

func (s *SQLStorage) SaveServer(ctx context.Context, server *storage.VPNServer) (serverID storage.ServerID, err error) {
	defer func() { e.WrapIfErr("can't save server", err) }()
	var id int64

	q := `SELECT * FROM servers WHERE id = ? LIMIT 1`

	var result sql.Result
	oldServer := &storage.VPNServer{}

	err = s.db.GetContext(ctx, oldServer, q, server.ID)
	if err == sql.ErrNoRows {
		q = `INSERT INTO servers (country_id, name, protocol, host, port, username, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

		server.CreatedAt = time.Now()
		server.UpdatedAt = server.CreatedAt
		result, err = s.db.ExecContext(ctx, q, server.CountryID, server.Name, server.Protocol, server.Host, server.Port, server.Username, server.Password, server.CreatedAt, server.UpdatedAt)
		if err != nil {
			return 0, e.Wrap("can't execute query", err)
		}
		id, err = result.LastInsertId()
		if err != nil {
			return 0, e.Wrap("can't get last inserted id", err)
		}
		fmt.Println("INSERT query executed, id:", id)
		serverID = storage.ServerID(id)
	} else if err != nil {
		return 0, e.Wrap("can't scan row", err)
	} else {
		q = `UPDATE servers SET country_id = ?, name = ?, protocol = ?, host = ?, port = ?, username = ?, password = ?, updated_at = ? WHERE id = ?`

		server.ParseDefaultsFrom(oldServer)

		result, err = s.db.ExecContext(ctx, q, server.CountryID, server.Name, server.Protocol, server.Host, server.Port, server.Username, server.Password, server.UpdatedAt, server.ID)
		if err != nil {
			return 0, e.Wrap("can't execute query", err)
		}
		serverID = oldServer.ID
		fmt.Println("UPDATE query executed, id:", id)
	}

	return serverID, err
}

func (s *SQLStorage) RemoveServerByID(ctx context.Context, id storage.ServerID) error {
	q := `DELETE FROM servers WHERE id = ?`

	result, err := s.db.ExecContext(ctx, q, id)

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

// func (s *SQLStorage) scanServersRows(rows *sql.Rows) (servers []*storage.VPNServer, err error) {
// 	for rows.Next() {
// 		user := &storage.VPNServer{}
// 		err = rows.Scan(user)
// 		if err != nil {
// 			return nil, e.Wrap("can't scan user", err)
// 		}
// 		servers = append(servers, user)
// 	}
// 	return servers, nil
// }
