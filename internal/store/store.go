package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/peeomid/prj/internal/config"
	"github.com/peeomid/prj/internal/project"
)

func Path() string {
	return filepath.Join(config.Dir(), "projects.json")
}

// Load reads the projects from disk.
func Load() ([]*project.Project, error) {
	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var projects []*project.Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

// Save writes projects to disk.
func Save(projects []*project.Project) error {
	if err := os.MkdirAll(config.Dir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0644)
}

// Merge upserts scanned projects into existing ones (by path).
func Merge(existing, scanned []*project.Project) []*project.Project {
	byPath := make(map[string]*project.Project)
	for _, p := range existing {
		byPath[p.Path] = p
	}
	for _, p := range scanned {
		byPath[p.Path] = p
	}
	result := make([]*project.Project, 0, len(byPath))
	for _, p := range byPath {
		result = append(result, p)
	}
	return result
}
