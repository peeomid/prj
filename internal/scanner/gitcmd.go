package scanner

import (
	"os/exec"
	"strings"
)

// Git runs a git command in the given directory and returns trimmed output.
func Git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// GitLines runs a git command and returns output split by newlines.
func GitLines(dir string, args ...string) ([]string, error) {
	out, err := Git(dir, args...)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

type CommitInfo struct {
	Hash    string `json:"hash"`
	Date    string `json:"date"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

// RecentCommits returns the last n commits.
func RecentCommits(dir string, n int) ([]CommitInfo, error) {
	lines, err := GitLines(dir, "log", "--format=%H|%aI|%an|%s", "-n", itoa(n))
	if err != nil {
		return nil, err
	}
	var commits []CommitInfo
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 4 {
			continue
		}
		commits = append(commits, CommitInfo{
			Hash:    parts[0],
			Date:    parts[1],
			Author:  parts[2],
			Message: parts[3],
		})
	}
	return commits, nil
}

// CommitCountSince returns number of commits since a date string (e.g. "8 months ago").
func CommitCountSince(dir, since string) int {
	out, err := Git(dir, "rev-list", "--count", "--since="+since, "HEAD")
	if err != nil {
		return 0
	}
	n := 0
	for _, c := range out {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

// Contributors returns unique author names.
func Contributors(dir string) []string {
	lines, err := GitLines(dir, "log", "--format=%an")
	if err != nil {
		return nil
	}
	seen := map[string]bool{}
	var result []string
	for _, name := range lines {
		name = strings.TrimSpace(name)
		if name != "" && !seen[name] {
			seen[name] = true
			result = append(result, name)
		}
	}
	return result
}

// Remote returns the origin remote URL.
func Remote(dir string) string {
	out, _ := Git(dir, "remote", "get-url", "origin")
	return out
}

// GitUserName returns the local or global git user.name.
func GitUserName(dir string) string {
	out, _ := Git(dir, "config", "user.name")
	return out
}

// GitHubUser returns the global github.user config.
func GitHubUser(dir string) string {
	out, _ := Git(dir, "config", "--global", "github.user")
	return out
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
