package cliapp

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) handleScreenInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.currScreenFocus {
	case MainMenuScreen:
		return m.handleMainMenuScreen(msg)
	case ViewTablesScreen:
		return m.handleViewTablesMenuScreen(msg)
	case ManageTableScreen:
		return m.handleManageTableMenuScreen(msg)
	case AddEntryScreen:
		return m.handleAddEntryScreen(msg)
	case SelectQueryScreen:
		// TODO:
	case TableScreen:
		return m.handleTableScreen(msg)
	}

	return m, cmd
}

func (m Model) handleMainMenuScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, selectKey):
		item, ok := m.mainMenuList.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}
		switch item.title {
		case "Manage tables":
			return m.navigateTo(m.currScreenLeft, ViewTablesScreen), nil

		case "Run queries":
			var cmd tea.Cmd
			return m, cmd
			// TODO: Implement this
		}
	}

	var cmd tea.Cmd
	m.mainMenuList, cmd = m.mainMenuList.Update(msg)
	return m, cmd
}

func (m Model) handleViewTablesMenuScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, selectKey):
		item, ok := m.viewTablesMenuList.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}

		m.currTableName = item.title
		m.currTableDesc = item.desc

		return m.navigateTo(m.currScreenLeft, ManageTableScreen), nil
	}

	var cmd tea.Cmd
	m.viewTablesMenuList, cmd = m.viewTablesMenuList.Update(msg)
	return m, cmd
}

func (m Model) handleManageTableMenuScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editEntryMode {
		return m.navigateTo(None, TableScreen), nil
	}

	switch {
	case key.Matches(msg, selectKey):
		item, ok := m.manageTableMenuList.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}

		switch item.title {
		case "View table":
			return m.navigateTo(m.currScreenLeft, TableScreen), nil

		case "Add entry":
			return m.navigateTo(m.currScreenLeft, AddEntryScreen), nil

		case "Edit entry":
			if len(m.table.Rows()) > 0 {
				m.editEntryMode = true
				return m.navigateTo(m.currScreenLeft, TableScreen), nil
			} else {
				m.status.Message = fmt.Sprintf("No entries in the table %s", m.currTableName)
			}
		}
	}

	var cmd tea.Cmd
	m.manageTableMenuList, cmd = m.manageTableMenuList.Update(msg)
	return m, cmd
}

func (m Model) handleTableScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.showConfirmModal {
		switch {
		case key.Matches(msg, nextEntryKey), key.Matches(msg, leftKey), key.Matches(msg, rightKey):
			m.confirmModalSelected = !m.confirmModalSelected
			return m, nil
		case key.Matches(msg, selectKey):
			m.showConfirmModal = false
			if m.confirmModalSelected {
				if err := m.confirmModalClosure(); err != nil {
					m.status.State = StatusBarStateErr
					m.status.Message = fmt.Sprintf("Error: %s", err)
				}
			}
			return m.navigateTo(None, TableScreen), nil
		}
	}

	switch {
	case key.Matches(msg, selectKey):
		// Draw a modal with the data in the current entry
		m.showTableModal = !m.showTableModal
	}

	if m.editEntryMode {
		switch {
		case key.Matches(msg, deleteEntryKey):
			m.showConfirmModal = true
			m.confirmModalSelected = false
			m.confirmModalText = "Are you sure you want to remove the selected entry?"
			m.confirmModalClosure = func() error {
				return m.backend.DeleteEntryFromTable(m.currTableName, m.table.SelectedRow())
			}

		case key.Matches(msg, editEntryKey):
			return m.navigateTo(None, AddEntryScreen), nil
		case key.Matches(msg, undoKey):
			if err := m.backend.UndoLatest(); err != nil {
				m.status.State = StatusBarStateErr
				m.status.Message = fmt.Sprintf("Error undoing operation: %s", err)
			}
			return m.navigateTo(None, TableScreen), nil
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	m.tableModalIndex = m.table.Cursor()
	return m, cmd
}

func (m Model) handleAddEntryScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.showConfirmModal {
		switch {
		case key.Matches(msg, nextEntryKey), key.Matches(msg, leftKey), key.Matches(msg, rightKey):
			m.confirmModalSelected = !m.confirmModalSelected
			return m, nil
		case key.Matches(msg, selectKey):
			m.showConfirmModal = false
			if m.confirmModalSelected {
				if err := m.confirmModalClosure(); err != nil {
					m.status.State = StatusBarStateErr
					m.status.Message = fmt.Sprintf("Error: %s", err)
				}
				m = m.navigateTo(None, ManageTableScreen).(Model)
				return m.navigateTo(None, TableScreen), nil
			}
			return m, nil
		}
	}
	switch {
	case key.Matches(msg, sendFormKey):

		if m.editEntryMode {
			m.showConfirmModal = true
			m.confirmModalSelected = false
			m.confirmModalText = "Are you sure you want to update the selected entry?"
			m.confirmModalClosure = func() error {
				return m.addOrEditEntry()
			}
			return m, nil
		}

		err := m.addOrEditEntry()
		if err != nil {
			m.status.State = StatusBarStateErr
			m.status.Message = fmt.Sprintf("Error: %s", err)
		}
		return m.NavigateBack(), nil
	}

	var cmd tea.Cmd
	m.addEntryForm, cmd = m.addEntryForm.Update(msg)
	return m, cmd
}

func (m *Model) addOrEditEntry() error {
	keys, values, err := m.addEntryForm.GetValues()
	if err != nil {
		return err
	}

	if m.editEntryMode {
		err = m.backend.UpdateRowInTable(m.currTableName, m.currRowFilterCol, m.currRowFilterVal, keys,
			m.table.SelectedRow(), values)
	} else {
		err = m.backend.AddEntryToTable(m.currTableName, keys, values)
	}
	if err != nil {
		return err
	}

	err = m.setupTable(false, false)
	if err != nil {
		m.status.State = StatusBarStateErr
		m.status.Message = fmt.Sprintf("Error updating table %s: %s", m.currTableName, err)
		return err
	}

	m.status.State = StatusBarStateOk
	if m.editEntryMode {
		m.status.Message = fmt.Sprintf("Updated entry in table %s", m.currTableName)
		m.editEntryMode = false
	} else {
		m.status.Message = fmt.Sprintf("Added entry to table %s", m.currTableName)
	}

	return nil
}

func (m Model) NavigateBack() tea.Model {
	if len(m.screenHistory) > 0 {
		screen := m.screenHistory[len(m.screenHistory)-1]
		m.screenHistory = m.screenHistory[:len(m.screenHistory)-1]
		return m.navigateTo(None, screen)
	}

	return m
}

func (m Model) navigateTo(fromScreen, toScreen Screen) tea.Model {
	if fromScreen != None {
		m.screenHistory = append(m.screenHistory, fromScreen)
	}

	m.currScreenFocus = toScreen

	switch toScreen {
	case MainMenuScreen:
		m.currScreenLeft = toScreen
	case ViewTablesScreen:
		m.currScreenLeft = toScreen
		m.viewTablesMenuList.SetItems(m.getTablesAsMenuItems())
	case ManageTableScreen:
		m.currScreenLeft = toScreen

		err := m.setupTable(fromScreen == ViewTablesScreen, fromScreen == ViewTablesScreen)
		if err != nil {
			m.status.State = StatusBarStateErr
			m.status.Message = fmt.Sprintf("Error setting up table %s: %s", m.currTableName, err)
			return m
		}

		// Set the style of the table to plain
		s := table.DefaultStyles()
		s.Selected = lipgloss.NewStyle()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		m.table.SetStyles(s)
	case TableScreen:
		err := m.setupTable(fromScreen == ViewTablesScreen, fromScreen == ViewTablesScreen)
		if err != nil {
			m.status.State = StatusBarStateErr
			m.status.Message = fmt.Sprintf("Error setting up table %s: %s", m.currTableName, err)
			return m
		}

		// Set the style of the selected row when navigating it
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = tableSelectedItemHighlight
		m.table.SetStyles(s)

	case AddEntryScreen:
		m.currScreenLeft = toScreen

		err := m.setupEntryForm()
		if err != nil {
			m.status.State = StatusBarStateErr
			m.status.Message = fmt.Sprintf("Error setting form for table %s: %s", m.currTableName, err)
			return m
		}
	}

	return m
}

func (m Model) getTablesAsMenuItems() []list.Item {
	tableSpecs := m.backend.GetTableSpecs()

	menuItems := make([]list.Item, len(tableSpecs))
	for i, spec := range tableSpecs {
		menuItems[i] = menuItem{
			title: spec.Table.Name,
			desc:  spec.Table.Description,
		}
	}
	return menuItems
}

func (m *Model) setupTable(showStatusOk bool, newTable bool) error {
	m.manageTableMenuList.Title = fmt.Sprintf("Manage table - %s", m.currTableName)
	m.showTable = true

	// Get the specs for columns in the table
	columnSpecs, err := m.backend.GetColumnSpecsInTable(m.currTableName)
	if err != nil {
		return fmt.Errorf("Error reading columns: %w", err)
	}

	// Load the data from the table
	result, err := m.backend.GetEntriesInTable(m.currTableName)
	if err != nil {
		return fmt.Errorf("Error reading table entries: %w", err)
	}
	if showStatusOk {
		m.status.State = StatusBarStateOk
		m.status.Message = fmt.Sprintf("Read table %v", m.currTableName)
	}

	tableColumns := make([]table.Column, len(columnSpecs))
	for i, colSpec := range columnSpecs {
		tableColumns[i] = table.Column{
			Title: colSpec.Name,
			Width: len(colSpec.Name) + 3,
		}
		if colSpec.WidthHint != 0 {
			tableColumns[i].Width = max(tableColumns[i].Width, colSpec.WidthHint)
		}
	}

	tableRows := make([]table.Row, len(result))

	if len(result) > 0 {
		for i, row := range result {
			tableRows[i] = row
		}
	}

	if newTable {
		m.table = table.New(
			table.WithFocused(true),
			table.WithHeight(m.height-12),
			table.WithColumns(tableColumns),
			table.WithRows(tableRows),
		)
	} else {
		m.table.SetColumns(tableColumns)
		m.table.SetRows(tableRows)
	}

	m.currTableNrEntries = len(tableRows)

	return nil
}

func (m *Model) setupEntryForm() error {
	columnSpecs, err := m.backend.GetColumnSpecsInTable(m.currTableName)
	if err != nil {
		return fmt.Errorf("Error getting the column specs for table %s: %w", m.currTableName, err)
	}

	m.addEntryForm.FocusedField = 0
	m.addEntryForm.Title = fmt.Sprintf("Add entry to the table %s", m.currTableName)

	// Populate the elements in the form
	m.addEntryForm.Fields = []FormField{}

	for _, colSpec := range columnSpecs {
		if colSpec.AutoIncrement {
			continue
		}

		var field FormField
		field.Label = colSpec.Name
		field.DataType = colSpec.Type
		field.DefaultValue = ""
		field.Required = (!colSpec.Nullable)

		field.Input = textinput.New()

		m.addEntryForm.Fields = append(m.addEntryForm.Fields, field)
	}

	if m.editEntryMode {
		// Populate the form with the current values
		fieldIdx := 0
		for i, value := range m.table.SelectedRow() {
			if columnSpecs[i].PrimaryKey {
				m.currRowFilterCol = columnSpecs[i].Name
				m.currRowFilterVal = value
			}
			if value != "NULL" {
				m.addEntryForm.Fields[fieldIdx].Input.SetValue(value)
			}
			fieldIdx++
		}
	}

	m.addEntryForm.Fields[0].Input.Focus()

	return nil
}
