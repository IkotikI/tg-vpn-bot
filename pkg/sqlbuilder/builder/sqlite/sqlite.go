package sqlite_builder

import (
	"fmt"
	"strings"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
)

var spf = fmt.Sprintf

type SQLiteBuilder struct {
}

func NewSQLiteBuilder() *SQLiteBuilder {
	return &SQLiteBuilder{}
}

func (b *SQLiteBuilder) BuildSelect(s []string) string {
	if len(s) == 0 {
		return spf("SELECT *")
	}
	return spf("SELECT %s", strings.Join(s, ", "))
}

func (b *SQLiteBuilder) BuildFrom(f string) string {
	if f == "" {
		return ""
	}
	return spf("FROM %s", f)
}

func (b *SQLiteBuilder) BuildWhere(w []builder.Where) string {
	if len(w) == 0 {
		return ""
	}
	whereStrings := make([]string, len(w))
	for i, where := range w {
		if where.Column == "" || where.Operator == "" || where.Value == "" {
			continue
		}
		whereStrings[i] = spf("%s %s %s", where.Column, where.Operator, where.Value)
	}
	return spf("WHERE %s", strings.Join(whereStrings, " AND "))
}

func (b *SQLiteBuilder) BuildGroupBy(g string) string {
	if g == "" {
		return ""
	}
	return spf("GROUP BY %s", g)
}

func (b *SQLiteBuilder) BuildOrderBy(o builder.OrderBy) string {
	if o.Column == "" {
		return ""
	}
	orderby := spf("ORDER BY %s", o.Column)
	var order string
	if o.Order == "" {
		order = "ASC"
	} else {
		order = o.Order
	}
	return spf("%s %s", orderby, order)
}

func (b *SQLiteBuilder) BuildLimit(l builder.Limit) string {
	if l.Limit <= 0 {
		return ""
	}
	limit := spf("LIMIT %d", l.Limit)
	if l.Offset <= 0 {
		return limit
	}
	offset := spf("OFFSET %d", l.Offset)
	return spf("%s %s", limit, offset)
}
