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

type Builder interface {
	// Build(Arguments) (string, []interface{})
	// BuildParts([]string, Arguments) (string, []interface{})
	SelectBuilder
	InsertBuilder
}

type SQLBuilder struct {
	Driver  string
	Builder Builder
}

// func NewSQLBuilder(driver string, b Builder) (*SQLBuilder, error) {
// 	return &SQLBuilder{Driver: driver, Builder: b}, nil
// }

type Arguments interface {
	BuildPartByName(partName string, b Builder) (q string, args []interface{})
	PartOrder() []string
}

func (b *SQLBuilder) BuildParts(parts []string, a Arguments) (string, []interface{}) {
	if a == nil {
		return "", nil
	}
	queryParts := make([]string, len(parts))
	sqlArgs := make([]interface{}, 0)
	var tempArgs []interface{}
	for _, partPositionedName := range a.PartOrder() {
		for i, partGivenName := range parts {
			if partGivenName == partPositionedName || partGivenName == partPositionedName+"s" {
				queryParts[i], tempArgs = a.BuildPartByName(partGivenName, b.Builder)
				if tempArgs != nil {
					sqlArgs = append(sqlArgs, tempArgs...)
				}
			}
		}
	}

	return strings.Join(queryParts, " "), sqlArgs
}

func (b *SQLBuilder) Build(a Arguments) (string, []interface{}) {
	if a == nil {
		return "", nil
	}
	queryParts := make([]string, len(a.PartOrder()))
	sqlArgs := make([]interface{}, 0)
	tempArgs := make([]interface{}, 0)
	for i, partName := range a.PartOrder() {
		queryParts[i], tempArgs = a.BuildPartByName(partName, b.Builder)
		sqlArgs = append(sqlArgs, tempArgs...)
	}
	return strings.Join(queryParts, " "), sqlArgs
}

/* ---- Basic types ---- */
type Column string

type Table string

type Value string

type Where struct {
	Column   string
	Operator string
	Value    string
}

type Limit struct {
	Offset int64
	Limit  int64
}

type GroupBy string

type OrderBy struct {
	Column string
	Order  string
}

/* ---- Parse ---- */

// func (c *Column) ParseURLValues(args []string, default []string) {

// 	if queryValue, ok := queryArgs[key]; ok {
// 		v, err := strconv.Atoi(queryValue[0])
// 			if err == nil {
// 				args.Limit.Limit = v
// 				continue
// 			}
// 		}
// 	}

// 	args.Limit.Limit = defaultArgs[key][0]

// }
