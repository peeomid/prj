package project

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DetectDeployment checks for deployment-related files and configs.
func DetectDeployment(dir string) []string {
	var result []string

	checks := []struct {
		path  string
		label string
	}{
		{"Dockerfile", "docker"},
		{"docker-compose.yml", "docker-compose"},
		{"docker-compose.yaml", "docker-compose"},
		{"Procfile", "heroku"},
		{"fly.toml", "fly.io"},
		{"vercel.json", "vercel"},
		{"netlify.toml", "netlify"},
		{"app.yaml", "gcp"},
		{"serverless.yml", "serverless"},
		{"serverless.yaml", "serverless"},
		{filepath.Join("bin", "deploy"), "deploy-script"},
		{".github/workflows", "github-actions"},
		{".circleci", "circleci"},
	}

	for _, c := range checks {
		if _, err := os.Stat(filepath.Join(dir, c.path)); err == nil {
			result = append(result, c.label)
		}
	}

	// Check package.json for deploy script
	pkgPath := filepath.Join(dir, "package.json")
	if data, err := os.ReadFile(pkgPath); err == nil {
		var pkg struct {
			Scripts map[string]string `json:"scripts"`
		}
		if json.Unmarshal(data, &pkg) == nil {
			if _, ok := pkg.Scripts["deploy"]; ok {
				result = append(result, "npm-deploy")
			}
		}
	}

	// Check README for deployment mentions
	readmePath := filepath.Join(dir, "README.md")
	if grepFile(readmePath, "deploy") {
		result = append(result, "readme-deploy")
	}

	return result
}
