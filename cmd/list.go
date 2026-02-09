package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/peeomid/prj/internal/display"
	"github.com/peeomid/prj/internal/project"
	"github.com/peeomid/prj/internal/store"
	"github.com/spf13/cobra"
)

var (
	listStatus string
	listType   string
	listTech   string
	listOwn    bool
	listForks  bool
	listSearch string
	listSort   string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all scanned projects in a table (with filters and sorting)",
	Long: `Display a table of all scanned projects. You can filter by status,
type, tech stack, ownership, or a search query. You can also sort
by name, last commit date, or commit count.

Statuses:
  active   — committed within the last 30 days
  wip      — active + recent commit messages contain "WIP"
  recent   — committed within the last 90 days
  paused   — no commits in 90+ days

Examples:
  prj list                          Show all projects (sorted by date)
  prj list --status active          Only active projects
  prj list --tech react             Only projects using React
  prj list --type go-app            Only Go projects
  prj list --own                    Only your own projects (not forks)
  prj list --forks                  Only forked projects
  prj list --search api             Search by name or path
  prj list --sort commits           Sort by commit count (most active first)
  prj list --sort name              Sort alphabetically
  prj list --status active --own    Combine multiple filters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := store.Load()
		if err != nil {
			return fmt.Errorf("load projects: %w", err)
		}

		// Apply filters
		filtered := filterProjects(projects)

		// Sort
		sortProjects(filtered, listSort)

		display.PrintTable(filtered)
		return nil
	},
}

func filterProjects(projects []*project.Project) []*project.Project {
	var result []*project.Project
	for _, p := range projects {
		if listStatus != "" && p.Status != listStatus {
			continue
		}
		if listType != "" && !strings.Contains(strings.ToLower(p.InferredType), strings.ToLower(listType)) {
			continue
		}
		if listTech != "" && !containsTech(p.TechStack, listTech) {
			continue
		}
		if listOwn && p.IsFork {
			continue
		}
		if listForks && !p.IsFork {
			continue
		}
		if listSearch != "" {
			q := strings.ToLower(listSearch)
			if !strings.Contains(strings.ToLower(p.Name), q) &&
				!strings.Contains(strings.ToLower(p.Path), q) {
				continue
			}
		}
		result = append(result, p)
	}
	return result
}

func containsTech(stack []string, tech string) bool {
	t := strings.ToLower(tech)
	for _, s := range stack {
		if strings.Contains(strings.ToLower(s), t) {
			return true
		}
	}
	return false
}

func sortProjects(projects []*project.Project, by string) {
	switch by {
	case "name":
		sort.Slice(projects, func(i, j int) bool {
			return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
		})
	case "commits":
		sort.Slice(projects, func(i, j int) bool {
			return projects[i].CommitCount8M > projects[j].CommitCount8M
		})
	default: // "date" or empty — sort by last commit desc
		sort.Slice(projects, func(i, j int) bool {
			return projects[i].LastCommitDate > projects[j].LastCommitDate
		})
	}
}

func init() {
	listCmd.Flags().StringVar(&listStatus, "status", "", "Filter by status (active/recent/paused/wip)")
	listCmd.Flags().StringVar(&listType, "type", "", "Filter by inferred type")
	listCmd.Flags().StringVar(&listTech, "tech", "", "Filter by tech stack")
	listCmd.Flags().BoolVar(&listOwn, "own", false, "Show only own projects (not forks)")
	listCmd.Flags().BoolVar(&listForks, "forks", false, "Show only forks")
	listCmd.Flags().StringVar(&listSearch, "search", "", "Search name/path")
	listCmd.Flags().StringVar(&listSort, "sort", "date", "Sort by: name, date, commits")
	rootCmd.AddCommand(listCmd)
}
