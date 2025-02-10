package sqlite

func (s *SQLStorage) createUsersTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            telegram_id TEXT NOT NULL,
            telegram_name TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
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
            protocol TEXT,
            host TEXT,
            port INTEGER,
            username TEXT,
            password TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        )
    `)
	return err
}

func (s *SQLStorage) createSubscriptionsTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS subscriptions (
            user_id INTEGER,
            server_id INTEGER,
            subscription_status TEXT, /* active | inactive */
            subscription_expired_at TIMESTAMP,
            UNIQUE(user_id, server_id)
        )
    `)
	return err
}

func (s *SQLStorage) createCountriesTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS countries (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            icon TEXT
        )
    `)
	return err
}

func (s *SQLStorage) createServersAuthorizationsTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS servers_authorizations (
            server_id INTEGER PRIMARY KEY,
            username TEXT,
            password TEXT,
            token TEXT,
            expired_at TIMESTAMP,
            meta TEXT
        )
    `)
	return err
}
