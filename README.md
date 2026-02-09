# prj — Local Git Repository Scanner & Project Dashboard for the Terminal

**Scan all your local git repos. See what you're working on, what's stale, and what tech each project uses — all from the command line.**

`prj` recursively finds every git repository on your machine, extracts rich metadata (tech stack, commit history, deployment config, project status), and gives you a fast, filterable overview. No web UI needed. No config files to write. Just point it at your dev folders and go.

Think of it as `git status` for your entire development life.

## Why prj?

If you're a developer with more than a handful of projects, you've probably hit these problems:

- **Lost track of what you were working on** — which repos have uncommitted changes? Which branch was I on?
- **Forgot to push before vacation** — is everything safely on the remote?
- **Can't remember the tech stack** of that project from 6 months ago
- **Side projects everywhere** — scattered across `~/Development`, `~/Projects`, `~/work`, with no overview
- **Context switching pain** — jumping between projects and losing track of where each one stands

`prj` gives you a single command to see all of it.

## Install

### Homebrew (macOS)

```bash
brew tap peeomid/tap
brew install prj
```

### Go Install

```bash
go install github.com/peeomid/prj@latest
```

### From Source

```bash
git clone https://github.com/peeomid/prj.git
cd prj
go install .
```

## Quick Start

```bash
# 1. Tell prj where your projects live
prj add ~/Development
prj add ~/Projects

# 2. Scan everything (finds repos recursively)
prj scan

# 3. See all your projects
prj list

# 4. Get a summary dashboard
prj status
```

That's it. `prj` recursively walks your folders, finds every `.git` repo, and extracts metadata automatically.

## Commands

### `prj add <folder>` — Register a folder to scan

```bash
prj add ~/Development         # Add your main dev folder
prj add ~/work/clients        # Add another folder
prj add .                     # Add current directory
```

Supports `~` expansion. The folder should be a **parent** that contains git repos inside it (at any depth).

### `prj scan` — Scan all folders and extract metadata

```bash
prj scan               # Scan and save results
prj scan --dry-run     # Preview what would be found (don't save)
```

Finds repos recursively. Skips `node_modules`, `vendor`, and hidden directories for speed. Extracts everything: git history, tech stack, deployment config, reference files, TODO counts.

### `prj list` — Show all projects in a table

```bash
prj list                          # All projects, sorted by last commit
prj list --status active          # Only active projects (committed in last 30 days)
prj list --status paused          # Stale projects (no commits in 90+ days)
prj list --tech react             # Filter by tech stack
prj list --type go-app            # Filter by project type
prj list --own                    # Only your own repos (exclude forks)
prj list --forks                  # Only forked repos
prj list --search api             # Search by name or path
prj list --sort commits           # Sort by commit count (most active first)
prj list --sort name              # Sort alphabetically
prj list --status active --own    # Combine filters
```

**Status levels:**
| Status | Meaning |
|--------|---------|
| `active` | Committed within 30 days |
| `wip` | Active + recent commit contains "WIP" |
| `recent` | Committed within 90 days |
| `paused` | No commits in 90+ days |

### `prj info <name>` — Full detail view for one project

```bash
prj info myapp           # Exact or partial name match
prj info api             # Finds "my-api-server", "api-gateway", etc.
```

Shows: description, tech stack, git history, recent commits, contributors, deployment methods, reference files, TODO counts, fork status, and more.

### `prj status` — Dashboard summary report

```bash
prj status
```

Shows:
- Total project count
- Breakdown by status (active / wip / recent / paused)
- Breakdown by type (rails-app, node-app, go-app, etc.)
- Own vs forked repos
- Top 5 most recently active projects
- Stalled projects (6+ months without a commit)

### `prj config` — View current settings

```bash
prj config
```

### `prj remove <folder>` — Stop scanning a folder

```bash
prj remove ~/old-projects
```

## What It Detects

### Tech Stack (auto-detected)

| Marker File | Detected As |
|-------------|-------------|
| `Gemfile` | Ruby (+ Rails if gem present) |
| `package.json` | Node (+ React / Next / Vue / TypeScript) |
| `requirements.txt` / `pyproject.toml` | Python |
| `go.mod` | Go |
| `Package.swift` / `*.xcodeproj` | Swift |
| `Cargo.toml` | Rust |

### Project Type (inferred)

`rails-app`, `next-app`, `react-app`, `vue-app`, `node-app`, `python-app`, `go-app`, `swift-app`, `rust-app`, `ruby-app`, `docs`, `script`, `unknown`

### Git Metadata

- Last commit date, message, and author
- 10 most recent commits
- Total commit count (last 8 months)
- All contributors
- Remote URL
- Fork detection (compares GitHub remote owner vs local git user)

### Deployment Detection

Dockerfile, docker-compose, Procfile (Heroku), fly.toml, vercel.json, netlify.toml, serverless.yml, GitHub Actions, CircleCI, deploy scripts in `bin/` or `package.json`

### Reference Files

Scans for: `README.md`, `CLAUDE.md`, `AGENT.md`, `CHANGELOG.md`, `TODO.md`, `.ai/`, `.cursor/`, `docs/`, `tasks/`

### Project Description

Extracted automatically with priority: `.ai/PROJECT_STATUS.md` > `CLAUDE.md` > `README.md` > folder name

### TODO Tracking

Counts open (`- [ ]`) and closed (`- [x]`) items in `TODO.md`

## Data Storage

Everything is stored locally as JSON:

```
~/.prj/
  config.json      # Tracked folders + settings
  projects.json    # All scanned project data
```

No database. No server. No cloud. Just files you can read, back up, or pipe into other tools.

## How It Compares

| Tool | What it does | How prj is different |
|------|-------------|---------------------|
| [mgitstatus](https://github.com/fboender/multi-git-status) | Shows git status across repos | prj adds tech stack, deployment, project type, dashboard |
| [gita](https://github.com/nosarthur/gita) | Manages multiple git repos | prj auto-discovers repos + extracts richer metadata |
| [mani](https://github.com/alajmo/mani) | Run commands across repos | prj focuses on scanning and reporting, not execution |
| [gh-dash](https://github.com/dlvhdr/gh-dash) | GitHub PR/issue dashboard | gh-dash is GitHub-only; prj works on local repos |

## Requirements

- Git (for repository scanning)
- macOS or Linux

## License

MIT
