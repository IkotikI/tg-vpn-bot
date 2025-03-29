package storage

import (
	"net/url"
	"strconv"
	"time"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
)

func (u *User) ParseURLValues(url url.Values, time_layout string) {
	id, err := strconv.ParseInt(url.Get("id"), 10, 64)
	if err != nil {
		id = 0
	}
	u.ID = UserID(id)

	telegram_id, err := strconv.ParseInt(url.Get("telegram_id"), 10, 64)
	if err != nil {
		telegram_id = 0
	}
	u.TelegramID = TelegramID(telegram_id)

	u.TelegramName = url.Get("telegram_name")

	created_at, err := time.Parse(time_layout, url.Get("created_at"))
	if err != nil {
		created_at = time.Time{}
	}
	u.CreatedAt = created_at

	updated_at, err := time.Parse(time_layout, url.Get("updated_at"))
	if err != nil {
		updated_at = time.Time{}
	}
	u.UpdatedAt = updated_at
}

func (u *User) ParseDefaultsFrom(from *User) {
	if u.ID == 0 {
		u.ID = from.ID
	}
	if u.TelegramID == 0 {
		u.TelegramID = from.TelegramID
	}
	if u.TelegramName == "" {
		u.TelegramName = from.TelegramName
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = from.CreatedAt
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = from.UpdatedAt
	}

}

func (s *VPNServer) ParseURLValues(url url.Values, time_layout string) {
	id, err := strconv.ParseInt(url.Get("id"), 10, 64)
	if err != nil {
		id = 0
	}
	s.ID = ServerID(id)

	country_id, err := strconv.ParseInt(url.Get("country_id"), 10, 64)
	if err != nil {
		country_id = 0
	}
	s.CountryID = CountryID(country_id)

	s.Name = url.Get("name")

	s.Protocol = url.Get("protocol")

	s.Host = url.Get("host")

	port, err := strconv.Atoi(url.Get("port"))
	if err != nil {
		port = 0
	}
	s.Port = port

	s.Username = url.Get("username")

	s.Password = url.Get("password")

	created_at, err := time.Parse(time_layout, url.Get("created_at"))
	if err != nil {
		created_at = time.Time{}
	}
	s.CreatedAt = created_at

	updated_at, err := time.Parse(time_layout, url.Get("updated_at"))
	if err != nil {
		updated_at = time.Time{}
	}
	s.UpdatedAt = updated_at
}

func (s *VPNServer) ParseDefaultsFrom(from *VPNServer) {
	if s.ID == 0 {
		s.ID = from.ID
	}
	if s.CountryID == 0 {
		s.CountryID = from.CountryID
	}
	if s.Protocol == "" {
		s.Protocol = from.Protocol
	}
	if s.Host == "" {
		s.Host = from.Host
	}
	if s.Port == 0 {
		s.Port = from.Port
	}
	if s.Username == "" {
		s.Username = from.Username
	}
	if s.Password == "" {
		s.Password = from.Password
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = from.CreatedAt
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = from.UpdatedAt
	}
}

func (s *Subscription) ParseDefaultsFrom(from *Subscription) {
	if s.UserID == 0 {
		s.UserID = from.UserID
	}
	if s.ServerID == 0 {
		s.ServerID = from.ServerID
	}
	if s.SubscriptionStatus == "" {
		s.SubscriptionStatus = from.SubscriptionStatus
	}
	if s.SubscriptionExpiredAt.IsZero() {
		s.SubscriptionExpiredAt = from.SubscriptionExpiredAt
	}
}

func (s *Subscription) ParseURLValues(url url.Values, time_layout string) {
	user_id, err := strconv.ParseInt(url.Get("user_id"), 10, 64)
	if err != nil {
		user_id = 0
	}
	s.UserID = UserID(user_id)

	server_id, err := strconv.ParseInt(url.Get("server_id"), 10, 64)
	if err != nil {
		server_id = 0
	}
	s.ServerID = ServerID(server_id)

	s.SubscriptionStatus = url.Get("subscription_status")

	s.SubscriptionExpiredAt, err = time.Parse(time_layout, url.Get("subscription_expired_at"))
	if err != nil {
		s.SubscriptionExpiredAt = time.Time{}
	}
}

func (s *Subscription) IsCorrectStatus() bool {
	switch s.SubscriptionStatus {
	case SubscriptionStatusActive, SubscriptionStatusExpired, SubscriptionStatusDisabled:
		return true
	}
	return false
}

func (s *VPNServerWithCountry) ParseDefaultsFrom(from *VPNServerWithCountry) {
	// storage.VPNServer
	if s.ID == 0 {
		s.ID = from.ID
	}
	if s.VPNServer.CountryID == 0 {
		s.VPNServer.CountryID = from.VPNServer.CountryID
	}
	if s.Protocol == "" {
		s.Protocol = from.Protocol
	}
	if s.Host == "" {
		s.Host = from.Host
	}
	if s.Port == 0 {
		s.Port = from.Port
	}
	if s.Username == "" {
		s.Username = from.Username
	}
	if s.Password == "" {
		s.Password = from.Password
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = from.CreatedAt
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = from.UpdatedAt
	}
	// storage.Country
	s.Country.CountryID = s.VPNServer.CountryID // Unify fields
	if s.CountryName == "" {
		s.CountryName = from.CountryName
	}
	if s.CountryCode == "" {
		s.CountryCode = from.CountryCode
	}
}

// func (p *Pagination) ParseURLValues(url url.Values, time_layout string) {
// 	page, err := strconv.Atoi(url.Get("page"))
// 	if err != nil || page < 1 {
// 		page = 0
// 	}
// 	p.Page = page

// 	per_page, err := strconv.Atoi(url.Get("per_page"))
// 	if err != nil || per_page < 1 {
// 		per_page = 0
// 	}
// 	p.PerPage = per_page
// }

func (a *QueryArguments) ParseDefaultsFrom(from *QueryArguments) {
	if a.Search == "" {
		a.Search = from.Search
	}
	if a.Page == 0 {
		a.Page = from.Page
	}
	if a.PerPage == 0 {
		a.PerPage = from.PerPage
	}
	if a.Order == "" {
		a.Order = from.Order
	}
	if a.OrderBy == "" {
		a.OrderBy = from.OrderBy
	}
}

func (a *QueryArguments) ParseURLValues(url url.Values, time_layout string) {
	a.Search = url.Get("search")
	a.Page = ParsePage(url, DefaultQueryArguments.Page)
	a.PerPage = ParsePerPage(url, DefaultQueryArguments.PerPage)
	a.Order = ParseOrder(url, DefaultQueryArguments.Order)
	a.OrderBy = ParseOrderBy(url, DefaultQueryArguments.OrderBy)
}

func DefaultArgs() url.Values {
	return url.Values{
		"per_page": []string{"10"},
		"page":     []string{"1"},
		"order":    []string{"DESC"},
		"order_by": []string{"created_at"},
	}
}

func ParseSelectQueryArgs(queryArgs url.Values) (argsPtr *QueryArgs) {
	limit := ParsePerPage(queryArgs, 10)
	args := QueryArgs{
		Limit:   limit,
		Offset:  (ParsePage(queryArgs, 1) - 1) * limit,
		Order:   Order(ParseOrder(queryArgs, "DESC")),
		OrderBy: ParseOrderBy(queryArgs, "created_at"),
	}

	return &args
}

func ParsePerPage(queryArgs url.Values, def int64) int64 {
	if queryValue, ok := queryArgs["per_page"]; ok {
		v, err := strconv.ParseInt(queryValue[0], 10, 64)
		if err == nil {
			return v
		}
	}
	return def
}

func ParsePage(queryArgs url.Values, def int64) int64 {
	if queryValue, ok := queryArgs["page"]; ok {
		v, err := strconv.ParseInt(queryValue[0], 10, 64)
		if err == nil {
			return max(1, v)
		}
	}
	return def
}

func ParseOrder(queryArgs url.Values, def string) string {
	if queryValue, ok := queryArgs["order"]; ok {
		if queryValue[0] == "DESC" {
			return "DESC"
		} else {
			return "ASC"
		}
	}
	return def
}

func ParseOrderBy(queryArgs url.Values, def string) string {
	if queryValue, ok := queryArgs["order_by"]; ok {
		return queryValue[0]
	}
	return def
}

func ParseQueryArgs(queryArgs url.Values) (args *QueryArguments) {
	args = &QueryArguments{}
	args.ParseURLValues(queryArgs, TimeLayout)
	return args
}

func (a *QueryArguments) ToSelectArgs() (argsPtr *builder.SelectArguments) {
	args := builder.SelectArguments{}

	args.Limit.Limit = a.PerPage

	args.Limit.Offset = (a.Page - 1) * args.Limit.Limit

	args.OrderBy.Order = a.Order

	args.OrderBy.Column = a.OrderBy

	return &args
}

func (a *QueryArguments) ToQueryArgs() (argsPtr *QueryArgs) {
	args := QueryArgs{}

	args.Limit = a.PerPage

	args.Offset = (a.Page - 1) * args.Limit

	if a.Order == string(OrderDECS) {
		args.Order = OrderDECS
	} else {
		args.Order = OrderASC
	}

	args.OrderBy = a.OrderBy

	return &args
}
