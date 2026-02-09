package display

import (
	"fmt"
	"sort"
	"time"

	"github.com/peeomid/prj/internal/project"
)

// PrintSummary renders a status report with counts and highlights.
func PrintSummary(projects []*project.Project) {
	if len(projects) == 0 {
		fmt.Println("No projects scanned yet. Run: prj scan")
		return
	}

	fmt.Printf("\n  %s  %d projects\n\n", Bold("Project Summary"), len(projects))

	// By status
	statusCounts := map[string]int{}
	for _, p := range projects {
		statusCounts[p.Status]++
	}
	fmt.Printf("  %s\n", Bold("By Status"))
	for _, s := range []string{"active", "wip", "recent", "paused"} {
		if c, ok := statusCounts[s]; ok {
			fmt.Printf("    %s  %d\n", StatusColor(s), c)
		}
	}

	// By type
	typeCounts := map[string]int{}
	for _, p := range projects {
		typeCounts[p.InferredType]++
	}
	fmt.Printf("\n  %s\n", Bold("By Type"))
	types := sortedKeys(typeCounts)
	for _, t := range types {
		fmt.Printf("    %-15s %d\n", t, typeCounts[t])
	}

	// Ownership
	ownCount, forkCount := 0, 0
	for _, p := range projects {
		if p.IsFork {
			forkCount++
		} else {
			ownCount++
		}
	}
	fmt.Printf("\n  %s\n", Bold("Ownership"))
	fmt.Printf("    Own:    %d\n", ownCount)
	fmt.Printf("    Forks:  %d\n", forkCount)

	// Top 5 most recent
	sorted := make([]*project.Project, len(projects))
	copy(sorted, projects)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].LastCommitDate > sorted[j].LastCommitDate
	})

	fmt.Printf("\n  %s\n", Bold("Most Recent"))
	limit := 5
	if len(sorted) < limit {
		limit = len(sorted)
	}
	for _, p := range sorted[:limit] {
		fmt.Printf("    %-25s %s  %s\n", p.Name, StatusColor(p.Status), formatAge(p.LastCommitDate))
	}

	// Stalled (no commits in 6+ months)
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	var stalled []*project.Project
	for _, p := range projects {
		if p.LastCommitDate == "" {
			stalled = append(stalled, p)
			continue
		}
		t, err := time.Parse(time.RFC3339, p.LastCommitDate)
		if err == nil && t.Before(sixMonthsAgo) {
			stalled = append(stalled, p)
		}
	}
	if len(stalled) > 0 {
		fmt.Printf("\n  %s  (%d)\n", Bold("Stalled (6+ months)"), len(stalled))
		limit := 10
		if len(stalled) < limit {
			limit = len(stalled)
		}
		for _, p := range stalled[:limit] {
			fmt.Printf("    %s  %s\n", Gray(p.Name), Gray(formatAge(p.LastCommitDate)))
		}
	}

	fmt.Println()
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
