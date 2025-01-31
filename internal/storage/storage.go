package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

/* ---- Storage Interface ---- */
// Basic Storage interface, thai includes Users, Servers, Subscriptions interfaces.
type Storage interface {
	Users
	Servers
	Subscriptions
}

/* ---- Simple Types ---- */
type TelegramID int64
type UserID int64
type ServerID int64

/* ---- Errors ---- */
var ErrNoSuchUser = errors.New("no such user")
var ErrNoSuchServer = errors.New("no such server")
var ErrNoSuchServerAuth = errors.New("no such server auth")

/* ---- Users Interface ---- */
// Represent VPN User
type User struct {
	ID           UserID     `db:"id"`
	TelegramID   TelegramID `db:"telegram_id"`
	TelegramName string     `db:"telegram_name"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

// Users interface, that includes:
// - Writer: Save and Remove User,
// - Reader: Basic Get User(s) operations,
// - Query: Multi repository (tables) queries for get User.
type Users interface {
	UserWriter
	UserReader
	UserQuery
}

type UserWriter interface {
	SaveUser(ctx context.Context, user *User) (UserID, error)
	RemoveUserByID(ctx context.Context, id UserID) error
}

type UserReader interface {
	GetAllUsers(ctx context.Context) (*[]User, error)
	GetUserByID(ctx context.Context, id UserID) (*User, error)
}

type UserQuery interface {
	GetUserByServerID(ctx context.Context, serverID ServerID) ([]*User, error)
}

/* ---- Servers Interface ---- */
// Represent VPN Server
type VPNServer struct {
	ID        ServerID  `db:"id"`
	CountryID int64     `db:"country_id"`
	Name      string    `db:"name"`
	Protocol  string    `db:"protocol"`
	Host      string    `db:"host"`
	Port      int       `db:"port"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Servers interface, that includes:
// - Writer: Save and Remove Server,
// - Reader: Basic Get Server(s) operations,
// - Query: Multi repository (tables) queries for get Server.
type Servers interface {
	ServerWriter
	ServerReader
	ServerQuery
}

type ServerWriter interface {
	SaveServer(ctx context.Context, server *VPNServer) (ServerID, error)
	RemoveServerByID(ctx context.Context, id ServerID) error
}

type ServerReader interface {
	GetAllServers(ctx context.Context) ([]*VPNServer, error)
	GetServerByID(ctx context.Context, id ServerID) (*VPNServer, error)
}

type ServerQuery interface {
	GetServerByUserID(ctx context.Context, userID UserID) ([]*VPNServer, error)
}

/* ---- Servers Authorization ---- */
// Represent Authorization data for VPN Server API.
type VPNServerAuthorization struct {
	ServerID  ServerID  `db:"server_id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Token     string    `db:"token"`
	ExpiredAt time.Time `db:"expired_at"`
	Meta      string    `db:"meta"`
}

// Server Authorization interface, that includes methods to work with Authorization data.
type ServerAuthorizations interface {
	SaveServerAuth(ctx context.Context, authorization *VPNServerAuthorization) (ServerID, error)
	GetServerAuthByServerID(ctx context.Context, ServerID ServerID) (*VPNServerAuthorization, error)
	GetServerAndAuthByServerID(ctx context.Context, ServerID ServerID) (*VPNServer, *VPNServerAuthorization, error)
	RemoveServerAuthByServerID(ctx context.Context, ServerID ServerID) error
}

/* ---- Subscriptions Interface ---- */
// Represent Subscription of VPN User to VPN Server
type Subscription struct {
	UserID                int64     `db:"user_id"`
	ServerID              int64     `db:"server_id"`
	SubscriptionStatus    string    `db:"subscription_status"`
	SubscriptionExpiredAt time.Time `db:"subscription_expired_at"`
}

// Subscriptions interface, that includes methods to work with Subscription data.
type Subscriptions interface {
	UpdateSubscription(ctx context.Context, subscription *Subscription) error
	RemoveSubscriptionByID(ctx context.Context, userID UserID, serverID ServerID) error
}

/* ---- SQLCompatible ---- */
// SQLCompatible interface, which allows to use full power of SQL Queries, by the cost of breaking encapsulation.
// SQL language can have difference between different databases, so recommend to use *sql.DB.DriverName()
// to identify concrete database driver under the interface.
type SQLCompatible interface {
	SQLStorageInstance() (*sql.DB, error)
}
