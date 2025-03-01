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
}

type UserQuery interface {
	GetUsers(ctx context.Context, args *QueryArgs) (*[]User, error)
	GetUserByServerID(ctx context.Context, serverID ServerID) (*[]User, error)
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
	GetServerByUserID(ctx context.Context, userID UserID) (*[]VPNServer, error)
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

// type SQLStorage struct {
// 	db *sql.DB
// }

// func NewSQLStorage(db *sql.DB) *SQLStorage {
// 	return &SQLStorage{db: db}
// }

// func (s *SQLStorage) SQLStorageInstance() (*sql.DB, error) {
// 	return s.db, nil
// }

type SQLStorage interface {
	Storage
	SQLCompatible
}

/* --- Tables ---- */
type Table string

const (
	TableUsers               Table = "users"
	TableServer              Table = "servers"
	TableSubscription        Table = "subscriptions"
	TableCountry             Table = "countries"
	TableServerAuthorization Table = "server_authorizations"
)

/* --- Query --- */
type QueryArgs struct {
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

/* ---- Builder Interface ---- */
// For making complex requests need to provide a Builder.
// Builder provide methods:
// - Build: build full requests;
// - BuildParts: build parts of requests, that would joined in PartsOrder of BuilderArguments.
// Both returns string query and slice of arguments.
//
// See: pkg/sqlbuilder
//
// Example Implementation:
// type ConcreteBuilder struct {}
// type ConcreteArguments struct {
// 	Arg1 any
// 	Arg2 any
// }
// func (b *ConcreteBuilder) Build(args interface{}) (query string, queryArgs []interface{}) {
// 	for _, partName := range args.PartsOrder() {
// 		queryPart, partArgs := args.BuildPartByName(partName, b)
// 		query += queryPart
// 		queryArgs = append(queryArgs, partArgs...)
// 	}
// 	return
// }
// // Like Build, but filter parts from give slice and join them in parts order.
// func (b *ConcreteBuilder) BuildParts(parts []string, args interface{}) (query string, queryArgs []interface{}) {...}
// func (a *ConcreteArguments) BuildPartByName(partName string, b Builder) (queryPart string, partArgs []interface{}) {
// 	switch partName {
// 	case "partName1":
// 		return b.BuildPart1(a.Arg1)
// 	case "partName2":
// 		return b.BuildPart2(a.Arg2)
// 	}
// }
// func (a *ConcreteArguments) PartsOrder() []string { return []string{"partName1", "partName2"} }
//
type Builder interface {
	Build(args interface{}) (query string, queryArgs []interface{})
	BuildParts(parts []string, args interface{}) (query string, queryArgs []interface{})
	ValidateArgs(args interface{}) error
}

// Builder Arguments provide methods:
//   - BuildPartByName: build part of request by string name,
//     returns string query part and slice of arguments;
//   - PartsOrder: should provide Names of available parts in correct order.
// type BuilderArguments interface {
// 	BuildPartByName(partName string, b Builder) (queryPart string, partArgs []interface{})
// 	PartsOrder() []string
// }
