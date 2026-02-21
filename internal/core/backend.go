package core

import (
	"fmt"
	"os"
	"path/filepath"

	"dbmanager/internal/core/spec-models"

	"gopkg.in/yaml.v3"
)

type RowOperation int

const (
	UpdateRow RowOperation = iota
	DeleteRow
)

// Used for the history of updates
type RowUpdateOperation struct {
	Table        string
	Columns      []string
	Values       []string
	RowFilterCol string
	RowFilterVal string
	Operation    RowOperation
}

type Backend struct {
	Store         SQLiteStore
	tableSpecs    []core.TableSpec
	mapTableSpecs map[string]*core.TableSpec
	// TODO: queries
	// TODO: rules

	UndoHistory []RowUpdateOperation
}

func NewBackend(dbPath *string, specsDir *string) (*Backend, error) {

	// TODO: Sanitize the database and specs paths

	store, err := NewSQLiteStore(*dbPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening the SQLite store: %w", err)
	}

	backend := Backend{
		Store: *store,
	}

	err = backend.parseSpecs(specsDir)
	if err != nil {
		return nil, err
	}

	// Create the different tables
	for _, tableSpec := range backend.tableSpecs {
		err := backend.CreateOrModifyTable(tableSpec)
		if err != nil {
			return nil, err
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

func (backend *Backend) GetColumnSpecsInTable(table string) ([]core.Column, error) {
	tableSpec, ok := backend.mapTableSpecs[table]
	if !ok {
		return nil, fmt.Errorf("table %s not found", table)
	}

	return tableSpec.Table.Columns, nil
}

func (backend Backend) CreateOrModifyTable(tableSpec core.TableSpec) error {
	queryTableInfo, err := tableSpec.Table.QueryTableInfo()
	if err != nil {
		return fmt.Errorf("create query to get table info: %w", err)
	}

	tableInfo, err := backend.Store.QueryTo2DimArray(queryTableInfo)
	if err != nil {
		return fmt.Errorf("get table info %s: %w", tableSpec.Table.Name, err)
	}

	if len(tableInfo) == 0 {
		queryCreate, err := tableSpec.Table.QueryCreate()
		if err != nil {
			return fmt.Errorf("get query to create table %s: %w", tableSpec.Table.Name, err)
		}

		_, err = backend.Store.RunQueryNoReturn(queryCreate)
		if err != nil {
			return fmt.Errorf("create table %s: %w", tableSpec.Table.Name, err)
		}
	} else {
		// Find the columns not existing in the table
		var newColumns []core.Column

		for _, colInSpec := range tableSpec.Table.Columns {
			exists := false
			for _, colInInfo := range tableInfo {
				if colInSpec.Name == colInInfo[1] {
					exists = true
					break
				}
			}
			if !exists {
				newColumns = append(newColumns, colInSpec)
			}
		}

		for _, newColSpec := range newColumns {
			queryAddNewCol, err := tableSpec.Table.QueryAddNewColumn(newColSpec)
			if err != nil {
				return fmt.Errorf("get query to add the new column %s: %w", newColSpec.Name, err)
			}

			_, err = backend.Store.RunQueryNoReturn(queryAddNewCol)
			if err != nil {
				return fmt.Errorf("add new column %s: %w", newColSpec.Name, err)
			}
		}
	}

	return nil
}

func (backend *Backend) GetEntriesInTable(table string) ([][]string, error) {
	query := backend.mapTableSpecs[table].Table.QueryGetAllEntries()
	result, err := backend.Store.QueryTo2DimArray(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (backend *Backend) AddEntryToTable(table string, columns []string, values []string) error {
	query := backend.mapTableSpecs[table].Table.QueryAddEntry(columns, values)
	_, err := backend.Store.RunQueryNoReturn(query)
	return err
}

func (backend *Backend) UpdateRowInTable(table string, filterCol string, filterVal string, columns []string,
	oldValues []string, newValues []string) error {
	thisOperation := RowUpdateOperation{
		Table:        table,
		Columns:      columns,
		Values:       backend.mapTableSpecs[table].Table.ValuesWithoutPrimaryKey(oldValues, true),
		RowFilterCol: filterCol,
		RowFilterVal: filterVal,
		Operation:    UpdateRow,
	}
	backend.UndoHistory = append(backend.UndoHistory, thisOperation)

	query := backend.mapTableSpecs[table].Table.QueryUpdateRow(filterCol, filterVal, columns, newValues)
	_, err := backend.Store.RunQueryNoReturn(query)
	return err
}

func (backend *Backend) DeleteEntryFromTable(table string, row []string) error {
	thisOperation := RowUpdateOperation{
		Table:     table,
		Values:    backend.mapTableSpecs[table].Table.ValuesWithoutPrimaryKey(row, false),
		Operation: DeleteRow,
	}
	backend.UndoHistory = append(backend.UndoHistory, thisOperation)

	query := backend.mapTableSpecs[table].Table.QueryDeleteEntry(row)
	_, err := backend.Store.RunQueryNoReturn(query)
	return err
}

func (backend *Backend) UndoLatest() error {
	if len(backend.UndoHistory) > 0 {
		operation := backend.UndoHistory[len(backend.UndoHistory)-1]
		backend.UndoHistory = backend.UndoHistory[:len(backend.UndoHistory)-1]

		switch operation.Operation {
		case UpdateRow:
			query := backend.mapTableSpecs[operation.Table].Table.QueryUpdateRow(operation.RowFilterCol,
				operation.RowFilterVal, operation.Columns, operation.Values)
			_, err := backend.Store.RunQueryNoReturn(query)
			return err
		case DeleteRow:
			var columns []string
			for _, col := range backend.mapTableSpecs[operation.Table].Table.Columns {
				if !col.PrimaryKey {
					columns = append(columns, col.Name)
				}
			}
			query := backend.mapTableSpecs[operation.Table].Table.QueryAddEntry(columns, operation.Values)
			_, err := backend.Store.RunQueryNoReturn(query)
			return err
		}
	}

	return nil
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
