package sqlite

import (
	"context"
	"database/sql"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"
)

func (s *SQLStorage) SaveCountry(ctx context.Context, country *storage.Country) (countryID storage.CountryID, err error) {
	var id int64

	q := `SELECT * FROM countries WHERE country_id = ? OR country_code = ? LIMIT 1`

	oldCountry := &storage.Country{}
	err = s.db.GetContext(ctx, oldCountry, q, country.CountryID, country.CountryCode)
	if err == sql.ErrNoRows {
		q = `
			INSERT INTO countries (country_name, country_code)
			VALUES (?,?)
		`
		result, err := s.db.ExecContext(ctx, q, country.CountryName, country.CountryCode)
		if err != nil {
			return 0, e.Wrap("can't execute query", err)
		}

		id, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
		countryID = storage.CountryID(id)
	} else if err != nil {
		return 0, e.Wrap("can't scan row", err)
	} else {
		q = `UPDATE countries SET country_name = ?, country_code = ? WHERE country_id = ?`

		_, err = s.db.ExecContext(ctx, q, country.CountryName, country.CountryCode, oldCountry.CountryID)
		if err != nil {
			return 0, e.Wrap("can't execute query", err)
		}
		countryID = oldCountry.CountryID
	}

	return countryID, nil
}

func (s *SQLStorage) GetCountries(ctx context.Context, args *storage.QueryArgs) (countries *[]storage.Country, err error) {
	q := `SELECT * FROM countries`

	var queryEnd string
	var queryArgs []interface{}
	if args != nil {

		queryEnd, queryArgs = s.buildParts([]string{"where", "order_by", "limit"}, args)
		q += " " + queryEnd
	}

	countries = &[]storage.Country{}
	err = s.db.SelectContext(ctx, countries, q, queryArgs...)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func (s *SQLStorage) GetCountryByID(ctx context.Context, countryID storage.CountryID) (country *storage.Country, err error) {
	q := `SELECT * FROM countries WHERE country_id = ? LIMIT 1`

	country = &storage.Country{}
	err = s.db.GetContext(ctx, country, q, countryID)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSuchCountry
	} else if err != nil {
		return nil, err
	}

	return country, nil
}

func (s *SQLStorage) RemoveCountryByID(ctx context.Context, countryID storage.CountryID) (err error) {
	q := `DELETE FROM countries WHERE country_id = ?`

	_, err = s.db.ExecContext(ctx, q, countryID)

	return err
}
