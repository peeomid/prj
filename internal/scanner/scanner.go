package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// skipDirs are directory names we never descend into.
var skipDirs = map[string]bool{
	"node_modules": true,
	"vendor":       true,
	".git":         true,
	"__pycache__":  true,
	".bundle":      true,
}

// FindRepos walks a directory tree recursively and returns paths that
// contain a .git folder. Once a repo is found, we don't descend into
// it (no nested repo scanning). Hidden dirs, node_modules, vendor, etc.
// are skipped for speed.
func FindRepos(root string) ([]string, error) {
	var repos []string
	found := map[string]bool{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip dirs we can't read
		}

		if !d.IsDir() {
			return nil
		}

		name := d.Name()

		// Skip hidden dirs (except the root itself)
		if strings.HasPrefix(name, ".") && path != root {
			return fs.SkipDir
		}

		// Skip known heavy dirs
		if skipDirs[name] {
			return fs.SkipDir
		}

		// If parent is already a found repo, skip (no nested repos)
		for repo := range found {
			if strings.HasPrefix(path, repo+string(os.PathSeparator)) {
				return fs.SkipDir
			}
		}

		// Check if this dir is a git repo
		gitDir := filepath.Join(path, ".git")
		if info, statErr := os.Stat(gitDir); statErr == nil && info.IsDir() {
			repos = append(repos, path)
			found[path] = true
			return fs.SkipDir // don't descend into the repo
		}

		return nil
	})

	return repos, err
}
