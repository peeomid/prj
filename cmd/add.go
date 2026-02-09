package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peeomid/prj/internal/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <folder>",
	Short: "Add a parent folder to the scan list",
	Long: `Add a parent folder that contains git repositories to the scan list.
When you run "prj scan", every immediate child directory with a .git
folder inside this parent will be detected as a project.

Supports ~ expansion and relative paths.

Examples:
  prj add ~/Development         Add your main dev folder
  prj add ~/Projects            Add another folder
  prj add .                     Add the current directory`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		folder := expandPath(args[0])

		info, err := os.Stat(folder)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("folder does not exist: %s", folder)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		if !cfg.AddFolder(folder) {
			fmt.Printf("Already tracked: %s\n", folder)
			return nil
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}

		fmt.Printf("Added: %s\n", folder)
		return nil
	},
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~/") || p == "~" {
		home, _ := os.UserHomeDir()
		p = filepath.Join(home, p[1:])
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

func init() {
	rootCmd.AddCommand(addCmd)
}
