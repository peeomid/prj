package cmd

import (
	"fmt"

	"github.com/Osimify/prj/internal/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <folder>",
	Short: "Stop tracking a folder (does not delete projects data)",
	Long: `Remove a folder from the scan list. Future scans will no longer look
in this folder for repos. Already-scanned projects stay in the data
file until they are overwritten by a new scan.

Use "prj config" to see which folders are currently tracked.

Examples:
  prj remove ~/old-projects     Stop scanning this folder
  prj remove ~/Development      Remove your main dev folder`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		folder := expandPath(args[0])

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		if !cfg.RemoveFolder(folder) {
			return fmt.Errorf("folder not tracked: %s", folder)
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}

		fmt.Printf("Removed: %s\n", folder)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
