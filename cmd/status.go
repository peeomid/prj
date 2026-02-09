package cmd

import (
	"fmt"

	"github.com/peeomid/prj/internal/display"
	"github.com/peeomid/prj/internal/store"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Dashboard-style summary: counts by status, type, ownership + highlights",
	Long: `Show a high-level overview of all your scanned projects:

  - Total project count
  - Breakdown by status (active, wip, recent, paused)
  - Breakdown by type (rails-app, node-app, go-app, etc.)
  - Ownership split (your own repos vs forks)
  - Top 5 most recently committed projects
  - Stalled projects (no commits in 6+ months)

Run "prj scan" first to populate the data.

Examples:
  prj status                Show the full summary report`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := store.Load()
		if err != nil {
			return fmt.Errorf("load projects: %w", err)
		}

		display.PrintSummary(projects)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
