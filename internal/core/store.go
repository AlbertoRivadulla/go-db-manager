package core

import (
    "fmt"

    "database/sql"
    _ "modernc.org/sqlite"
)

type SQLiteStore struct {
    db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }

    return &SQLiteStore{db: db}, nil
}

func (store *SQLiteStore) RunQueryNoReturn(query string) (sql.Result, error) {
    result, err := store.db.Exec(query)
    if err != nil {
        return nil, fmt.Errorf("running query %s: %w", query, err)
    }

    return result, nil
}

func (store *SQLiteStore) RunQuery(query string) (map[string]any, error) {
    rows, err := store.db.Query(
        query,
    )
    if err != nil {
        return nil, fmt.Errorf("running query %s: %w", query, err)
    }
    defer rows.Close()

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }

    // Create a slice of any to hold the values
    values := make([]any, len(columns))
    valuePtrs := make([]any, len(columns))
    for i := range values {
        valuePtrs[i] = &values[i]
    }

    // Get the first row
    if !rows.Next() {
        return nil, nil
        // return nil, sql.ErrNoRows
    }

    // Scan into the value pointers
    if err := rows.Scan(valuePtrs...); err != nil {
        return nil, err
    }

    // Create map with column names as keys
    result := make(map[string]any)
    for i, col := range columns {
        result[col] = values[i]
    }

    return result, nil
}

func (store *SQLiteStore) QueryToMap(query string, args ...any) ([]map[string]any, error) {
    rows, err := store.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }

    var results []map[string]any

    for rows.Next() {
        // Create a slice of any for this row
        values := make([]any, len(columns))
        valuePtrs := make([]any, len(columns))
        for i := range values {
            valuePtrs[i] = &values[i]
        }

        if err := rows.Scan(valuePtrs...); err != nil {
            return nil, err
        }

        // Create map for this row
        rowMap := make(map[string]any)
        for i, col := range columns {
            val := values[i]

            // Convert []byte to string for better usability
            if b, ok := val.([]byte); ok {
                rowMap[col] = string(b)
            } else {
                rowMap[col] = val
            }
        }

        results = append(results, rowMap)
    }

    return results, rows.Err()
}

func (store *SQLiteStore) QueryTo2DimArray(query string, args ...any) ([][]string, error) {
    rows, err := store.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }

    var results [][]string

    for rows.Next() {
        // Create a slice of any for this row
        values := make([]any, len(columns))
        valuePtrs := make([]any, len(columns))
        for i := range values {
            valuePtrs[i] = &values[i]
        }

        if err := rows.Scan(valuePtrs...); err != nil {
            return nil, err
        }

        // Convert to []string
        row := make([]string, len(columns))
        for i, val := range values {
            if val == nil {
                row[i] = "NULL"
            } else {
                row[i] = fmt.Sprintf("%v", val)
            }
        }

        results = append(results, row)
    }

    return results, rows.Err()
}
