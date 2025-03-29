package storage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"
)

var ErrNeedSQLCompatibleInterface = errors.New("SQL storage must be SQLCompatible")
var ErrDatabaseDriverIsNotSupported = errors.New("This database driver is not supported")

/* ---- Storage Interface ---- */
// Basic Storage interface, thai includes Users, Servers, Subscriptions interfaces.
type Storage interface {
	Users
	Servers
	Subscriptions
	Countries
	Utilities
	Queries
}

/* ---- Simple Types ---- */
type TelegramID int64
type UserID int64
type ServerID int64
type CountryID int64

func (id TelegramID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func (id UserID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func (id ServerID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func (id CountryID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

/* ---- Errors ---- */
var ErrNoSuchUser = errors.New("no such user")
var ErrNoSuchServer = errors.New("no such server")
var ErrNoSuchServerAuth = errors.New("no such server auth")
var ErrNoSuchSubscription = errors.New("no such subscription")
var ErrNoSuchCountry = errors.New("no such country")

var ErrZeroUserID = errors.New("zero user id")
var ErrZeroServerID = errors.New("zero server id")
var ErrZeroTelegramID = errors.New("zero telegram id")

/* ---- Utilities Interface ---- */
// Utilities interface, that includes:
// - Count: Imlement counting records in table. Accept "where", "from" QueryArgs.
type Utilities interface {
	Count(ctx context.Context, args *QueryArgs) (int64, error)
}

/* ---- Users Interface ---- */
// Represent VPN User
type User struct {
	ID           UserID     `db:"id" json:"id"`
	TelegramID   TelegramID `db:"telegram_id" json:"telegram_id"`
	TelegramName string     `db:"telegram_name" json:"telegram_name"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
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
	GetUserByID(ctx context.Context, id UserID) (*User, error)
	GetUserByTelegramID(ctx context.Context, telegramID TelegramID) (*User, error)
}

type UserQuery interface {
	GetUsers(ctx context.Context, args *QueryArgs) (*[]User, error)
	GetUsersByServerID(ctx context.Context, serverID ServerID) (*[]User, error)
}

/* ---- Servers Interface ---- */
// Represent VPN Server
type VPNServer struct {
	ID        ServerID  `db:"id" json:"id"`
	CountryID CountryID `db:"country_id" json:"country_id"`
	Name      string    `db:"name" json:"name"`
	Protocol  string    `db:"protocol" json:"protocol"`
	Host      string    `db:"host" json:"host"`
	Port      int       `db:"port" json:"port"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
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
	GetServerByID(ctx context.Context, id ServerID) (*VPNServer, error)
}

type ServerQuery interface {
	GetServers(ctx context.Context, args *QueryArgs) (*[]VPNServer, error)
	GetServersByUserID(ctx context.Context, userID UserID) (*[]VPNServer, error)
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

const (
	SubscriptionStatusActive   = "active"
	SubscriptionStatusExpired  = "expired"
	SubscriptionStatusDisabled = "disabled"
)

// Represent Subscription of VPN User to VPN Server
type Subscription struct {
	UserID                UserID    `db:"user_id" json:"user_id"`
	ServerID              ServerID  `db:"server_id" json:"server_id"`
	SubscriptionStatus    string    `db:"subscription_status" json:"subscription_status"`
	SubscriptionExpiredAt time.Time `db:"subscription_expired_at" json:"subscription_expired_at"`
}

// Subscriptions interface, that includes methods to work with Subscription data.
type Subscriptions interface {
	SaveSubscription(ctx context.Context, subscription *Subscription) error
	RemoveSubscriptionByIDs(ctx context.Context, userID UserID, serverID ServerID) error
	GetSubscriptionsByUserID(ctx context.Context, userID UserID) (*[]Subscription, error)
	GetSubscriptionsServerID(ctx context.Context, serverID ServerID) (*[]Subscription, error)
	GetSubscriptionByIDs(ctx context.Context, userID UserID, serverID ServerID) (*Subscription, error)
}

/* ---- Countries Interface ---- */
type Countries interface {
	SaveCountry(ctx context.Context, country *Country) (CountryID, error)
	GetCountries(ctx context.Context, args *QueryArgs) (*[]Country, error)
	GetCountryByID(ctx context.Context, id CountryID) (*Country, error)
	RemoveCountryByID(ctx context.Context, id CountryID) error
}

// Represent Country
type Country struct {
	CountryID   CountryID `db:"country_id" json:"-"`
	CountryName string    `db:"country_name" json:"name"`
	CountryCode string    `db:"country_code" json:"code"`
}

/* ---- SQLCompatible ---- */
// SQLCompatible interface, which allows to use full power of SQL Queries, by the cost of breaking encapsulation.
// SQL language can have difference between different databases, so recommend to use *sql.DB.DriverName()
// to identify concrete database driver under the interface.
type SQLCompatible interface {
	SQLStorageInstance() (*sql.DB, error)
}

type SQLStorage interface {
	Storage
	SQLCompatible
}

/* ---- Tables ---- */
type Table string

const (
	TableUsers                Table = "users"
	TableServers              Table = "servers"
	TableSubscriptions        Table = "subscriptions"
	TableCountries            Table = "countries"
	TableServerAuthorizations Table = "server_authorizations"
)

/* ---- Query ---- */
// Provide abstract arguments for making SQL queries.
// Concrete implementation lies on chosen
type QueryArgs struct {
	From    Table
	Where   []Where
	Order   Order
	OrderBy string
	Limit   int64
	Offset  int64
}

type Operator string

const (
	OpEqual       Operator = "="
	OpNotEqual    Operator = "!="
	OpLess        Operator = "<"
	OpMore        Operator = ">"
	OpLessOrEqual Operator = "<="
	OpMoreOrEqual Operator = ">="
)

type Where struct {
	Column   string
	Operator Operator
	Value    interface{}
}

type Order string

const (
	OrderASC  Order = "ASC"
	OrderDECS Order = "DESC"
)
