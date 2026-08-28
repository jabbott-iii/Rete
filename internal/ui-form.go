/*
Copyright 2026 Joseph Anthony Abbott III

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

//------------------------------------------TUI dashboard----------------------------------------------//

// DashboardModel is the BubbleTea model for the interactive TUI dashboard.
type DashboardModel struct {
	db         *Database
	categories []*FeatureCategory
	recentJobs []*ScanJob
	cursor     int
	err        error
	loading    bool
	width      int
	height     int
}

// dashboardLoadedMsg is sent when the initial data fetch completes.
type dashboardLoadedMsg struct {
	categories []*FeatureCategory
	recentJobs []*ScanJob
}

// NewDashboardModel creates the TUI dashboard.
func NewDashboardModel(db *Database) *DashboardModel {
	return &DashboardModel{db: db, loading: true}
}

// NewFormModel is an alias kept for compatibility with main.go.
func NewFormModel(db *Database) *DashboardModel {
	return NewDashboardModel(db)
}

func (m *DashboardModel) Init() tea.Cmd {
	return func() tea.Msg {
		cats, err := m.db.ListCategories()
		if err != nil {
			return errMsg{err}
		}
		jobs, err := m.db.ListJobs(10)
		if err != nil {
			return errMsg{err}
		}
		return dashboardLoadedMsg{categories: cats, recentJobs: jobs}
	}
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case dashboardLoadedMsg:
		m.categories = msg.categories
		m.recentJobs = msg.recentJobs
		m.loading = false

	case errMsg:
		m.err = msg.error
		m.loading = false

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			m.cursor++
		}
	}
	return m, nil
}

func (m *DashboardModel) View() string {
	var sb strings.Builder

	sb.WriteString(styleBanner.Render("  rete") + "  network diagnostic & penetration testing\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n\n")

	if m.loading {
		sb.WriteString("  Loading…\n")
		return sb.String()
	}

	if m.err != nil {
		sb.WriteString(fmt.Sprintf("  Error: %v\n", m.err))
		return sb.String()
	}

	// Feature catalog
	if len(m.categories) == 0 {
		sb.WriteString("  No features loaded.\n\n")
	} else {
		sb.WriteString("  Features\n")
		sb.WriteString("  " + strings.Repeat("─", 40) + "\n")
		for _, cat := range m.categories {
			sb.WriteString("  " + styleCategory.Render(cat.Name) + "\n")
			for _, f := range cat.Features {
				line := fmt.Sprintf("    %-34s  rete %s", f.Name, f.CobraCommand)
				sb.WriteString(styleFeature.Render(line) + "\n")
			}
		}
		sb.WriteString("\n")
	}

	// Recent jobs
	sb.WriteString("  Recent Jobs\n")
	sb.WriteString("  " + strings.Repeat("─", 40) + "\n")
	if len(m.recentJobs) == 0 {
		sb.WriteString("  No jobs yet.\n")
	} else {
		for _, j := range m.recentJobs {
			sb.WriteString(fmt.Sprintf("  #%-4d  %-12s  %s\n", j.ID, j.Status, j.Command))
		}
	}

	sb.WriteString("\n")
	sb.WriteString(styleHelp.Render("  q quit  ↑/↓ scroll") + "\n")

	return sb.String()
}

type errMsg struct{ error }