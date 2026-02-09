package project

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Osimify/prj/internal/scanner"
)

// InferState determines project status based on commit recency, TODOs, and commit messages.
// Returns (status, todoOpen, todoClosed).
func InferState(dir, lastCommitDate string, commits []scanner.CommitInfo) (string, int, int) {
	todoOpen, todoClosed := parseTodoFile(filepath.Join(dir, "TODO.md"))

	// Check commit messages for WIP keywords
	hasWIP := false
	for _, c := range commits {
		msg := strings.ToLower(c.Message)
		if strings.Contains(msg, "wip") || strings.Contains(msg, "work in progress") {
			hasWIP = true
			break
		}
	}

	if hasWIP && isRecent(lastCommitDate, 30) {
		return "wip", todoOpen, todoClosed
	}

	if isRecent(lastCommitDate, 30) {
		return "active", todoOpen, todoClosed
	}

	if isRecent(lastCommitDate, 90) {
		return "recent", todoOpen, todoClosed
	}

	return "paused", todoOpen, todoClosed
}

func isRecent(dateStr string, days int) bool {
	if dateStr == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return false
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	return t.After(cutoff)
}

func parseTodoFile(path string) (int, int) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0
	}
	defer f.Close()

	open, closed := 0, 0
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "- [ ]") || strings.HasPrefix(line, "* [ ]") {
			open++
		} else if strings.HasPrefix(line, "- [x]") || strings.HasPrefix(line, "- [X]") ||
			strings.HasPrefix(line, "* [x]") || strings.HasPrefix(line, "* [X]") {
			closed++
		}
	}
	return open, closed
}
