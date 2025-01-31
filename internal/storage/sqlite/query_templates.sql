
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id TEXT NOT NULL,
    telegram_name TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)

CREATE TABLE IF NOT EXISTS vpn_servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    country_id INTEGER, /* ? */
    name TEXT NOT NULL,
    protocol TEXT,
    ip_address TEXT,
    port INTEGER,
    login TEXT,
    password TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)

CREATE TABLE IF NOT EXISTS users_vpn_servers (
    user_id INTEGER,
    server_id INTEGER,
    subscription_status TEXT, /* active | inactive */
    subscription_expired_at TIMESTAMP
)

CREATE TABLE IF NOT EXISTS countries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    icon TEXT
)

CREATE TABLE IF NOT EXISTS vpn_servers_authorization (
    server_id INTEGER PRIMARY KEY,
    username TEXT,
    password TEXT,
    token TEXT,
    expired_at TIMESTAMP,
    meta TEXT
)