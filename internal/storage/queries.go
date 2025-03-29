package storage

import (
	"context"
)

var TimeLayout = "2006-01-02 15:04:05"

type UserWithSubscription struct {
	User
	Subscription
}

type ServerWithAuthorization struct {
	VPNServer
	VPNServerAuthorization
	Country
}

type VPNServerWithCountry struct {
	VPNServer
	Country
}

type SubscriptionWithServer struct {
	VPNServerWithCountry
	Subscription
}

type SubscriptionWithUser struct {
	User
	Subscription
}

type SubscriptionWithUserAndServer struct {
	User
	VPNServerWithCountry
	Subscription
}

/* ---- Queries Interface ---- */
type Queries interface {
	GetServerWithCountryByID(ctx context.Context, serverID ServerID) (*VPNServerWithCountry, error)
	GetServersWithCountries(ctx context.Context, args *QueryArgs) (*[]VPNServerWithCountry, error)
	MakePagination(ctx context.Context, db Storage, table Table, queryArgs *QueryArguments, args *QueryArgs) (Pagination, error)
	GetSubscriptionsWithServersByUserID(ctx context.Context, user_id UserID, args *QueryArgs) (*[]SubscriptionWithServer, error)
	GetSubscriptionsWithUsersByServerID(ctx context.Context, user_id ServerID, args *QueryArgs) (*[]SubscriptionWithUser, error)
	GetSubscriptionsWithUsersAndServers(ctx context.Context, args *QueryArgs) (*[]SubscriptionWithUserAndServer, error)
	GetSubscriptionWithUserAndServerByIDs(ctx context.Context, userID UserID, serverID ServerID) (*SubscriptionWithUserAndServer, error)
	CountWithBuilder(ctx context.Context, args *QueryArgs) (n int64, err error)
}

/* ---- Pagination queries ---- */
type QueryArguments struct {
	Search  string
	Page    int64
	PerPage int64
	Order   string
	OrderBy string
}

var DefaultQueryArguments = &QueryArguments{Search: "", Page: 1, PerPage: 10, Order: "DESC", OrderBy: "created_at"}

type Pagination struct {
	Table        Table
	RecordsCount int64
	TotalPages   int64
	Page         int64
	PerPage      int64
}

type PaginationLink struct {
	Link string
	Num  int64
}
