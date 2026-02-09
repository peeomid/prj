package project

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DetectTechStack checks for language/framework markers in the repo.
func DetectTechStack(dir string) []string {
	var stack []string

	// Ruby / Rails
	if fileExists(dir, "Gemfile") {
		stack = append(stack, "ruby")
		if grepFile(filepath.Join(dir, "Gemfile"), "rails") {
			stack = append(stack, "rails")
		}
	}

	// Node / React / Next / Vue
	if fileExists(dir, "package.json") {
		stack = append(stack, "node")
		deps := readPackageJSONDeps(filepath.Join(dir, "package.json"))
		if _, ok := deps["react"]; ok {
			stack = append(stack, "react")
		}
		if _, ok := deps["next"]; ok {
			stack = append(stack, "next")
		}
		if _, ok := deps["vue"]; ok {
			stack = append(stack, "vue")
		}
		if _, ok := deps["typescript"]; ok {
			stack = append(stack, "typescript")
		}
	}

	// Python
	if fileExists(dir, "requirements.txt") || fileExists(dir, "pyproject.toml") || fileExists(dir, "setup.py") {
		stack = append(stack, "python")
	}

	// Go
	if fileExists(dir, "go.mod") {
		stack = append(stack, "go")
	}

	// Swift
	if fileExists(dir, "Package.swift") || dirHasGlob(dir, "*.xcodeproj") {
		stack = append(stack, "swift")
	}

	// Rust
	if fileExists(dir, "Cargo.toml") {
		stack = append(stack, "rust")
	}

	return stack
}

// InferType guesses the project type from tech stack and directory contents.
func InferType(dir string, stack []string) string {
	for _, t := range stack {
		switch t {
		case "rails":
			return "rails-app"
		case "next":
			return "next-app"
		case "react":
			return "react-app"
		case "vue":
			return "vue-app"
		}
	}
	for _, t := range stack {
		switch t {
		case "node":
			return "node-app"
		case "python":
			return "python-app"
		case "go":
			return "go-app"
		case "swift":
			return "swift-app"
		case "rust":
			return "rust-app"
		case "ruby":
			return "ruby-app"
		}
	}

	// Check for docs-only or scripts
	if dirHasGlob(dir, "*.md") && !fileExists(dir, "main.go") && !fileExists(dir, "index.js") {
		mdCount := countGlob(dir, "*.md")
		totalFiles := countGlob(dir, "*")
		if totalFiles > 0 && float64(mdCount)/float64(totalFiles) > 0.5 {
			return "docs"
		}
	}

	if dirHasGlob(dir, "*.sh") || dirHasGlob(dir, "*.py") {
		return "script"
	}

	return "unknown"
}

func fileExists(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

func dirHasGlob(dir, pattern string) bool {
	matches, _ := filepath.Glob(filepath.Join(dir, pattern))
	return len(matches) > 0
}

func countGlob(dir, pattern string) int {
	matches, _ := filepath.Glob(filepath.Join(dir, pattern))
	return len(matches)
}

func grepFile(path, substr string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return containsCI(string(data), substr)
}

func containsCI(s, substr string) bool {
	sl := len(s)
	subl := len(substr)
	if subl > sl {
		return false
	}
	for i := 0; i <= sl-subl; i++ {
		match := true
		for j := 0; j < subl; j++ {
			a, b := s[i+j], substr[j]
			if a >= 'A' && a <= 'Z' {
				a += 32
			}
			if b >= 'A' && b <= 'Z' {
				b += 32
			}
			if a != b {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func readPackageJSONDeps(path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var pkg struct {
		Dependencies    map[string]interface{} `json:"dependencies"`
		DevDependencies map[string]interface{} `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}
	merged := make(map[string]interface{})
	for k, v := range pkg.Dependencies {
		merged[k] = v
	}
	for k, v := range pkg.DevDependencies {
		merged[k] = v
	}
	return merged
}
