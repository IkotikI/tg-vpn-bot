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

func (b *SQLiteBuilder) BuildSelect(s []builder.Column) (string, []interface{}) {
	if len(s) == 0 {
		return "SELECT *", nil
	}
	str, args := joinWithCommas(s)
	return spf("SELECT %s", str), args
}

func (b *SQLiteBuilder) BuildFrom(f builder.Table) (string, []interface{}) {
	if f == "" {
		return "", nil
	}
	return "FROM ?", []interface{}{f}
}

func (b *SQLiteBuilder) BuildWhere(w []builder.Where) (string, []interface{}) {
	if len(w) == 0 {
		return "", nil
	}
	size := len(w)
	str := make([]string, size)
	args := make([]interface{}, size*2)
	for i, where := range w {
		if where.Column == "" || where.Operator == "" || where.Value == "" {
			size--
			continue
		}
		// Each iteration 2 args will be added
		str[i] = spf("? %s ?", where.Operator)
		args[i*2] = where.Column
		args[i*2+1] = where.Value
	}
	// Reduce slice sizes
	str = str[:size]
	args = args[:size*2]
	return spf("WHERE %s", strings.Join(str, " AND ")), args
}

func (b *SQLiteBuilder) BuildGroupBy(g builder.GroupBy) (string, []interface{}) {
	if g == "" {
		return "", nil
	}
	return "GROUP BY ?", []interface{}{g}
}

func (b *SQLiteBuilder) BuildOrderBy(o builder.OrderBy) (string, []interface{}) {
	if o.Column == "" {
		return "", nil
	}
	var order string
	if o.Order == "" {
		order = "ASC"
	} else {
		order = o.Order
	}
	return "ORDER BY ? ?", []interface{}{o.Column, order}
}

func (b *SQLiteBuilder) BuildLimit(l builder.Limit) (string, []interface{}) {
	if l.Limit <= 0 {
		return "", nil
	}
	limit := "LIMIT ?"
	if l.Offset <= 0 {
		return limit, []interface{}{l.Limit}
	}
	offset := "OFFSET ?"
	return spf("%s %s", limit, offset), []interface{}{l.Limit, l.Offset}
}

func (b *SQLiteBuilder) BuildInsertInto(t builder.Table) (string, []interface{}) {
	if len(t) == 0 {
		return "", nil
	}
	return "INSERT INTO ?", []interface{}{t}
}

func (b *SQLiteBuilder) BuildInsertColumns(c []builder.Column) (string, []interface{}) {
	if len(c) == 0 {
		return "", nil
	}
	str, args := joinWithCommas(c)
	return fmt.Sprintf("(%s)", str), args
}

func (b *SQLiteBuilder) BuildValues(t []builder.Value) (string, []interface{}) {
	if len(t) == 0 {
		return "", nil
	}
	str, args := joinWithCommas(t)
	return fmt.Sprintf("VALUES (%s)", str), args
}

func joinWithCommas[T ~string](s []T) (string, []interface{}) {
	str, _ := strings.CutPrefix(strings.Repeat(",?", len(s)), ",")
	args := make([]interface{}, len(s))
	for i, v := range s {
		args[i] = v
	}
	return str, args
}
