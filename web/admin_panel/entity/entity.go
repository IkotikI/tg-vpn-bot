package entity

import "vpn-tg-bot/internal/storage"

var TimeLayout = "2006-01-02 15:04:05"

type UserWithSubscription struct {
	storage.User
	storage.Subscription
}

type ServerWithAuthorization struct {
	storage.VPNServer
	storage.VPNServerAuthorization
	storage.Country
}

type User struct {
	storage.User
}

type Server struct {
	storage.VPNServer
	storage.Country
}

type Subscription struct {
	storage.Subscription
}

type Country struct {
	storage.Country
}

type SubscriptionWithServer struct {
	Server
	Subscription
}

type SubscriptionWithUser struct {
	User
	Subscription
}

type SubscriptionWithUserAndServer struct {
	User
	Server
	Subscription
}

type QueryArguments struct {
	Search  string
	Page    int64
	PerPage int64
	Order   string
	OrderBy string
}

var DefaultQueryArguments = &QueryArguments{Search: "", Page: 1, PerPage: 10, Order: "DESC", OrderBy: "created_at"}

type Pagination struct {
	Table        storage.Table
	RecordsCount int64
	TotalPages   int64
	Page         int64
	PerPage      int64
}

type PaginationLink struct {
	Link string
	Num  int64
}
