package sqlbuilder

import (
	"errors"
	"vpn-tg-bot/pkg/sqlbuilder/builder"
	sqlite_builder "vpn-tg-bot/pkg/sqlbuilder/builder/sqlite"
)

var drivers = []string{"sqlite3"}

var ErrNoSuchDriver = errors.New("no such driver")

func Drivers() []string {
	dest := make([]string, len(drivers))
	copy(dest, drivers)
	return dest
}

func NewSQLBuilder(driver string) (*builder.SQLBuilder, error) {
	switch driver {
	case "sqlite3":
		return &builder.SQLBuilder{Driver: driver, Builder: sqlite_builder.NewSQLiteBuilder()}, nil
	default:
		return nil, ErrNoSuchDriver
	}

}
