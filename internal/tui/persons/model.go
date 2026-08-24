// Package persons provides the interactive bubbletea persons browser.
package persons

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/mylinden-tech/linden-cli/internal/client"
)

// LoadedMsg carries persons loaded from the API.
type LoadedMsg struct {
	Persons []client.Person
	Err     error
}

// Model is the bubbletea model for the persons browser.
// Implements tea.Model (charm.land/bubbletea/v2).
type Model struct {
	persons  []client.Person
	filtered []client.Person
	cursor   int
	filter   string
	loading  bool
	err      error
	selected *client.Person
	detail   bool
	quitting bool
}

// New creates a new Model in loading state.
func New() Model {
	return Model{loading: true}
}

// Init satisfies tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update satisfies tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoadedMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.persons = msg.Persons
		m.filtered = msg.Persons

	case tea.KeyPressMsg:
		if m.quitting {
			return m, tea.Quit
		}
		if m.loading || m.err != nil {
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}
		if m.detail {
			return m.updateDetail(msg)
		}
		return m.updateList(msg)
	}

	return m, nil
}

func (m Model) updateList(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.quitting = true
		return m, tea.Quit

	case "enter":
		if len(m.filtered) > 0 {
			sel := m.filtered[m.cursor]
			m.selected = &sel
			m.detail = true
		}

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}

	case "backspace":
		if len(m.filter) > 0 {
			m.filter = m.filter[:len(m.filter)-1]
			m.applyFilter()
			if m.cursor >= len(m.filtered) {
				m.cursor = maxInt(0, len(m.filtered)-1)
			}
		}

	default:
		if len(msg.String()) == 1 {
			m.filter += msg.String()
			m.applyFilter()
			m.cursor = 0
		}
	}

	return m, nil
}

func (m Model) updateDetail(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.detail = false
		m.selected = nil
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "enter":
		// Keep selected and quit — caller will print JSON
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

// View satisfies tea.Model.
func (m Model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	if m.loading {
		v := tea.NewView("\n  Loading persons…")
		v.AltScreen = true
		return v
	}
	if m.err != nil {
		v := tea.NewView(fmt.Sprintf("\n  Error: %s\n\n  Press q to quit.", m.err))
		v.AltScreen = true
		return v
	}
	if m.detail && m.selected != nil {
		v := tea.NewView(m.renderDetail())
		v.AltScreen = true
		return v
	}
	v := tea.NewView(m.renderList())
	v.AltScreen = true
	return v
}

// Selected returns the person chosen by the user (nil if none / quit).
func (m Model) Selected() *client.Person {
	return m.selected
}

func (m Model) renderList() string {
	var sb strings.Builder

	title := lipgloss.NewStyle().Bold(true).Render("  Persons")
	sb.WriteString(title + "\n")

	if m.filter != "" {
		filterLine := lipgloss.NewStyle().Faint(true).Render(
			fmt.Sprintf("  Filter: %s", m.filter))
		sb.WriteString(filterLine + "\n")
	}
	sb.WriteString("\n")

	if len(m.filtered) == 0 {
		sb.WriteString("  No persons found.\n")
	}

	for i, p := range m.filtered {
		name := p.FirstName + " " + p.LastName
		rel := ""
		if p.RelationshipType != nil {
			rel = *p.RelationshipType
		}
		email := ""
		if p.Email != nil {
			email = *p.Email
		}

		row := fmt.Sprintf("  %-28s %-20s %s", name, rel, email)
		if i == m.cursor {
			row = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Render(row)
		}
		sb.WriteString(row + "\n")
	}

	sb.WriteString("\n")
	hint := lipgloss.NewStyle().Faint(true).Render(
		"  ↑/↓ navigate · enter view · type to filter · q quit")
	sb.WriteString(hint + "\n")
	return sb.String()
}

func (m Model) renderDetail() string {
	p := m.selected
	var sb strings.Builder

	title := lipgloss.NewStyle().Bold(true).Render(
		fmt.Sprintf("  %s %s", p.FirstName, p.LastName))
	sb.WriteString(title + "\n\n")

	field := func(label, value string) {
		if value == "" {
			return
		}
		l := lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("  %-20s", label+":"))
		sb.WriteString(l + " " + value + "\n")
	}

	field("ID", p.ID)
	if p.RelationshipType != nil {
		field("Relationship", *p.RelationshipType)
	}
	if p.Email != nil {
		field("Email", *p.Email)
	}
	if p.Phone != nil {
		field("Phone", *p.Phone)
	}
	if p.Birthday != nil {
		field("Birthday", *p.Birthday)
	}
	if p.Notes != nil && *p.Notes != "" {
		field("Notes", *p.Notes)
	}
	if addr := buildAddress(p); addr != "" {
		field("Address", addr)
	}

	sb.WriteString("\n")
	hint := lipgloss.NewStyle().Faint(true).Render(
		"  enter select · esc back · q quit")
	sb.WriteString(hint + "\n")
	return sb.String()
}

func (m *Model) applyFilter() {
	if m.filter == "" {
		m.filtered = m.persons
		return
	}
	lower := strings.ToLower(m.filter)
	filtered := m.filtered[:0]
	for _, p := range m.persons {
		name := strings.ToLower(p.FirstName + " " + p.LastName)
		email := ""
		if p.Email != nil {
			email = strings.ToLower(*p.Email)
		}
		if strings.Contains(name, lower) || strings.Contains(email, lower) {
			filtered = append(filtered, p)
		}
	}
	m.filtered = filtered
}

func buildAddress(p *client.Person) string {
	var parts []string
	add := func(s *string) {
		if s != nil && *s != "" {
			parts = append(parts, *s)
		}
	}
	add(p.AddressStreet)
	add(p.AddressApt)
	if p.AddressCity != nil && *p.AddressCity != "" {
		city := *p.AddressCity
		if p.AddressState != nil && *p.AddressState != "" {
			city += ", " + *p.AddressState
		}
		parts = append(parts, city)
	}
	add(p.AddressZip)
	add(p.AddressCountry)
	return strings.Join(parts, " · ")
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
