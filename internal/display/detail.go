package display

import (
	"fmt"
	"strings"

	"github.com/Osimify/prj/internal/project"
)

// PrintDetail renders a full detail view for a single project.
func PrintDetail(p *project.Project) {
	fmt.Println()
	fmt.Printf("  %s  %s\n", Bold(p.Name), StatusColor(p.Status))
	fmt.Printf("  %s\n", Gray(p.Path))
	fmt.Println()

	if p.Description != "" {
		fmt.Printf("  %s\n\n", p.Description)
	}

	section("Type", p.InferredType)
	section("Tech", strings.Join(p.TechStack, ", "))
	section("Remote", p.GitRemote)

	if p.IsFork {
		section("Fork", "yes")
	}

	fmt.Println()
	fmt.Printf("  %s\n", Bold("Git"))
	section("  Last commit", fmt.Sprintf("%s — %s (%s)", formatAge(p.LastCommitDate), p.LastCommitMessage, p.LastCommitAuthor))
	section("  Commits (8m)", fmt.Sprintf("%d", p.CommitCount8M))
	section("  Contributors", strings.Join(p.Contributors, ", "))

	if len(p.RecentCommits) > 0 {
		fmt.Printf("\n  %s\n", Bold("Recent Commits"))
		for _, c := range p.RecentCommits {
			fmt.Printf("    %s %s %s\n", Gray(c.Date[:10]), Cyan(c.Hash[:7]), c.Message)
		}
	}

	if p.TodoOpen > 0 || p.TodoClosed > 0 {
		fmt.Printf("\n  %s  open:%s  closed:%s\n", Bold("TODOs"), Green(fmt.Sprintf("%d", p.TodoOpen)), Gray(fmt.Sprintf("%d", p.TodoClosed)))
	}

	if len(p.Deployment) > 0 {
		fmt.Printf("\n  %s  %s\n", Bold("Deployment"), strings.Join(p.Deployment, ", "))
	}

	if len(p.ReferenceFiles.Root) > 0 {
		fmt.Printf("\n  %s\n", Bold("Reference Files"))
		printRefList("Root", p.ReferenceFiles.Root)
		printRefList("AI", p.ReferenceFiles.AI)
		printRefList("Cursor", p.ReferenceFiles.Cursor)
		printRefList("Docs", p.ReferenceFiles.Docs)
		printRefList("Tasks", p.ReferenceFiles.Tasks)
	}

	if len(p.NestedRepos) > 0 {
		fmt.Printf("\n  %s  %s\n", Bold("Nested Repos"), strings.Join(p.NestedRepos, ", "))
	}

	if len(p.Errors) > 0 {
		fmt.Printf("\n  %s\n", Red("Errors"))
		for _, e := range p.Errors {
			fmt.Printf("    %s\n", Red(e))
		}
	}

	fmt.Printf("\n  %s %s\n\n", Gray("Scanned:"), Gray(p.ScannedAt))
}

func section(label, value string) {
	if value == "" {
		return
	}
	fmt.Printf("  %-16s %s\n", Bold(label), value)
}

func printRefList(label string, files []string) {
	if len(files) == 0 {
		return
	}
	fmt.Printf("    %-10s %s\n", label+":", strings.Join(files, ", "))
}
