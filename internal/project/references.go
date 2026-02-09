package project

import (
	"os"
	"path/filepath"
)

// FindReferences scans for reference files across the repo.
func FindReferences(dir string) ReferenceFiles {
	ref := ReferenceFiles{}

	// Root reference files
	rootFiles := []string{
		"README.md", "CLAUDE.md", "AGENT.md", "CHANGELOG.md",
		"TODO.md", "CONTRIBUTING.md", "LICENSE",
	}
	for _, name := range rootFiles {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			ref.Root = append(ref.Root, name)
		}
	}

	// .ai/ directory
	ref.AI = listDir(filepath.Join(dir, ".ai"))

	// .cursor/ directory
	ref.Cursor = listDir(filepath.Join(dir, ".cursor"))

	// docs/ directory
	ref.Docs = listDir(filepath.Join(dir, "docs"))

	// tasks/ directory
	ref.Tasks = listDir(filepath.Join(dir, "tasks"))

	return ref
}

func listDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.Name()[0] == '.' {
			continue
		}
		names = append(names, e.Name())
	}
	return names
}
