# dwm — DBeaver Workspace Manager

> CLI + MCP Server for automating DBeaver workspace configurations, connections, SQL scripts, and snippets.

[![CI](https://github.com/jinkp/dbeaver-go-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/jinkp/dbeaver-go-mcp/actions/workflows/ci.yml)
[![Release](https://github.com/jinkp/dbeaver-go-mcp/actions/workflows/release.yml/badge.svg)](https://github.com/jinkp/dbeaver-go-mcp/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/jinkp/dbeaver-go-mcp)](go.mod)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## What it does

`dwm` eliminates the repetitive manual setup of DBeaver environments. Instead of clicking through DBeaver's UI to create dozens of connections, scripts, and snippets for every project and environment, you describe your setup in a single YAML file and let `dwm` generate everything — including direct installation into DBeaver.

It also exposes all its functionality as an **MCP server**, so AI assistants (OpenCode, Claude Code) can orchestrate your DBeaver workspace generation directly.

---

## Install

### Windows (PowerShell)

```powershell
iex (irm https://raw.githubusercontent.com/jinkp/dbeaver-go-mcp/main/install.ps1)
```

### macOS / Linux

```sh
curl -fsSL https://raw.githubusercontent.com/jinkp/dbeaver-go-mcp/main/install.sh | sh
```

### Homebrew (macOS / Linux)

```sh
brew install jinkp/tap/dwm
```

### Go install

```sh
go install github.com/jinkp/dbeaver-go-mcp/cmd/dwm@latest
```

---

## Quick start

**1. Create a config file** (`myproject.yaml`):

```yaml
project: SRPayments

environments:
  - DEV
  - QA
  - PROD

connections:
  - name: "{{.Project}} {{.Env}} - PostgreSQL"
    type: postgresql
    host: "{{.Env | lower}}-postgres.internal"
    port: 5432
    database: srpayments
    folder: "{{.Env}}"

  - name: "{{.Project}} {{.Env}} - SQL Server"
    type: sqlserver
    host: "{{.Env | lower}}-sqlserver.internal"
    port: 1433
    database: srpayments
    folder: "{{.Env}}"

engines:
  - postgresql
  - sqlserver
```

**2. Generate everything:**

```sh
# Create workspace structure
dwm create-workspace SRPayments

# Generate connection files (DBeaver 26.0 format, no credentials)
dwm create-connections myproject.yaml

# Generate SQL maintenance scripts for all engines
dwm generate-scripts myproject.yaml

# Generate organized SQL snippets
dwm sync-snippets myproject.yaml
```

**3. Install directly into DBeaver** (DBeaver must be closed):

```sh
dwm install SRPayments
```

Or import manually: DBeaver → **File** → **Import** → **DBeaver Projects** → select `./output/SRPayments/`

---

## Commands

### `create-workspace`

```sh
dwm create-workspace <name> [--environments DEV,QA,STAGE,PROD] [--output ./output]
```

Generates the DBeaver project folder structure:

```
output/SRPayments/
  .dbeaver/
  Scripts/
  Diagrams/
  Bookmarks/
```

**Example:**
```sh
dwm create-workspace SRPayments --environments DEV,QA,PROD
```

---

### `create-project`

```sh
dwm create-project <name> --workspace <ws-name> [--output ./output]
```

Creates a project folder nested inside an existing workspace.

**Example:**
```sh
dwm create-project Analytics --workspace SRPayments
```

---

### `create-connections`

```sh
dwm create-connections config.yaml [--split-by-env] [--dry-run] [--output ./output]
```

Generates `data-sources.json` in DBeaver 26.0 format. Connections are **never stored with credentials** (`save-password: false`).

- `--split-by-env` — generates one `data-sources-<env>.json` per environment (DBeaver loads all of them)
- `--dry-run` — prints what would be generated without writing files

**Stable connection IDs:** same config always generates the same connection IDs — safe to version-control the output.

**Example:**
```sh
# Generate all environments in one file
dwm create-connections myproject.yaml

# Split into per-environment files (DEV/QA/PROD)
dwm create-connections myproject.yaml --split-by-env

# Preview without writing
dwm create-connections myproject.yaml --dry-run
```

---

### `generate-scripts`

```sh
dwm generate-scripts config.yaml [--scripts health-check,deadlocks-check] [--output ./output]
```

Generates SQL maintenance scripts under `Scripts/scripts/<engine>/<name>.sql`.

**Available scripts:**

| Script | Description |
|--------|-------------|
| `health-check` | DB version, current user, server time, DB size |
| `deadlocks-check` | Blocked queries and sessions |
| `long-running-queries` | Queries running over 30 seconds |
| `table-size-analysis` | Top tables by disk size |
| `compare-row-counts` | Estimated row counts per table |
| `migration-validation` | Schema/table listing for migration verification |
| `failed-jobs` | Failed scheduled jobs (SQL Server Agent / pg_cron) |

**Example:**
```sh
# Generate all scripts for all engines in config
dwm generate-scripts myproject.yaml

# Generate only specific scripts
dwm generate-scripts myproject.yaml --scripts health-check,deadlocks-check
```

---

### `sync-snippets`

```sh
dwm sync-snippets config.yaml [--engines postgresql,sqlserver,mysql,shared] [--output ./output]
```

Generates reusable SQL snippets under `Scripts/snippets/<engine>/<name>.sql`.

**Included snippets per engine:**
- `postgresql` — list-tables, list-indexes, table-row-count, active-connections
- `sqlserver` — list-tables, list-indexes, table-row-count, active-connections
- `mysql` — list-tables, table-row-count, active-connections
- `shared` — compare-row-counts, migration-validation

**Example:**
```sh
dwm sync-snippets myproject.yaml
dwm sync-snippets myproject.yaml --engines postgresql,shared
```

---

### `create-template`

```sh
dwm create-template <name> [--output ./output]
```

Scaffolds a reusable project template in `./output/templates/<name>/template.yaml`.

**Example:**
```sh
dwm create-template payments-multi-env
```

---

### `apply-template`

```sh
dwm apply-template <name> [--set Key=Value] [--output ./output]
```

Applies a named template (built-in or local) with variable substitution.

**Built-in templates:**
- `multi-env` — Generic multi-environment project (workspace + connections + scripts + snippets)

**Example:**
```sh
# Apply built-in multi-env template
dwm apply-template multi-env --set Project=Analytics --set DBEngine=postgresql

# Apply with multiple environments
dwm apply-template multi-env \
  --set Project=ControlInterno \
  --set Environments=DEV,QA,PROD \
  --set DBEngine=sqlserver
```

---

### `install`

```sh
dwm install <workspace-name> [flags]
```

Copies the generated output directly into DBeaver's workspace. **DBeaver must be closed** before running.

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--output` | `./output` | Source directory |
| `--workspace-path` | _(auto)_ | Override DBeaver workspace path |
| `--overwrite` | `false` | Replace files instead of merging |
| `--dry-run` | `false` | Preview without writing |
| `--force` | `false` | Skip DBeaver-is-running check |
| `--backup` | `false` | Backup target workspace before writing |

**Auto-detected DBeaver paths:**

| OS | Path |
|----|------|
| Windows | `%APPDATA%\DBeaverData\workspace6\` |
| macOS | `~/Library/DBeaverData/workspace6/` |
| Linux | `$XDG_DATA_HOME/DBeaverData/workspace6/` |

**Safety rules:**
- `.dbeaver/credentials-config.json` is **NEVER** modified
- `data-sources.json` is **merged by default** — your existing connections are preserved
- DBeaver running → error (use `--force` to bypass, but DBeaver may overwrite on close)

**Examples:**
```sh
# Preview install plan
dwm install SRPayments --dry-run

# Install (DBeaver must be closed)
dwm install SRPayments

# Install with backup (recommended when using --overwrite)
dwm install SRPayments --overwrite --backup

# Non-standard DBeaver path
dwm install SRPayments --workspace-path "D:\DBeaver\workspace6"
```

---

## MCP Server

`dwm` exposes all its commands as MCP tools so AI assistants can orchestrate your DBeaver setup directly.

### Setup (interactive wizard)

```sh
# Register with OpenCode
dwm setup opencode

# Register with Claude Code
dwm setup claude
```

The wizard guides you through global vs local scope selection. For CI/scripting, use flags directly:

```sh
dwm setup opencode --global   # ~/.config/opencode/opencode.json
dwm setup opencode --local    # ./opencode.json
dwm setup claude --global     # ~/.claude.json
dwm setup claude --local      # ./.claude/settings.json
```

### Start MCP server (manual)

```sh
dwm mcp
```

### Available MCP tools

| Tool | Required args | Optional args |
|------|--------------|--------------|
| `create_workspace` | `name` | `environments`, `output`, `conflict` |
| `create_project` | `name`, `workspace` | `output`, `conflict` |
| `create_connections` | `config` (YAML path) | `output`, `split_by_env`, `conflict` |
| `generate_scripts` | `config` (YAML path) | `scripts`, `output`, `conflict` |
| `sync_snippets` | `config` (YAML path) | `engines`, `output`, `conflict` |
| `create_template` | `name` | `output` |
| `apply_template` | `name` | `vars` (object), `output`, `conflict` |

All tools accept `conflict: "error" | "overwrite" | "skip"` (default: `"error"`).

### Example AI interactions

Once registered, an AI assistant can handle natural language requests:

> *"Create a DBeaver workspace for SRPayments with DEV, QA, and PROD environments using PostgreSQL and SQL Server"*

The AI will:
1. Write a YAML config file
2. Call `create_workspace` with `name="SRPayments"`
3. Call `create_connections` with the config path
4. Call `generate_scripts` and `sync_snippets`

> *"Generate all maintenance scripts for my Analytics project and install them into DBeaver"*

The AI will call `generate_scripts` then suggest `dwm install Analytics` (or call it via the tool).

---

## Global flags

```
--output string    Output directory (default: "./output")
--overwrite        Replace existing files
--skip             Skip existing files (warn and continue)
--dry-run          Print plan without writing
--merge            Not supported in v1 (returns error)
```

---

## Supported database engines

| Engine | Provider | Driver |
|--------|----------|--------|
| `postgresql` | postgresql | postgres-jdbc |
| `sqlserver` | sqlserver | sqlserver |
| `mysql` | mysql | mysql8 |
| `oracle` | oracle | oracle_thin |
| `mongodb` | mongodb | mongodb |
| `dynamodb` | dynamodb | dynamodb |

---

## Config file reference

```yaml
project: MyApp                   # Project identifier

environments:                    # List of environments
  - DEV
  - QA
  - STAGE
  - PROD

connections:
  - name: "{{.Project}} {{.Env}} - PostgreSQL"  # Go template
    type: postgresql                              # Engine (see supported engines)
    host: "{{.Env | lower}}-db.internal"
    port: 5432
    database: "{{.Project | lower}}"
    folder: "{{.Env}}"           # DBeaver folder/group name

engines:                         # Engines for scripts/snippets generation
  - postgresql
  - sqlserver
```

**Template variables:**

| Variable | Type | Description |
|----------|------|-------------|
| `{{.Project}}` | string | Project name |
| `{{.Env}}` | string | Current environment (DEV, QA, etc.) |
| `{{.DBEngine}}` | string | Engine type |
| `{{.Env \| lower}}` | string | Lowercase environment |

---

## Roadmap

- **v1.x** — Exit code 2 for internal errors; `--merge` for connection updates
- **v2.0** — `dwm install` direct DBeaver copy ✅ (current)
- **v3.0** — DBeaver Plugin (Eclipse/Java) with UI buttons
- **v4.0** — Knowledge Graph + AI copilot for SQL schema understanding
