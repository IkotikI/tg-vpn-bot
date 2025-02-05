package sqlbuilder_test

import (
	"testing"
	"vpn-tg-bot/pkg/sqlbuilder"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
)

func TestSQLite(t *testing.T) {
	args := &builder.SelectArguments{
		Select:  []string{"id", "name"},
		From:    "users",
		Where:   []builder.Where{{Column: "id", Operator: "=", Value: "5"}},
		GroupBy: "id",
		OrderBy: builder.OrderBy{Column: "id", Order: "DESC"},
		Limit:   builder.Limit{Offset: 10, Limit: 10},
	}

	builder, err := sqlbuilder.NewSQLBuilder("sqlite3")
	if err != nil {
		t.Fatalf("Can't create sqlite builder: %v", err)
	}

	q := builder.BuildParts([]string{"select", "from", "where", "group_by", "order_by", "limit"}, args)

	t.Logf("q: %s", q)
}
