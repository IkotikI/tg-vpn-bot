package sqlite

func (s *SQLStorage) createUsersTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            telegram_id TEXT NOT NULL,
            telegram_name TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL,
            updated_at TIMESTAMP NOT NULL
        )
    `)
	return err
}

func (s *SQLStorage) createServersTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS servers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            country_id INTEGER, /* ? */
            name TEXT NOT NULL,
            protocol TEXT NOT NULL,
            host TEXT NOT NULL,
            port INTEGER NOT NULL,
            username TEXT NOT NULL,
            password TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL,
            updated_at TIMESTAMP NOT NULL
        )
    `)
	return err
}

func (s *SQLStorage) createSubscriptionsTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS subscriptions (
            user_id INTEGER NOT NULL,
            server_id INTEGER NOT NULL,
            subscription_status TEXT NOT NULL, /* active | inactive */
            subscription_expired_at TIMESTAMP NOT NULL,
            UNIQUE(user_id, server_id)
        )
    `)
	return err
}

func (s *SQLStorage) createCountriesTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS countries (
            country_id INTEGER PRIMARY KEY AUTOINCREMENT,
            country_name TEXT NOT NULL,
            country_code TEXT NOT NULL
        )
    `)
	return err
}

func (s *SQLStorage) createServersAuthorizationsTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS servers_authorizations (
            server_id INTEGER PRIMARY KEY,
            username TEXT NOT NULL,
            password TEXT NOT NULL,
            token TEXT NOT NULL,
            expired_at TIMESTAMP NOT NULL,
            meta TEXT NOT NULL
        )
    `)
	return err
}
