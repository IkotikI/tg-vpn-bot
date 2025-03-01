package entity

import (
	"net/url"
	"strconv"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
)

func (s *Server) ParseDefaultsFrom(from *Server) {
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

func (a *QueryArguments) SelectArgs() (argsPtr *builder.SelectArguments) {
	args := builder.SelectArguments{}

	args.Limit.Limit = a.PerPage

	args.Limit.Offset = (a.Page - 1) * args.Limit.Limit

	args.OrderBy.Order = a.Order

	args.OrderBy.Column = a.OrderBy

	return &args
}
