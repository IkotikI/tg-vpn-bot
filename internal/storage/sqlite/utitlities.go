package sqlite

import (
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
)

func (s *SQLStorage) parseQueryArgs(args *storage.QueryArgs) (a *builder.SelectArguments) {
	selectArgs := builder.SelectArguments{
		OrderBy: builder.OrderBy{
			Column: args.OrderBy,
			Order:  string(args.Order),
		},
		Limit: builder.Limit{
			Offset: args.Offset,
			Limit:  args.Limit,
		},
	}
	where := make([]builder.Where, 0, len(args.Where))
	for _, w := range args.Where {
		where = append(selectArgs.Where, builder.Where{
			Column:   w.Column,
			Operator: string(w.Operator),
			Value:    w.Value,
		})
	}
	selectArgs.Where = where

	return &selectArgs
}
