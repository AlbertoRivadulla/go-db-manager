package core

import (
	"fmt"

	// _ "modernc.org/sqlite"
)

type Backend struct {
	store SQLiteStore
	// TODO: different specs of the database
	// TODO: queries
}

func NewBackend(dbPath *string, specsDir *string) (*Backend, error) {

	// TODO: Sanitize the database and specs paths

	store, err := NewSQLiteStore(*dbPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening the SQLite store: %w", err)
	}
	fmt.Printf("Created database\n")


	// // NOTE: Example of adding data to the database
	// _, err = store.RunQuery(`
	//        CREATE TABLE IF NOT EXISTS entries (
	//            title TEXT NOT NULL,
	//            content TEXT NOT NULL
	//        );
	//    `)
	// if err != nil {
	// 	fmt.Printf("Error %w\n", err)
	// }
	// _, err = store.RunQuery(`
	// 	INSERT INTO entries VALUES ("my title", "my content")
	// `)
	// if err != nil {
	// 	fmt.Printf("Error %w\n", err)
	// }
	// _, err = store.RunQuery(`
	// 	INSERT INTO entries VALUES ("my title 2", "my content 2")
	// `)
	// if err != nil {
	// 	fmt.Printf("Error %w\n", err)
	// }
	// _, err = store.RunQuery(`
	// 	INSERT INTO entries VALUES ("my title 3", "my content 3")
	// `)
	// if err != nil {
	// 	fmt.Printf("Error %w\n", err)
	// }
	// result, err := store.RunQuery(`
	// 	SELECT * FROM entries
	// `)
	// if err != nil {
	// 	fmt.Printf("Error %w\n", err)
	// }
	//
	// fmt.Printf("Result:\n")
	// for k, v := range result {
	// 	fmt.Printf("\t%s : %s\n", k, v)
	//    }
	//
	// resultMapList, err := store.QueryToMap(`
	// 	SELECT * FROM entries
	// `)
	// if err != nil {
	// 	fmt.Printf("Error %w\n", err)
	// }
	//
	// fmt.Printf("Result map:\n")
	// for key := range(resultMapList[0]) {
	// 	fmt.Printf("\t%s:", key)
	// 	for _, row := range resultMapList{
	// 		fmt.Printf("\t%s", row[key])
	// 	}
	// 	fmt.Println("")
	// }



	// TODO: Load specs

	return &Backend{
		store: *store,
	}, nil
}

func (backend *Backend) Cleanup() {
	fmt.Printf("Closing database\n")
	backend.store.db.Close()
}

