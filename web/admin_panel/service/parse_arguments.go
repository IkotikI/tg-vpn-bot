package service

import (
	"net/url"
	"strconv"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
)

func DefaultSelectArgs() url.Values {
	return url.Values{
		"per_page": []string{"10"},
		"page":     []string{"1"},
		"order":    []string{"DESC"},
		"order_by": []string{"created_at"},
	}
}

func ParseSelectQueryArgs(queryArgs url.Values) (argsPtr *builder.SelectArguments) {
	args := builder.SelectArguments{
		Limit:   builder.Limit{Limit: 10, Offset: 0},
		OrderBy: builder.OrderBy{Column: "created_at", Order: "DESC"},
	}

	args.Limit.Limit = ParsePerPage(queryArgs, 10)

	args.Limit.Offset = (ParsePage(queryArgs, 1) - 1) * args.Limit.Limit

	args.OrderBy.Order = ParseOrder(queryArgs, "DESC")

	args.OrderBy.Column = ParseOrderBy(queryArgs, "created_at")

	return &args
}

// func ParseSelectQueryArgs(queryArgs url.Values) (argsPtr *builder.SelectArguments) {
// 	args := builder.SelectArguments{
// 		Limit:   builder.Limit{Limit: 10, Offset: 0},
// 		OrderBy: builder.OrderBy{Column: "created_at", Order: "DESC"},
// 	}

// 	if queryValue, ok := queryArgs["per_page"]; ok {
// 		v, err := strconv.ParseInt(queryValue[0], 10, 64)
// 		if err == nil {
// 			args.Limit.Limit = v
// 		}
// 	}

// 	if queryValue, ok := queryArgs["page"]; ok {
// 		v, err := strconv.ParseInt(queryValue[0], 10, 64)
// 		v = v - 1
// 		if v < 0 {
// 			v = 0
// 		}
// 		if err == nil {
// 			args.Limit.Offset = v * args.Limit.Limit
// 		}
// 	}

// 	if queryValue, ok := queryArgs["order"]; ok {
// 		if queryValue[0] == "DESC" {
// 			args.OrderBy.Order = "DESC"
// 		} else {
// 			args.OrderBy.Order = "ASC"
// 		}
// 	}

// 	if queryValue, ok := queryArgs["order_by"]; ok {
// 		args.OrderBy.Column = queryValue[0]
// 	}

// 	return &args
// }

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
			return min(1, v)
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
