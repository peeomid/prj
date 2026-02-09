package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Osimify/prj/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration (tracked folders, settings)",
	Long: `Print the current prj configuration as JSON. This includes:

  - folders:     list of parent directories being scanned
  - cutoff_days: how many days of inactivity before a project is "paused"

Config is stored at ~/.prj/config.json.

Examples:
  prj config               Print the config file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
