package builder

import (
	"strings"
)

// var drivers = []string{"sqlite3"}

// var ErrNoSuchDriver = errors.New("no such driver")
// func Drivers() []string {
// 	dest := make([]string, len(drivers))
// 	copy(dest, drivers)
// 	return dest
// }

var partOrder = []string{"select", "from", "where", "group_by", "order_by", "limit"}

type Builder interface {
	BuildSelect([]string) string
	BuildFrom(string) string
	BuildWhere([]Where) string
	BuildGroupBy(string) string
	BuildOrderBy(OrderBy) string
	BuildLimit(Limit) string
}

type SQLBuilder struct {
	Driver  string
	Builder Builder
}

// func NewSQLBuilder(driver string, b Builder) (*SQLBuilder, error) {
// 	return &SQLBuilder{Driver: driver, Builder: b}, nil
// }

type Arguments interface {
	BuildPartByName(partName string, b *SQLBuilder) string
}

type SelectArguments struct {
	Arguments

	Select  []string
	From    string
	Where   []Where
	GroupBy string
	OrderBy OrderBy
	Limit   Limit
}

type Where struct {
	Column   string
	Operator string
	Value    string
}

type Limit struct {
	Offset int64
	Limit  int64
}

type OrderBy struct {
	Column string
	Order  string
}

func (b *SQLBuilder) BuildParts(parts []string, a Arguments) string {
	queryParts := make([]string, len(parts))
	for _, partPositionedName := range partOrder {
		for i, partGivenName := range parts {
			if partGivenName == partPositionedName {
				queryParts[i] = a.BuildPartByName(partGivenName, b)
			}
		}
	}

	return strings.Join(queryParts, " ")
}

func (a *SelectArguments) BuildPartByName(partName string, b *SQLBuilder) string {
	switch partName {
	case "select":
		return b.Builder.BuildSelect(a.Select)
	case "from":
		return b.Builder.BuildFrom(a.From)
	case "where":
		return b.Builder.BuildWhere(a.Where)
	case "group_by":
		return b.Builder.BuildGroupBy(a.GroupBy)
	case "order_by":
		return b.Builder.BuildOrderBy(a.OrderBy)
	case "limit":
		return b.Builder.BuildLimit(a.Limit)
	default:
		return ""
	}
}

func (b *SQLBuilder) BuildAll(a Arguments) string {
	queryParts := make([]string, len(partOrder))
	for i, partName := range partOrder {
		queryParts[i] = a.BuildPartByName(partName, b)
	}
	return strings.Join(queryParts, " ")
}
