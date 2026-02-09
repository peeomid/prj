package project

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ExtractDescription returns (description, claudeDescription) from project files.
// Priority: .ai/PROJECT_STATUS.md → CLAUDE.md → README.md → fallback to dir name.
func ExtractDescription(dir string) (string, string) {
	var claudeDesc string

	// 1. Try .ai/PROJECT_STATUS.md
	if desc := extractFromFile(filepath.Join(dir, ".ai", "PROJECT_STATUS.md")); desc != "" {
		return desc, claudeDesc
	}

	// 2. Try CLAUDE.md
	claudePath := filepath.Join(dir, "CLAUDE.md")
	if desc := extractFromFile(claudePath); desc != "" {
		claudeDesc = desc
		return desc, claudeDesc
	}

	// 3. Try README.md
	if desc := extractFromFile(filepath.Join(dir, "README.md")); desc != "" {
		return desc, claudeDesc
	}

	// 4. Fallback
	return filepath.Base(dir), claudeDesc
}

// extractFromFile reads a markdown file and returns a short description.
// Looks for: first heading content, then first paragraph.
func extractFromFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var lines []string
	foundHeading := false
	pastHeading := false

	for s.Scan() {
		line := s.Text()

		// Skip empty lines before content
		if len(lines) == 0 && strings.TrimSpace(line) == "" {
			continue
		}

		// Skip badges/images at top
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "![") || strings.HasPrefix(trimmed, "<img") {
			continue
		}

		// Found a heading
		if strings.HasPrefix(trimmed, "#") {
			if foundHeading && pastHeading {
				// Hit second heading, stop
				break
			}
			foundHeading = true
			continue
		}

		if foundHeading {
			if strings.TrimSpace(line) == "" {
				if pastHeading && len(lines) > 0 {
					break
				}
				continue
			}
			pastHeading = true
			lines = append(lines, trimmed)
			// Grab up to 3 lines of first paragraph
			if len(lines) >= 3 {
				break
			}
		}
	}

	desc := strings.Join(lines, " ")
	// Truncate to ~200 chars
	if len(desc) > 200 {
		desc = desc[:197] + "..."
	}
	return desc
}
