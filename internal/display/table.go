package display

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/peeomid/prj/internal/project"
	"github.com/rodaine/table"
)

// PrintTable renders projects as a table.
func PrintTable(projects []*project.Project) {
	if len(projects) == 0 {
		fmt.Println("No projects found.")
		return
	}

	tbl := table.New("Name", "Type", "Status", "Tech", "Last Commit", "Commits(8m)")
	tbl.WithWriter(os.Stdout)

	for _, p := range projects {
		tech := strings.Join(p.TechStack, ",")
		if len(tech) > 20 {
			tech = tech[:17] + "..."
		}

		lastCommit := formatAge(p.LastCommitDate)

		tbl.AddRow(
			p.Name,
			p.InferredType,
			StatusColor(p.Status),
			tech,
			lastCommit,
			p.CommitCount8M,
		)
	}

	tbl.Print()
	fmt.Printf("\n%s projects\n", Bold(fmt.Sprintf("%d", len(projects))))
}

func formatAge(dateStr string) string {
	if dateStr == "" {
		return "never"
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr[:10]
	}
	d := time.Since(t)
	switch {
	case d < 24*time.Hour:
		return "today"
	case d < 48*time.Hour:
		return "yesterday"
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dw ago", int(d.Hours()/(24*7)))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy ago", int(d.Hours()/(24*365)))
	}
}
