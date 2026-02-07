package core

import (
	"fmt"
	"os"
	"path/filepath"

	// _ "modernc.org/sqlite"

	"carmaintenance/internal/core/spec-models"

	"gopkg.in/yaml.v3"
)

type Backend struct {
	Store         SQLiteStore
	tableSpecs    []core.TableSpec
	mapTableSpecs map[string]*core.TableSpec
	// TODO: queries
	// TODO: rules
}

func NewBackend(dbPath *string, specsDir *string) (*Backend, error) {

	// TODO: Sanitize the database and specs paths

	store, err := NewSQLiteStore(*dbPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening the SQLite store: %w", err)
	}

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

	backend := Backend{
		Store: *store,
	}

	err = backend.parseSpecs(specsDir)
	if err != nil {
		return nil, err
	}

	// Create the different tables
	for _, tableSpec := range backend.tableSpecs {
		queryCreate, err := tableSpec.Table.QueryCreate()
		if err != nil {
			return nil, fmt.Errorf("get query to create table %s: %w", tableSpec.Table.Name, err)
		}

		_, err = backend.Store.RunQueryNoReturn(queryCreate)
		if err != nil {
			return nil, fmt.Errorf("create table %s: %w", tableSpec.Table.Name, err)
		}
	}

	// FIXME: Delete this
	for _ = range 100 {
		_, err = store.RunQuery(`
			INSERT INTO Cars (model, year, license_plate, chassis_nr, ITV_date, ITV_interval)
			VALUES ("mod", "2000", "abcd", "algo123", "1/1/2001", 12)
		`)
		if err != nil {
			fmt.Printf("Error %s\n", err)
		}
	}

	return &backend, nil
}

func (backend *Backend) GetTableSpecs() []core.TableSpec {
	return backend.tableSpecs
}

func (backend *Backend) GetColumnsInTable(table string) ([]string, error) {
	tableSpec, ok := backend.mapTableSpecs[table]
	if !ok {
		return nil, fmt.Errorf("table %s not found", table)
	}

	var columns []string
	for _, col := range tableSpec.Table.Columns {
		columns = append(columns, col.Name)
	}

	return columns, nil
}

func (backend *Backend) GetEntriesInTable(table string) ([][]string, error) {
	query := backend.mapTableSpecs[table].Table.QueryGetAllEntries()
	result, err := backend.Store.QueryTo2DimArray(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (backend *Backend) Cleanup() {
	backend.Store.db.Close()
}

func (backend *Backend) parseSpecs(specsDir *string) error {
	entries, err := os.ReadDir(*specsDir)
	if err != nil {
		return fmt.Errorf("read specs directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if entry.Name() == "tables" {
			tablesDir := filepath.Join(*specsDir, entry.Name())

			backend.tableSpecs, err = loadYAMLSpecs[core.TableSpec](tablesDir)
			if err != nil {
				return fmt.Errorf("read tables directory: %w", err)
			}

			backend.mapTableSpecs = map[string]*core.TableSpec{}
			for i, tableSpec := range backend.tableSpecs {
				backend.mapTableSpecs[tableSpec.Table.Name] = &backend.tableSpecs[i]
			}

		} else if entry.Name() == "queries" {
			// TODO: Read queries
			// queriesDir := filepath.Join(*specsDir, entry.Name())
			//
			// backend.querySpecs, err = loadYAMLSpecs[core.querySpec](queriesDir)
			// if err != nil {
			// 	return fmt.Errorf("read queries directory")
			// }

		} else if entry.Name() == "rules" {
			// TODO: Read rules
			// rulesDir := filepath.Join(*specsDir, entry.Name())
			//
			// backend.ruleSpecs, err = loadYAMLSpecs[core.ruleSpec](rulesDir)
			// if err != nil {
			// 	return fmt.Errorf("read rules directory")
			// }
		}
	}

	return nil
}

type Validatable interface {
	Validate() error
}

func loadYAMLSpecs[T any](directory string) ([]T, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var specs []T

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if ext := filepath.Ext(entry.Name()); ext != ".yaml" && ext != ".yml" {
			continue
		}

		path := filepath.Join(directory, entry.Name())
		spec, err := loadSingleYAMLSpec[T](path)
		if err != nil {
			return nil, err
		}

		specs = append(specs, *spec)
	}

	return specs, nil
}

func loadSingleYAMLSpec[T any](filePath string) (*T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filePath, err)
	}

	var spec T

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parse %s: %w", filePath, err)
	}

	if v, ok := any(&spec).(Validatable); ok {
		if err := v.Validate(); err != nil {
			return nil, fmt.Errorf("validate %s: %w", filePath, err)
		}
	}

	return &spec, nil
}
