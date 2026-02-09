# prj — Go CLI for Git Project Scanning & Reporting

## Context

Luan wants a CLI version of the `your-project-dashboard` Rails app — same scanning/metadata but no web dashboard. Just scan, track, and report from terminal. Language: Go. Storage: JSON files. Global install via `go install`.

## Project Structure

```
prj/
  go.mod                           # module github.com/Osimify/prj
  main.go                          # entry point → cmd.Execute()
  README.md                        # project description

  cmd/
    root.go                        # cobra root + version
    add.go                         # prj add <folder>
    remove.go                      # prj remove <folder>
    scan.go                        # prj scan [--dry-run]
    list.go                        # prj list [--status, --type, --tech, --own, --forks, --search, --sort]
    info.go                        # prj info <name>
    status.go                      # prj status (summary report)
    config_cmd.go                  # prj config (show config)

  internal/
    config/config.go               # load/save ~/.prj/config.json
    scanner/
      scanner.go                   # walk dirs, find .git repos
      gitcmd.go                    # run git commands, parse output
    project/
      project.go                   # Project struct + ExtractFromPath()
      techstack.go                 # detect Ruby/Node/Python/Go/Swift
      description.go               # extract from README/CLAUDE.md etc
      state.go                     # infer active/paused/wip from commits+TODOs
      deployment.go                # detect Dockerfile/Procfile/deploy scripts
      references.go                # find reference files (root/ai/cursor/docs/tasks)
    store/store.go                 # load/save/merge ~/.prj/projects.json
    display/
      colors.go                   # status colors (green/blue/yellow/gray)
      table.go                    # table formatter for list
      detail.go                   # single project detail view
      summary.go                  # summary report (prj status)
```

## Data Storage

- `~/.prj/config.json` — folders list + cutoff_days (default 240)
- `~/.prj/projects.json` — flat array of project objects

## CLI Commands

| Command | What it does |
|---------|-------------|
| `prj add <folder>` | Add folder to scan list (expands ~, validates exists) |
| `prj remove <folder>` | Remove folder from scan list |
| `prj scan` | Scan all folders, upsert into projects.json |
| `prj scan --dry-run` | Scan without saving |
| `prj list` | Table of all projects (filterable) |
| `prj list --status active` | Filter: active/recent/paused/wip |
| `prj list --type node-app` | Filter by inferred type |
| `prj list --tech react` | Filter by tech stack |
| `prj list --own` / `--forks` | Filter by ownership |
| `prj list --search <q>` | Search name/path |
| `prj list --sort name\|date\|commits` | Sort (default: date desc) |
| `prj info <name>` | Detailed view of one project |
| `prj status` | Summary: counts by type, status, ownership, top recent, stalled |
| `prj config` | Show current config |

## Metadata Extracted Per Repo

Exact port from reference project (`project_data.rb`):

- **Git data:** last_commit_date, last_commit_message, last_commit_author, recent_commits (10), commit_count_8m, contributors, git_remote
- **Fork detection:** compare GitHub remote owner vs local git user.name / github.user
- **Tech stack:** Gemfile→ruby/rails, package.json→node/react/next/vue, requirements.txt/pyproject.toml→python, go.mod→go, Package.swift→swift
- **Inferred type:** rails-app, node-app, python-app, go-app, docs, script, unknown
- **Reference files:** root (README.md, CLAUDE.md, AGENT.md, CHANGELOG.md), .ai/, .cursor/, docs/, tasks/
- **Description:** priority chain .ai/PROJECT_STATUS.md → CLAUDE.md → README.md (section headers → first paragraph → fallback to dir name)
- **Current state:** TODO.md open/closed counts + commit recency + commit message keywords (WIP/done)
- **Deployment:** Dockerfile, docker-compose.yml, Procfile, bin/deploy, package.json deploy script, README mentions
- **Other:** nested_repos, plans_count, ai_docs_count, claude_description, errors array

## Dependencies

- `github.com/spf13/cobra` — CLI framework (subcommands, flags, help)
- `github.com/fatih/color` — terminal colors
- `github.com/rodaine/table` — table printing

Everything else from Go stdlib: `os/exec` for git, `encoding/json`, `filepath.WalkDir`, `bufio`, `regexp`.

## Build Order

**Phase 1 — Foundation**
1. `go.mod` + `main.go`
2. `internal/config/config.go`
3. `cmd/root.go` + `cmd/config_cmd.go`
4. `cmd/add.go` + `cmd/remove.go`

**Phase 2 — Core Scanning**
5. `internal/scanner/gitcmd.go`
6. `internal/project/project.go` (struct + ExtractFromPath)
7. `internal/project/techstack.go`
8. `internal/project/description.go`
9. `internal/project/state.go`
10. `internal/project/deployment.go`
11. `internal/project/references.go`
12. `internal/scanner/scanner.go`
13. `internal/store/store.go`

**Phase 3 — CLI Commands + Display**
14. `internal/display/colors.go`
15. `cmd/scan.go` + progress output
16. `internal/display/table.go` + `cmd/list.go`
17. `internal/display/detail.go` + `cmd/info.go`
18. `internal/display/summary.go` + `cmd/status.go`

**Phase 4 — README + install test**
19. `README.md`
20. `go build && go install` verification

## Install

```bash
# dev
cd ~/Development/Osimify/prj && go install .

# users
go install github.com/Osimify/prj@latest
```

Single binary in `$GOPATH/bin/prj`.

## Verification

1. `prj config` — should show empty folder list
2. `prj add ~/Development` — should add folder
3. `prj scan` — should find repos, print progress, save to ~/.prj/projects.json
4. `prj list` — should show table of projects
5. `prj list --status active --own` — should filter correctly
6. `prj info <project>` — should show full detail
7. `prj status` — should show summary report
8. `prj scan --dry-run` — should scan but not save
