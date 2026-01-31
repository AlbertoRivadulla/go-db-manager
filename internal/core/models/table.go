package core

import (
	"fmt"
	"strings"
)

type TableSpec struct {
	Table Table `yaml:"table"`
}

type Table struct {
	Name string `yaml:"name"`
	Schema string `yaml:"schema"`
	PrimaryKey string `yaml:"primary_key"`
	ForeignKeys []ForeignKey `yaml:"foreign_keys"`
	Columns []Column `yaml:"columns"`
}

type ForeignKey struct {
	Column string `yaml:"column"`
	References Reference `yaml:"references"`
	OnDelete string `yaml:"on_delete"`
	OnUpdate string `yaml:"on_update"`
}

type Reference struct {
	Table string `yaml:"table"`
	Column string `yaml:"column"`
}

type Column struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	AutoIncrement bool `yaml:"auto_increment"`
	Nullable bool `yaml:"nullable"`
	Default *string `yaml:"default"`
	Check *string `yaml:"check"`
	Comment string `yaml:"comment"`
}

func (s TableSpec) Validate() error {
	return s.Table.Validate()
}

func (t Table) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("table name is required")
	}

	columnSet := make(map[string]struct{})

	// Validate columns
	for _, c := range t.Columns {
		if _, exists := columnSet[c.Name]; exists {
			return fmt.Errorf("duplicate column with name %s", c.Name)
		}
		columnSet[c.Name] = struct{}{}

		if err := c.Validate(); err != nil {
			return fmt.Errorf("column %s: %w", c.Name, err)
		}
	}

	// Check that the primary key exists
	if _, ok := columnSet[t.PrimaryKey]; !ok {
		return fmt.Errorf("primary key column %s does not exist", t.PrimaryKey)
	}

	return nil
}

func (c Column) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}

	if c.Type == "" {
		return fmt.Errorf("type is required")
	}

	if c.AutoIncrement && !strings.Contains(strings.ToLower(c.Type), "int") {
		return fmt.Errorf("auto_increment requires type int")
	}

	return nil
}
