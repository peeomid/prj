package cmd

import (
	"fmt"
	"strings"

	"github.com/Osimify/prj/internal/display"
	"github.com/Osimify/prj/internal/store"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show full details for a single project",
	Long: `Display everything known about a project: description, tech stack,
git history, recent commits, contributors, deployment config,
reference files, TODO counts, and more.

Supports partial name matching — you don't need the exact name.

Examples:
  prj info myapp           Exact match on project name "myapp"
  prj info api             Partial match — finds "my-api-server" etc.
  prj info openclaw        Full detail view for openclaw`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.ToLower(args[0])

		projects, err := store.Load()
		if err != nil {
			return fmt.Errorf("load projects: %w", err)
		}

		for _, p := range projects {
			if strings.ToLower(p.Name) == name {
				display.PrintDetail(p)
				return nil
			}
		}

		// Partial match
		for _, p := range projects {
			if strings.Contains(strings.ToLower(p.Name), name) {
				display.PrintDetail(p)
				return nil
			}
		}

		return fmt.Errorf("project not found: %s", args[0])
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
