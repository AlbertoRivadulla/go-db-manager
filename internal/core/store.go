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

func (store *SQLiteStore) RunQuery(query string) (map[string]interface{}, error) {
	fmt.Printf("\nRunning query: %s\n", query)

	// result, err := store.db.Exec(
	rows, err := store.db.Query(
		query,
    )
	defer rows.Close()
    if err != nil {
		fmt.Printf("Error running query: %w\n", err)
        return nil, err
    }
    
    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }

	// Create a slice of interface{} to hold the values
    values := make([]interface{}, len(columns))
    valuePtrs := make([]interface{}, len(columns))
    for i := range values {
        valuePtrs[i] = &values[i]
    }

    // Get the first row
    if !rows.Next() {
        return nil, sql.ErrNoRows
    }

    // Scan into the value pointers
    if err := rows.Scan(valuePtrs...); err != nil {
        return nil, err
    }

    // Create map with column names as keys
    result := make(map[string]interface{})
    for i, col := range columns {
        result[col] = values[i]
    }

    return result, nil
}

func (store *SQLiteStore) QueryToMap(query string, args ...interface{}) ([]map[string]interface{}, error) {
	fmt.Printf("\nRunning query: %s\n", query)
    rows, err := store.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    
    var results []map[string]interface{}
    
    for rows.Next() {
        // Create a slice of interface{} for this row
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range values {
            valuePtrs[i] = &values[i]
        }
        
        if err := rows.Scan(valuePtrs...); err != nil {
            return nil, err
        }
        
        // Create map for this row
        rowMap := make(map[string]interface{})
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


