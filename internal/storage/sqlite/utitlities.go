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
	queryEnd, queryArgs = s.buildParts([]string{"from", "where"}, args)
	q += queryEnd

	log.Printf("query: `%s` args: %+v\n", q, queryArgs)
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

func (s *SQLStorage) buildParts(parts []string, args *storage.QueryArgs) (query string, queryArgs []interface{}) {
	if args == nil {
		return "", []interface{}{}
	}
	bulderArgs := s.parseQueryArgs(args)
	return s.builder.BuildParts(parts, bulderArgs)
}

func (s *SQLStorage) MakePagination(ctx context.Context, db storage.Storage, table storage.Table, queryArgs *storage.QueryArguments, args *storage.QueryArgs) (storage.Pagination, error) {
	if queryArgs == nil {
		queryArgs = storage.DefaultQueryArguments
	}
	if args == nil {
		args = queryArgs.ToQueryArgs()
	}
	args.From = storage.Table(table)

	n, err := db.CountWithBuilder(ctx, args)
	if err != nil {
		return storage.Pagination{}, err
	}
	total_pages := n / queryArgs.PerPage
	if n-total_pages > 0 {
		total_pages += 1
	}
	return storage.Pagination{
		Table:        table,
		RecordsCount: n,
		TotalPages:   total_pages,
		Page:         queryArgs.Page,
		PerPage:      queryArgs.PerPage,
	}, nil
}
