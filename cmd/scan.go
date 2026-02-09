package cmd

import (
	"fmt"

	"github.com/peeomid/prj/internal/config"
	"github.com/peeomid/prj/internal/display"
	"github.com/peeomid/prj/internal/project"
	"github.com/peeomid/prj/internal/scanner"
	"github.com/peeomid/prj/internal/store"
	"github.com/spf13/cobra"
)

var scanDryRun bool

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan all tracked folders, extract metadata from every git repo",
	Long: `Walk through every folder in your scan list, find git repositories,
and extract rich metadata from each one:

  - Git history: last commit, recent commits, contributor list
  - Tech stack: Ruby, Node, Python, Go, Swift, Rust, and frameworks
  - Project type: rails-app, node-app, go-app, docs, script, etc.
  - Status: active (< 30d), recent (< 90d), wip, or paused
  - Deployment: Docker, Heroku, Vercel, GitHub Actions, etc.
  - Reference files: README, CLAUDE.md, .ai/, .cursor/, docs/, tasks/
  - TODO counts: open/closed items from TODO.md

Results are merged into ~/.prj/projects.json (existing projects are
updated, new ones are added).

Examples:
  prj scan               Scan and save all project data
  prj scan --dry-run     Scan but don't save (preview what would happen)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		if len(cfg.Folders) == 0 {
			return fmt.Errorf("no folders configured. Run: prj add <folder>")
		}

		var allRepos []string
		for _, folder := range cfg.Folders {
			fmt.Printf("Scanning %s ...\n", folder)
			repos, err := scanner.FindRepos(folder)
			if err != nil {
				fmt.Printf("  %s: %s\n", display.Red("error"), err)
				continue
			}
			fmt.Printf("  found %d repos\n", len(repos))
			allRepos = append(allRepos, repos...)
		}

		if len(allRepos) == 0 {
			fmt.Println("No repos found.")
			return nil
		}

		fmt.Printf("\nExtracting metadata from %d repos...\n", len(allRepos))
		var scanned []*project.Project
		for i, repoPath := range allRepos {
			fmt.Printf("  [%d/%d] %s\n", i+1, len(allRepos), repoPath)
			p := project.ExtractFromPath(repoPath)
			scanned = append(scanned, p)
		}

		if scanDryRun {
			fmt.Printf("\n%s — %d projects scanned (not saved)\n", display.Yellow("dry-run"), len(scanned))
			return nil
		}

		existing, err := store.Load()
		if err != nil {
			return fmt.Errorf("load store: %w", err)
		}

		merged := store.Merge(existing, scanned)
		if err := store.Save(merged); err != nil {
			return fmt.Errorf("save store: %w", err)
		}

		fmt.Printf("\n%s — %d projects saved to %s\n", display.Green("done"), len(merged), store.Path())
		return nil
	},
}

func init() {
	scanCmd.Flags().BoolVar(&scanDryRun, "dry-run", false, "Scan without saving")
	rootCmd.AddCommand(scanCmd)
}
