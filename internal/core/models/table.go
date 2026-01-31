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

func (t Table) QueryCreate() (string, error) {
	var defs []string

	for _, col := range t.Columns {
		defs = append(defs, col.QueryDefinition())
	}

	if t.PrimaryKey != "" {
		defs = append(
			defs,
			fmt.Sprintf(
				"PRIMARY KEY (%s)",
				t.PrimaryKey,
			),
		)
	}

	for _, fk := range t.ForeignKeys {
		defs = append(defs, fk.QueryDefinition())
	}

	stmt := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n  %s\n);",
		t.Name,
		strings.Join(defs, ",\n  "),
	)

	return stmt, nil
}

func (c Column) QueryDefinition() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("%s %s", c.Name, c.Type))

	if !c.Nullable {
		parts = append(parts, "NOT NULL")
	}

	if c.Default != nil {
		parts = append(parts, fmt.Sprintf("DEFAULT '%s'", *c.Default))
	}

	if c.Check != nil {
		parts = append(parts, fmt.Sprintf("CHECK (%s)", *c.Check))
	}

	return strings.Join(parts, " ")
}

func (fk ForeignKey) QueryDefinition() string {
	stmt := fmt.Sprintf(
		"FOREIGN KEY (%s) REFERENCES %s(%s)",
		fk.Column,
		fk.References.Table,
		fk.References.Column,
	)

	if fk.OnDelete != "" {
		stmt += " ON DELETE " + strings.ToUpper(fk.OnDelete)
	}

	if fk.OnUpdate != "" {
		stmt += " ON UPDATE " + strings.ToUpper(fk.OnUpdate)
	}

	return stmt
}
