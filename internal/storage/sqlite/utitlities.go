package sqlite

import (
	"context"
	"fmt"
	"log"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
)

func (s *SQLStorage) Count(ctx context.Context, args *storage.QueryArgs) (n int64, err error) {
	q := `
		SELECT count(*) AS count 
	`

	var queryEnd string
	var queryArgs []interface{}
	if args != nil {
		selectArgs := s.parseQueryArgs(args)
		queryEnd, queryArgs = s.builder.BuildParts([]string{"from", "where"}, selectArgs)
		q += queryEnd
	}

	log.Printf("query: `%s` args: %+v", q, queryArgs)
	count := &[]int64{}
	err = s.db.SelectContext(ctx, count, q, queryArgs...)
	if err != nil {
		return -1, err
	}

	fmt.Printf("got count %+v", count)

	return (*count)[0], nil
}

func (s *SQLStorage) parseQueryArgs(args *storage.QueryArgs) (a *builder.SelectArguments) {
	selectArgs := builder.SelectArguments{
		From: builder.Table(args.From),
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
