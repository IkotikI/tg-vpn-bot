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

	if queryValue, ok := queryArgs["per_page"]; ok {
		v, err := strconv.ParseInt(queryValue[0], 10, 64)
		if err == nil {
			args.Limit.Limit = v
		}
	}

	if queryValue, ok := queryArgs["page"]; ok {
		v, err := strconv.ParseInt(queryValue[0], 10, 64)
		v = v - 1
		if v < 0 {
			v = 0
		}
		if err == nil {
			args.Limit.Offset = v * args.Limit.Limit
		}
	}

	if queryValue, ok := queryArgs["order"]; ok {
		if queryValue[0] == "DESC" {
			args.OrderBy.Order = "DESC"
		} else {
			args.OrderBy.Order = "ASC"
		}
	}

	if queryValue, ok := queryArgs["order_by"]; ok {
		args.OrderBy.Column = queryValue[0]
	}

	return &args
}
