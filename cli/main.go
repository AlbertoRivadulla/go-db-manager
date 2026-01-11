package main

import (
	// "bufio"
	"fmt"
	// "log"
	"os"
	// "strings"

	"carmaintenance/internal/core"
)

func main() {

	// TODO: Remove this
	// reader := bufio.NewReader(os.Stdin)

	if len(os.Args) < 2 {
		fmt.Println("Usage:\n\tcli [add|list|get] [args...]")
		os.Exit(0)
	}

	store, err := core.NewSQLiteStore("./data.db")
	if err != nil {
		fmt.Printf("Error opening the SQLite store: %w\n", err)
		os.Exit(1)
	}

	if len(os.Args) >= 2 {
		command := os.Args[1]

		switch command {
		case "add":
			// TODO:
			fmt.Println("Command add")
		case "list":
			// TODO:
			fmt.Println("Command list")
		case "get":
			// TODO:
			fmt.Println("Command get")
		default:
			fmt.Printf("Unknown command: %v\n", command)
		}
	}

	// NOTE: Example of adding data to the database
	_, err = store.RunQuery(`
	       CREATE TABLE IF NOT EXISTS entries (
	           title TEXT NOT NULL,
	           content TEXT NOT NULL
	       );
	   `)
	if err != nil {
		fmt.Printf("Error %w\n", err)
	}
	_, err = store.RunQuery(`
		INSERT INTO entries VALUES ("my title", "my content")
	`)
	if err != nil {
		fmt.Printf("Error %w\n", err)
	}
	_, err = store.RunQuery(`
		INSERT INTO entries VALUES ("my title 2", "my content 2")
	`)
	if err != nil {
		fmt.Printf("Error %w\n", err)
	}
	_, err = store.RunQuery(`
		INSERT INTO entries VALUES ("my title 3", "my content 3")
	`)
	if err != nil {
		fmt.Printf("Error %w\n", err)
	}
	result, err := store.RunQuery(`
		SELECT * FROM entries
	`)
	if err != nil {
		fmt.Printf("Error %w\n", err)
	}

	fmt.Printf("Result:\n")
	for k, v := range result {
		fmt.Printf("\t%s : %s\n", k, v)
    }

	resultMapList, err := store.QueryToMap(`
		SELECT * FROM entries
	`)
	if err != nil {
		fmt.Printf("Error %w\n", err)
	}

	fmt.Printf("Result map:\n")
	for key := range(resultMapList[0]) {
		fmt.Printf("\t%s:", key)
		for _, row := range resultMapList{
			fmt.Printf("\t%s", row[key])
		}
		fmt.Println("")
	}
}
