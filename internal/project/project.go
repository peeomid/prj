package project

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/Osimify/prj/internal/scanner"
)

type Project struct {
	Name              string              `json:"name"`
	Path              string              `json:"path"`
	Description       string              `json:"description"`
	ClaudeDescription string              `json:"claude_description,omitempty"`
	TechStack         []string            `json:"tech_stack"`
	InferredType      string              `json:"inferred_type"`
	Status            string              `json:"status"`
	IsFork            bool                `json:"is_fork"`
	GitRemote         string              `json:"git_remote"`
	LastCommitDate    string              `json:"last_commit_date"`
	LastCommitMessage string              `json:"last_commit_message"`
	LastCommitAuthor  string              `json:"last_commit_author"`
	RecentCommits     []scanner.CommitInfo `json:"recent_commits"`
	CommitCount8M     int                 `json:"commit_count_8m"`
	Contributors      []string            `json:"contributors"`
	ReferenceFiles    ReferenceFiles       `json:"reference_files"`
	TodoOpen          int                 `json:"todo_open"`
	TodoClosed        int                 `json:"todo_closed"`
	Deployment        []string            `json:"deployment"`
	NestedRepos       []string            `json:"nested_repos,omitempty"`
	PlansCount        int                 `json:"plans_count"`
	AIDocsCount       int                 `json:"ai_docs_count"`
	Errors            []string            `json:"errors,omitempty"`
	ScannedAt         string              `json:"scanned_at"`
}

type ReferenceFiles struct {
	Root   []string `json:"root"`
	AI     []string `json:"ai"`
	Cursor []string `json:"cursor"`
	Docs   []string `json:"docs"`
	Tasks  []string `json:"tasks"`
}

// ExtractFromPath scans a git repo at the given path and returns a Project.
func ExtractFromPath(repoPath string) *Project {
	p := &Project{
		Name:      filepath.Base(repoPath),
		Path:      repoPath,
		ScannedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Git data
	commits, err := scanner.RecentCommits(repoPath, 10)
	if err != nil {
		p.Errors = append(p.Errors, "git log: "+err.Error())
	} else {
		p.RecentCommits = commits
		if len(commits) > 0 {
			p.LastCommitDate = commits[0].Date
			p.LastCommitMessage = commits[0].Message
			p.LastCommitAuthor = commits[0].Author
		}
	}

	p.CommitCount8M = scanner.CommitCountSince(repoPath, "8 months ago")
	p.Contributors = scanner.Contributors(repoPath)
	p.GitRemote = scanner.Remote(repoPath)

	// Fork detection
	p.IsFork = detectFork(repoPath, p.GitRemote)

	// Tech stack + type
	p.TechStack = DetectTechStack(repoPath)
	p.InferredType = InferType(repoPath, p.TechStack)

	// Reference files
	p.ReferenceFiles = FindReferences(repoPath)
	p.AIDocsCount = len(p.ReferenceFiles.AI)
	p.PlansCount = len(p.ReferenceFiles.Tasks)

	// Description
	p.Description, p.ClaudeDescription = ExtractDescription(repoPath)

	// State
	p.Status, p.TodoOpen, p.TodoClosed = InferState(repoPath, p.LastCommitDate, p.RecentCommits)

	// Deployment
	p.Deployment = DetectDeployment(repoPath)

	// Nested repos
	p.NestedRepos = findNestedRepos(repoPath)

	return p
}

func detectFork(dir, remote string) bool {
	if remote == "" {
		return false
	}

	// Extract owner from GitHub remote
	owner := extractGitHubOwner(remote)
	if owner == "" {
		return false
	}

	localUser := scanner.GitUserName(dir)
	ghUser := scanner.GitHubUser(dir)

	ownerLower := strings.ToLower(owner)
	if localUser != "" && strings.ToLower(localUser) == ownerLower {
		return false
	}
	if ghUser != "" && strings.ToLower(ghUser) == ownerLower {
		return false
	}

	return true
}

func extractGitHubOwner(remote string) string {
	// SSH: git@github.com:owner/repo.git
	if strings.Contains(remote, "github.com:") {
		parts := strings.SplitN(remote, ":", 2)
		if len(parts) == 2 {
			path := strings.TrimSuffix(parts[1], ".git")
			segments := strings.Split(path, "/")
			if len(segments) >= 1 {
				return segments[0]
			}
		}
	}
	// HTTPS: https://github.com/owner/repo.git
	if strings.Contains(remote, "github.com/") {
		idx := strings.Index(remote, "github.com/")
		path := remote[idx+len("github.com/"):]
		path = strings.TrimSuffix(path, ".git")
		segments := strings.Split(path, "/")
		if len(segments) >= 1 {
			return segments[0]
		}
	}
	return ""
}

func findNestedRepos(dir string) []string {
	var nested []string
	entries, err := filepath.Glob(filepath.Join(dir, "*", ".git"))
	if err != nil {
		return nil
	}
	for _, e := range entries {
		nested = append(nested, filepath.Dir(e))
	}
	return nested
}
