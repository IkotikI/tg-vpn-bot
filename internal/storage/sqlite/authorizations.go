package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"
)

/* ---- Interface implementation ---- */
func (s *SQLStorage) SaveServerAuth(ctx context.Context, auth *storage.VPNServerAuthorization) (serverID storage.ServerID, err error) {
	defer func() { e.WrapIfErr("can't save server authorization", err) }()
	var id int64

	q := `SELECT * FROM servers_authorizations WHERE server_id = ? LIMIT 1`

	var result sql.Result
	oldAuth := &storage.VPNServerAuthorization{}

	err = s.db.GetContext(ctx, oldAuth, q, auth.ServerID)
	if err == sql.ErrNoRows {
		q := `INSERT INTO servers_authorizations (server_id, username, password, token, expired_at, meta) VALUES (?, ?, ?, ?, ?, ?)`

		result, err = s.db.ExecContext(ctx, q, auth.ServerID, auth.Username, auth.Password, auth.Token, auth.ExpiredAt, auth.Meta)
		if err != nil {
			return 0, e.Wrap("can't execute query", err)
		}

		id, err = result.LastInsertId()
		serverID = storage.ServerID(id)
	} else if err != nil {
		return 0, e.Wrap("can't scan row", err)
	} else {
		q := `UPDATE servers_authorizations SET username = ?, password = ?, token = ?, expired_at = ?, meta = ? WHERE server_id = ?`

		_, err = s.db.ExecContext(ctx, q, auth.Username, auth.Password, auth.Token, auth.ExpiredAt, auth.Meta, auth.ServerID)
		if err != nil {
			return 0, e.Wrap("can't execute query", err)
		}

		serverID = oldAuth.ServerID
	}

	return serverID, nil
}

func (s *SQLStorage) GetServerAuthByServerID(ctx context.Context, ServerID storage.ServerID) (auth *storage.VPNServerAuthorization, err error) {
	defer func() { e.WrapIfErr("can't get server authorization by server id", err) }()

	q := `SELECT * FROM servers_authorizations WHERE server_id = ? LIMIT 1`

	auth = &storage.VPNServerAuthorization{}
	err = s.db.GetContext(ctx, auth, q, ServerID)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSuchServerAuth
	} else if err != nil {
		return nil, err
	}

	return auth, nil
}

func (s *SQLStorage) GetServerAndAuthByServerID(ctx context.Context, ServerID storage.ServerID) (server *storage.VPNServer, auth *storage.VPNServerAuthorization, err error) {
	defer func() { e.WrapIfErr("can't get server and authorization by server id", err) }()

	q := `
		SELECT * FROM servers AS s
		JOIN servers_authorizationss AS sa ON sa.server_id = s.id
		WHERE server_id = ? LIMIT 1
	`

	server = &storage.VPNServer{}
	auth = &storage.VPNServerAuthorization{}
	row := s.db.QueryRowxContext(ctx, q, ServerID)
	err = row.Err()
	if err == sql.ErrNoRows {
		return nil, nil, storage.ErrNoSuchServerAuth
	} else if err != nil {
		return nil, nil, e.Wrap("can't execute query", row.Err())
	}

	err = row.StructScan(&server)
	if err != nil {
		return nil, nil, e.Wrap("can't scan server", err)
	}

	err = row.StructScan(&auth)
	if err != nil {
		return nil, nil, e.Wrap("can't scan authorization", err)
	}

	return server, auth, nil
}

func (s *SQLStorage) RemoveServerAuthByServerID(ctx context.Context, ServerID storage.ServerID) (err error) {
	defer func() { e.WrapIfErr("can't remove server authorization by server id", err) }()

	q := `
		DELETE FROM servers_authorizations WHERE server_id = ?
	`

	result, err := s.db.ExecContext(ctx, q, ServerID)

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
