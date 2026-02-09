# prj

A fast CLI tool that scans your local git repositories, extracts metadata, and reports their status — all from the terminal.

## What it does

You tell `prj` which folders contain your projects (e.g. `~/Development`). It recursively finds every git repo under those folders, pulls out useful info (tech stack, last commit, activity level, deployment status, docs inventory), and saves it locally as JSON.

Then you can list, filter, search, and inspect your projects without leaving the terminal.

## Why

When you have dozens of projects scattered across folders, it's hard to remember what exists, what's active, and what's stalled. `prj` gives you a bird's-eye view of all your work.

## How it works

1. **Add folders:** `prj add ~/Development` — tells prj where to look
2. **Scan:** `prj scan` — finds git repos, extracts metadata, saves to `~/.prj/projects.json`
3. **View:** `prj list`, `prj info <name>`, `prj status` — browse and filter your projects

## Tech

- **Go** — single binary, no runtime needed
- **Storage:** JSON files at `~/.prj/` (config + project data)
- **No database, no server** — pure CLI
- **Install:** `go install github.com/Osimify/prj@latest`

## Inspired by

[your-project-dashboard](https://github.com/aviflombaum/your-project-dashboard) — a Rails app that does similar scanning with a web dashboard. `prj` is the CLI-only version.

## Commands

```
prj add <folder>       Add a folder to scan
prj remove <folder>    Remove a folder
prj scan               Scan all folders, save results
prj scan --dry-run     Scan without saving
prj list               List all projects (with filters)
prj info <name>        Detailed view of one project
prj status             Summary report
prj config             Show current configuration
```

## Essential commands

```bash
go build -o prj .      # Build
go install .           # Install globally
go test ./...          # Run tests
```
