package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

//-----------------------------------------core---------------------------------------------------------//

// NewRootCmd tui main entry point
func NewRootCmd(db *Database) *cobra.Command {
	cmd := &cobra.Command{
		Use: "rete",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(NewFormModel(db), tea.WithAltScreen())
			_, err := p.Run()
			return err
		},
	}



	return cmd
}