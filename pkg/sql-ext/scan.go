package sqlext

import (
	"database/sql"
	"reflect"
)

// ScanRowToStruct scans a single SQL row into a struct of any type.
func ScanRowToStruct[T any](row *sql.Row) (T, error) {
	var item T

	// Get a reference to the struct
	structValue := reflect.ValueOf(&item).Elem()
	numFields := structValue.NumField()

	// Prepare a slice of pointers to struct fields for scanning
	columns := make([]interface{}, numFields)
	for i := 0; i < numFields; i++ {
		field := structValue.Field(i)
		columns[i] = field.Addr().Interface()
	}

	// Scan the single row into the struct
	err := row.Scan(columns...)
	if err != nil {
		return item, err
	}

	return item, nil
}

// ScanRowsToStructSlice scans SQL rows into a slice of structs of any type.
func ScanRowsToStructSlice[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close() // Ensure rows are closed after processing.

	// Result slice
	var results []T

	for rows.Next() {
		// Create a new instance of the target struct
		var item T
		structValue := reflect.ValueOf(&item).Elem()
		numFields := structValue.NumField()

		// Prepare a slice of pointers to struct fields for scanning
		columns := make([]interface{}, numFields)
		for i := 0; i < numFields; i++ {
			field := structValue.Field(i)
			columns[i] = field.Addr().Interface()
		}

		// Scan the current row into the struct
		err := rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		// Append the struct to the result slice
		results = append(results, item)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
