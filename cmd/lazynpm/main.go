package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Elcris1/lazynpm/internal/app"
	"github.com/Elcris1/lazynpm/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	// ANSI codes to enter/exit the alternate screen buffer
	enterAltScreen = "\x1b[?1049h"
	exitAltScreen  = "\x1b[?1049l"
)

type packageView int

const (
	PackageSelectionView packageView = iota
	InfoView
	ScriptsView
	DependenciesView
	DevDependenciesView

	numberOfViews
)

type model struct {
	choices []app.Package
	cursor  int
	//struct{} => 0 bytes in memory, bool 1 byte
	selected       map[int]struct{}
	selectedChoice int
	selectedView   packageView
}

func initialModel(packages []app.Package) model {
	return model{
		choices:        packages,
		selected:       make(map[int]struct{}),
		selectedChoice: 0,
		selectedView:   PackageSelectionView,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		stringKey := msg.String()
		m.selectViewBindings(stringKey)
		m.cursorBindings(&stringKey, &m.cursor, m.calculateAmountOptions())
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) getSelectedChoice() app.Package {
	return m.choices[m.selectedChoice]
}

func (m model) getEntrFromSelectedView() app.PacakgeEntry {
	switch m.selectedView {
	case InfoView:
		return m.getSelectedChoice().Info
	case ScriptsView:
		return m.getSelectedChoice().Scripts
	case DependenciesView:
		return m.getSelectedChoice().Dependecies
	case DevDependenciesView:
		return m.getSelectedChoice().DevDependencies
	}
	return app.PacakgeEntry{}
}

func (m *model) selectViewBindings(key string) {

	switch key {
	case "1":
		m.selectedView = PackageSelectionView
		m.cursor = m.selectedChoice
	case "2":
		m.changeView(m.selectedView, m.cursor)
		m.selectedView = InfoView
	case "3":
		m.changeView(m.selectedView, m.cursor)
		m.selectedView = ScriptsView
	case "4":
		m.changeView(m.selectedView, m.cursor)
		m.selectedView = DependenciesView
	case "5":
		m.changeView(m.selectedView, m.cursor)
		m.selectedView = DevDependenciesView
	}
}

func (m *model) changeView(originalView packageView, cursor int) {
	if originalView == PackageSelectionView {
		m.selectedChoice = cursor
	}
	m.cursor = 0
}

func (m model) cursorBindings(key *string, cursor *int, cap int) {
	switch *key {
	case "up", "k":
		if *cursor > 0 {
			*cursor--
		} else if *cursor-1 < 0 {
			*cursor = cap
		}
	case "down", "j":
		if *cursor < cap {
			*cursor++
		} else if *cursor+1 > cap {
			*cursor = 0
		}
	}
}

func (m model) calculateAmountOptions() int {
	switch m.selectedView {
	case PackageSelectionView:
		return len(m.choices) - 1
	case InfoView:
		return 4
	case ScriptsView:
		return len(m.getSelectedChoice().Scripts.Keys) - 1
	case DependenciesView:
		return len(m.getSelectedChoice().Dependecies.Keys) - 1
	case DevDependenciesView:
		return len(m.getSelectedChoice().DevDependencies.Keys) - 1
	}
	return 0
}

func (m model) View() string {
	leftBox := ""
	var leftBoxWidth int = ui.GetLeftColumnWidth()
	var rightBoxWidth int = ui.GetRightColumnWidth()

	for i := range numberOfViews {
		view := packageView(i)
		leftBox += m.renderView(view, &leftBoxWidth)
	}

	rightBox := m.renderValueView(&rightBoxWidth)

	s := ui.TwoColumn(&leftBox, &rightBox) + "\nPress q to quit.\n"
	return s
}

func (m model) renderView(view packageView, width *int) string {
	switch view {
	case PackageSelectionView:
		return m.renderPackageSelectionView(width) + "\n"
	case InfoView:
		return m.renderInfoView(width) + "\n"
	case ScriptsView:
		return m.renderScriptsView(width) + "\n"
	case DependenciesView:
		return m.renderDependenciesView(width) + "\n"
	case DevDependenciesView:
		return m.renderDevDependenciesView(width) + "\n"
	}
	return ""
}

func (m model) renderPackageSelectionView(width *int) string {
	s := ""
	var isSelected bool = m.selectedView == PackageSelectionView
	if isSelected {
		for i, val := range m.choices {

			s += ui.RenderText(val.Info.Data[string(app.FileLocation)], m.cursor == i)
		}
	} else {

		s += ui.RenderText(m.getSelectedChoice().Info.Data[string(app.FileLocation)], false)
	}
	return ui.RenderBox("[1] Packages in repo", s, *width, isSelected)
}

func (m model) renderInfoView(width *int) string {

	var isSelected bool = m.selectedView == InfoView

	s := fmt.Sprintf("%s\n", m.getSelectedChoice().Info.Data[string(app.Name)])
	if isSelected {
		s += fmt.Sprintf("%s\n", m.getSelectedChoice().Info.Data[string(app.FileLocation)])
		s += fmt.Sprintf("%s\n", m.getSelectedChoice().Info.Data[string(app.Version)])
		s += fmt.Sprintf("%s\n", m.getSelectedChoice().Info.Data[string(app.PackageManager)])
	}

	return ui.RenderBox("[2] Package Information", s, *width, isSelected)
}

func (m model) getTextForKeyView(view packageView) string {
	var keys []string
	switch view {
	case ScriptsView:
		keys = m.getSelectedChoice().Scripts.Keys
	case DependenciesView:
		keys = m.getSelectedChoice().Dependecies.Keys
	case DevDependenciesView:
		keys = m.getSelectedChoice().DevDependencies.Keys
	}

	if len(keys) > 0 {
		return keys[0]
	}
	return "No available content"
}

func (m model) getTextForValueView() string {
	entry := m.getEntrFromSelectedView()
	if len(entry.Keys) > 0 {
		return entry.Data[entry.Keys[m.cursor]]
	}
	return "No available information"
}

func (m model) renderScriptsView(width *int) string {
	s := ""
	var isSelected bool = m.selectedView == ScriptsView
	if isSelected {
		for i, key := range m.getSelectedChoice().Scripts.Keys {
			s += ui.RenderText(key, m.cursor == i)
		}
	} else {
		s += ui.RenderText(m.getTextForKeyView(ScriptsView), false)
	}
	return ui.RenderBox("[3] Scripts", s, *width, isSelected)
}

func (m model) renderDependenciesView(width *int) string {
	s := ""
	var isSelected bool = m.selectedView == DependenciesView
	if isSelected {
		for i, key := range m.getSelectedChoice().Dependecies.Keys {
			s += ui.RenderText(key, m.cursor == i)
		}
	} else {
		s += ui.RenderText(m.getTextForKeyView(DependenciesView), false)
	}
	return ui.RenderBox("[4] Dependencies", s, *width, isSelected)
}

func (m model) renderDevDependenciesView(width *int) string {
	s := ""
	var isSelected bool = m.selectedView == DevDependenciesView
	if isSelected {
		for i, key := range m.getSelectedChoice().DevDependencies.Keys {
			s += ui.RenderText(key, m.cursor == i)
		}
	} else {
		s += ui.RenderText(m.getTextForKeyView(DevDependenciesView), false)
	}
	return ui.RenderBox("[5] DevDependencies", s, *width, isSelected)
}

func (m model) renderValueView(width *int) string {
	s := ""
	switch m.selectedView {
	case PackageSelectionView, InfoView:
		s += ui.RenderText("Nothing to show", false)
	default:
		s += ui.RenderText(m.getTextForValueView(), false)
	}

	return ui.RenderBox("[ENTER] Additional information", s, *width, false)
}

func findFiles(root, pattern string) ([]string, error) {
	//make is not used to maker reolcation easier because its not likely to have multiple package json in directory
	var matches []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		matched, err := filepath.Match(pattern, d.Name())
		if err != nil {
			return err
		}

		if matched {
			matches = append(matches, path)
		}
		return nil
	})
	return matches, err
}

func main() {

	files, error := findFiles(".", "package.json")
	if error != nil {
		fmt.Println("There was an error finding package files")
		return
	}

	var packages []app.Package
	for _, val := range files {
		var pkg app.Package
		err := app.LoadPackage(val, &pkg)
		packages = append(packages, pkg)
		if err != nil {
			fmt.Println("Parsing file failed")
			return
		}
	}

	fmt.Print(enterAltScreen)
	defer fmt.Print(exitAltScreen)

	p := tea.NewProgram(initialModel(packages))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
