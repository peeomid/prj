package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "prj",
	Short: "Git project scanner & reporter",
	Long: `prj scans your local git repositories, extracts metadata (tech stack,
commit history, deployment config, status, etc.), and lets you search,
filter, and report on all your projects from the terminal.

Quick start:
  prj add ~/Development      Add a folder to scan
  prj scan                   Scan all folders for git repos
  prj list                   Show all projects in a table
  prj status                 Show a summary report

Data is stored in ~/.prj/ as JSON files.`,
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
