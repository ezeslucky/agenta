# CLI and Agent Daemon Guide

The `agenta` CLI connects your local machine to Agenta. It handles authentication, workspace management, issue tracking, and runs the agent daemon that executes AI tasks locally.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install agenta-ai/tap/agenta
```

### Build from Source

```bash
git clone https://github.com/ezeslucky/agenta.git
cd agenta
make build
cp server/bin/agenta /usr/local/bin/agenta
```

### Update

```bash
brew upgrade agenta-ai/tap/agenta
```

For install script or manual installs, use:

```bash
agenta update
```

`agenta update` auto-detects your installation method and upgrades accordingly.

## Quick Start

```bash
# One-command setup: configure, authenticate, and start the daemon
agenta setup

# For self-hosted (local) deployments:
agenta setup self-host
```

Or step by step:

```bash
# 1. Authenticate (opens browser for login)
agenta login

# 2. Start the agent daemon
agenta daemon start

# 3. Done — agents in your watched workspaces can now execute tasks on your machine
```

`agenta login` automatically discovers all workspaces you belong to and adds them to the daemon watch list.

## Authentication

### Browser Login

```bash
agenta login
```

Opens your browser for OAuth authentication, creates a 90-day personal access token, and auto-configures your workspaces.

### Token Login

```bash
agenta login --token <mul_...>
```

Authenticate using a personal access token directly. Useful for headless environments. Pass `--token=` with an empty value to be prompted interactively (so the token never lands in shell history).

### Check Status

```bash
agenta auth status
```

Shows your current server, user, and token validity.

### Logout

```bash
agenta auth logout
```

Removes the stored authentication token.

## Agent Daemon

The daemon is the local agent runtime. It detects available AI CLIs on your machine, registers them with the Agenta server, and executes tasks when agents are assigned work.

### Start

```bash
agenta daemon start
```

By default, the daemon runs in the background and logs to `~/.agenta/daemon.log`.

To run in the foreground (useful for debugging):

```bash
agenta daemon start --foreground
```

### Stop

```bash
agenta daemon stop
```

### Status

```bash
agenta daemon status
agenta daemon status --output json
```

Shows PID, uptime, detected agents, and watched workspaces.

### Logs

```bash
agenta daemon logs              # Last 50 lines
agenta daemon logs -f           # Follow (tail -f)
agenta daemon logs -n 100       # Last 100 lines
```

### Supported Agents

The daemon auto-detects these AI CLIs on your PATH:

| CLI | Command | Description |
|-----|---------|-------------|
| [Claude Code](https://docs.anthropic.com/en/docs/claude-code) | `claude` | Anthropic's coding agent |
| [Codex](https://github.com/openai/codex) | `codex` | OpenAI's coding agent |
| [GitHub Copilot CLI](https://docs.github.com/en/copilot) | `copilot` | GitHub's coding agent (model routed by your GitHub entitlement) |
| OpenCode | `opencode` | Open-source coding agent |
| OpenClaw | `openclaw` | Open-source coding agent |
| Hermes | `hermes` | Nous Research coding agent |
| Gemini | `gemini` | Google's coding agent |
| [Pi](https://pi.dev/) | `pi` | Pi coding agent |
| [Cursor Agent](https://cursor.com/) | `cursor-agent` | Cursor's headless coding agent |
| Kimi | `kimi` | Moonshot coding agent |
| Kiro CLI | `kiro-cli` | Kiro ACP coding agent |

You need at least one installed. The daemon registers each detected CLI as an available runtime.

### How It Works

1. On start, the daemon detects installed agent CLIs and registers a runtime for each agent in each watched workspace
2. It polls the server at a configurable interval (default: 3s) for claimed tasks
3. When a task arrives, it creates an isolated workspace directory, spawns the agent CLI, and streams results back
4. Heartbeats are sent periodically (default: 15s) so the server knows the daemon is alive
5. On shutdown, all runtimes are deregistered

### Configuration

Daemon behavior is configured via flags or environment variables:

| Setting | Flag | Env Variable | Default |
|---------|------|--------------|---------|
| Poll interval | `--poll-interval` | `AGENTA_DAEMON_POLL_INTERVAL` | `3s` |
| Heartbeat interval | `--heartbeat-interval` | `AGENTA_DAEMON_HEARTBEAT_INTERVAL` | `15s` |
| Agent timeout | `--agent-timeout` | `AGENTA_AGENT_TIMEOUT` | `0` (no cap; bounded by the watchdogs) |
| Codex semantic inactivity timeout | `--codex-semantic-inactivity-timeout` | `AGENTA_CODEX_SEMANTIC_INACTIVITY_TIMEOUT` | `10m` |
| Max concurrent tasks | `--max-concurrent-tasks` | `AGENTA_DAEMON_MAX_CONCURRENT_TASKS` | `20` |
| Daemon ID | `--daemon-id` | `AGENTA_DAEMON_ID` | hostname |
| Device name | `--device-name` | `AGENTA_DAEMON_DEVICE_NAME` | hostname |
| Runtime name | `--runtime-name` | `AGENTA_AGENT_RUNTIME_NAME` | `Local Agent` |
| Workspaces root | — | `AGENTA_WORKSPACES_ROOT` | `~/agenta_workspaces` |
| GC enabled | — | `AGENTA_GC_ENABLED` | `true` (set `false`/`0` to disable) |
| GC scan interval | — | `AGENTA_GC_INTERVAL` | `1h` |
| GC TTL (done/cancelled issues) | — | `AGENTA_GC_TTL` | `24h` |
| GC orphan TTL (no `.gc_meta.json`) | — | `AGENTA_GC_ORPHAN_TTL` | `72h` |
| GC artifact TTL (open issues) | — | `AGENTA_GC_ARTIFACT_TTL` | `12h` (set `0` to disable) |
| GC artifact patterns | — | `AGENTA_GC_ARTIFACT_PATTERNS` | `node_modules,.next,.turbo` |

#### Workspace garbage collection

The daemon periodically scans `AGENTA_WORKSPACES_ROOT` and reclaims disk space in three modes:

- **Full task cleanup** — when an issue's status is `done` or `cancelled` and has been idle for `AGENTA_GC_TTL`, the entire task directory is removed.
- **Orphan cleanup** — task directories with no `.gc_meta.json` (e.g. left over from a daemon crash) are removed once they exceed `AGENTA_GC_ORPHAN_TTL`.
- **Artifact-only cleanup** — when a task has been completed for at least `AGENTA_GC_ARTIFACT_TTL` but the issue is still open, regenerable build outputs whose directory basename matches `AGENTA_GC_ARTIFACT_PATTERNS` are removed; the rest of the workdir (source, `.git`, `output/`, `logs/`, `.gc_meta.json`) is preserved so the agent can resume the same workdir on the next task.

Patterns are basename-only — entries containing `/` or `\` are silently dropped — and `.git` subtrees are never descended into. The default list (`node_modules`, `.next`, `.turbo`) is intentionally narrow; extend it per deployment if your repos consistently produce other regenerable directories (for example, `AGENTA_GC_ARTIFACT_PATTERNS=node_modules,.next,.turbo,target,__pycache__`). To disable artifact cleanup entirely, set `AGENTA_GC_ARTIFACT_TTL=0`.

Agent-specific overrides:

| Variable | Description |
|----------|-------------|
| `AGENTA_CLAUDE_PATH` | Custom path to the `claude` binary |
| `AGENTA_CLAUDE_MODEL` | Override the Claude model used |
| `AGENTA_CLAUDE_ARGS` | Default extra arguments for Claude Code runs |
| `AGENTA_CODEX_PATH` | Custom path to the `codex` binary |
| `AGENTA_CODEX_MODEL` | Override the Codex model used |
| `AGENTA_CODEX_ARGS` | Default extra arguments for Codex runs |
| `AGENTA_COPILOT_PATH` | Custom path to the `copilot` binary |
| `AGENTA_COPILOT_MODEL` | Override the Copilot model used (note: GitHub Copilot routes models through your account entitlement, so this may not be honoured) |
| `AGENTA_OPENCODE_PATH` | Custom path to the `opencode` binary |
| `AGENTA_OPENCODE_MODEL` | Override the OpenCode model used |
| `AGENTA_OPENCLAW_PATH` | Custom path to the `openclaw` binary |
| `AGENTA_OPENCLAW_MODEL` | Override the OpenClaw model used |
| `AGENTA_HERMES_PATH` | Custom path to the `hermes` binary |
| `AGENTA_HERMES_MODEL` | Override the Hermes model used |
| `AGENTA_GEMINI_PATH` | Custom path to the `gemini` binary |
| `AGENTA_GEMINI_MODEL` | Override the Gemini model used |
| `AGENTA_PI_PATH` | Custom path to the `pi` binary |
| `AGENTA_PI_MODEL` | Override the Pi model used |
| `AGENTA_CURSOR_PATH` | Custom path to the `cursor-agent` binary |
| `AGENTA_CURSOR_MODEL` | Override the Cursor Agent model used |
| `AGENTA_KIMI_PATH` | Custom path to the `kimi` binary |
| `AGENTA_KIMI_MODEL` | Override the Kimi model used |
| `AGENTA_KIRO_PATH` | Custom path to the `kiro-cli` binary |
| `AGENTA_KIRO_MODEL` | Override the Kiro model used |

`AGENTA_CLAUDE_ARGS` and `AGENTA_CODEX_ARGS` are parsed with POSIX shellword quoting, so values such as `--model "gpt-5.1 codex" --sandbox read-only` are split like a shell command line. Agent arguments are applied in this order: hardcoded Agenta defaults, daemon-wide env defaults, then per-agent `custom_args` from the task.

### Self-Hosted Server

When connecting to a self-hosted Agenta instance, the easiest approach is:

```bash
# One command — configures for localhost, authenticates, starts daemon
agenta setup self-host

# Or for on-premise with custom domains:
agenta setup self-host --server-url https://api.example.com --app-url https://app.example.com
```

Or configure manually:

```bash
# Set URLs individually
agenta config set server_url http://localhost:8080
agenta config set app_url http://localhost:3000

# For production with TLS:
# agenta config set server_url https://api.example.com
# agenta config set app_url https://app.example.com

agenta login
agenta daemon start
```

### Profiles

Profiles let you run multiple daemons on the same machine — for example, one for production and one for a staging server.

```bash
# Set up a staging profile
agenta setup self-host --profile staging --server-url https://api-staging.example.com --app-url https://staging.example.com

# Start its daemon
agenta daemon start --profile staging

# Default profile runs separately
agenta daemon start
```

Each profile gets its own config directory (`~/.agenta/profiles/<name>/`), daemon state, health port, and workspace root.

## Workspaces

### Working with multiple workspaces

Every command runs against a single workspace. The CLI resolves which one in this order (highest priority first):

1. `--workspace-id <id>` flag on the command
2. `AGENTA_WORKSPACE_ID` environment variable
3. The default workspace stored in your current profile (set by `agenta workspace switch` or `agenta login`)

`agenta workspace switch <id|slug>` is the day-to-day way to change the default workspace. For scripting and headless setups where you don't want any stored state, prefer the `--workspace-id` flag or the env variable. `agenta config set workspace_id <id>` is the low-level equivalent of `switch` (it writes the same setting but skips the access check).

If you need full isolation between organizations or accounts — separate tokens, separate daemons, separate config dirs — use `--profile <name>` instead. Each profile keeps its own default workspace.

### List Workspaces

```bash
agenta workspace list
agenta workspace list --full-id
agenta workspace list --output json
```

The current default workspace is marked with `*`. Table output shows short UUID prefixes — pass `--full-id` when you need the canonical UUIDs.

### Switch Default Workspace

```bash
agenta workspace switch <workspace-id>
agenta workspace switch <slug>
```

Verifies you have access to the workspace, then sets it as the default for the current profile. Subsequent commands without `--workspace-id` and `AGENTA_WORKSPACE_ID` target this workspace. Pair `--profile` if you want to change a non-default profile's workspace.

### Get Details

```bash
agenta workspace get <workspace-id>
agenta workspace get <workspace-id> --output json
```

Passing no `<workspace-id>` resolves to the current default workspace, so `agenta workspace get` doubles as "what workspace am I on?".

### List Members

```bash
agenta workspace member list <workspace-id>
```

## Issues

### List Issues

```bash
agenta issue list
agenta issue list --status in_progress
agenta issue list --priority urgent --assignee "Agent Name"
agenta issue list --assignee-id 5fb87ac7-23b5-4a7a-81fa-ed295a54545d
agenta issue list --full-id
agenta issue list --limit 20 --output json
```

Table output shows a routable issue `KEY` such as `MUL-123`; copy that key into follow-up commands like `issue get`, `issue comment list`, `issue status`, or `--parent`. Add `--full-id` when you need canonical UUIDs. Available filters: `--status`, `--priority`, `--assignee` / `--assignee-id`, `--project`, `--metadata`, `--limit`. Use `--assignee-id <uuid>` for unambiguous filtering when names overlap.

Use `--metadata key=value` (repeatable; combined with AND) to filter by per-issue metadata. The value is JSON-parsed: `true`/`false` become bool, numbers become numbers, anything else is a string. Wrap as `'"42"'` to force a string when the value would otherwise sniff as a number:

```bash
agenta issue list --metadata pipeline_status=waiting_review
agenta issue list --metadata pr_number=482 --metadata is_blocked=true
```

### Get Issue

```bash
agenta issue get <id>
agenta issue get <id> --output json
```

### Create Issue

```bash
agenta issue create --title "Fix login bug" --description "..." --priority high --assignee "Lambda"
agenta issue create --title "Fix login bug" --assignee-id 5fb87ac7-23b5-4a7a-81fa-ed295a54545d
```

Flags: `--title` (required), `--description`, `--status`, `--priority`, `--assignee` / `--assignee-id`, `--parent`, `--project`, `--due-date`. Pass `--assignee-id <uuid>` (mutually exclusive with `--assignee`) when scripting against the IDs returned by `agenta workspace member list --output json` / `agenta agent list --output json`.

### Update Issue

```bash
agenta issue update <id> --title "New title" --priority urgent
```

### Assign Issue

```bash
agenta issue assign <id> --to "Lambda"
agenta issue assign <id> --to-id 5fb87ac7-23b5-4a7a-81fa-ed295a54545d
agenta issue assign <id> --unassign
```

Pass `--to-id <uuid>` to assign by canonical UUID (mutually exclusive with `--to`); useful when names overlap across members and agents.

### Change Status

```bash
agenta issue status <id> in_progress
```

Valid statuses: `backlog`, `todo`, `in_progress`, `in_review`, `done`, `blocked`, `cancelled`.

### Comments

```bash
# List comments — flat timeline, chronological. Hard cap of 2000 rows; on
# long-running issues prefer one of the thread-aware reads below to keep
# context windows tight.
agenta issue comment list <issue-id>

# Single thread (root + every descendant). Anchor may be the root itself
# or any reply inside the thread — the server walks up to the root.
agenta issue comment list <issue-id> --thread <comment-id>

# Single thread, capped to the N most recent replies. The thread root is
# always included (even with --tail 0), so an agent landing on a long
# thread keeps the "what is this about" context without dragging hundreds
# of replies into its prompt.
agenta issue comment list <issue-id> --thread <comment-id> --tail 30

# Scroll older replies inside the same thread. --before / --before-id are
# the reply cursor that the previous response emitted on stderr as
# `Next reply cursor: --before <ts> --before-id <reply-id>`.
agenta issue comment list <issue-id> --thread <comment-id> --tail 30 \
    --before <ts> --before-id <reply-id>

# Most recently active threads (root + every descendant), grouped by
# thread. Returns N complete conversational arcs, oldest-active first so
# the freshest thread sits closest to "now" in an agent prompt.
agenta issue comment list <issue-id> --recent 20

# Scroll older threads. Under --recent, --before / --before-id are a
# THREAD cursor (thread last_activity_at + root id), emitted on stderr as
# `Next thread cursor: --before <ts> --before-id <root-id>`.
agenta issue comment list <issue-id> --recent 20 \
    --before <ts> --before-id <root-id>

# Incremental polling. Combines with --thread or --recent; filters out
# replies created on or before <ts> from the page (the thread root is
# exempt so the agent always gets context).
agenta issue comment list <issue-id> --thread <comment-id> --tail 30 \
    --since <RFC3339-timestamp>

# Add a comment
agenta issue comment add <issue-id> --content "Looks good, merging now"

# Reply to a specific comment
agenta issue comment add <issue-id> --parent <comment-id> --content "Thanks!"

# Delete a comment
agenta issue comment delete <comment-id>
```

**`--before` / `--before-id` semantics depend on the paging mode**, by
design — same flag, different scope:

| Mode | What the cursor walks | stderr label |
| --- | --- | --- |
| `--recent N` | Older *threads* (last_activity_at, root_id) | `Next thread cursor` |
| `--thread <id> --tail N` | Older *replies* inside that thread (created_at, id) | `Next reply cursor` |

Outside those two modes (`--thread` without `--tail`, or no `--thread`
and no `--recent`) the cursor flags are rejected so they cannot silently
no-op. The server emits the cursor headers (`X-Agenta-Next-Before` /
`X-Agenta-Next-Before-Id`) only when an older page actually exists —
exact-boundary pages (e.g. `--tail 3` on a thread with exactly 3
replies) intentionally return no cursor so callers stop paginating.

When `--since` is combined with `--recent` or `--thread --tail`, the
server additionally suppresses the cursor once the cursor target itself
is older than `since`. Older pages walk strictly older rows, so they
cannot satisfy `> since` either — emitting a cursor there would just
hand back root-only pages until the caller reaches the start of the
thread / issue. Incremental polling stops at the first page whose
cursor target falls before the watermark.

### Metadata

Per-issue metadata is a small KV map agents use to track pipeline state (PR number, pipeline status, waiting_on, ...). Keys match `^[a-zA-Z_][a-zA-Z0-9_.-]{0,63}$`, values are primitives (string / number / bool), max 50 keys per issue, blob capped at 8KB.

The bar for writing is high: pin a value only when it is materially important to the issue AND likely to be re-read by future runs on this same issue (the PR URL, the deploy URL, what we're blocked on). Most runs write zero new keys — that's the expected case. Don't pin runtime bookkeeping like `attempts`, single-run investigation notes, large logs, secrets/tokens, or description/comment copies — see the agent runtime prompt for the full anti-pattern list.

```bash
# List every key on an issue
agenta issue metadata list <issue-id>

# Read a single key
agenta issue metadata get <issue-id> --key pipeline_status

# Write a single key — value auto-typed (true/false → bool, numbers → number, else string)
agenta issue metadata set <issue-id> --key pipeline_status --value waiting_review
agenta issue metadata set <issue-id> --key pr_number --value 482
agenta issue metadata set <issue-id> --key is_blocked --value true

# Force a specific type when sniffing would pick the wrong one
agenta issue metadata set <issue-id> --key code --value 42 --type string

# Remove a key
agenta issue metadata delete <issue-id> --key pipeline_status
```

All writes are single-key atomic — concurrent agents writing different keys do not lose each other's updates. To query, use `agenta issue list --metadata key=value` (see *List Issues* above).

### Subscribers

```bash
# List subscribers of an issue
agenta issue subscriber list <issue-id>

# Subscribe yourself to an issue
agenta issue subscriber add <issue-id>

# Subscribe another member or agent by name
agenta issue subscriber add <issue-id> --user "Lambda"

# Unsubscribe yourself
agenta issue subscriber remove <issue-id>

# Unsubscribe another member or agent
agenta issue subscriber remove <issue-id> --user "Lambda"
```

Subscribers receive notifications about issue activity (new comments, status changes, etc.). Without `--user`, the command acts on the caller.

### Execution History

```bash
# List all execution runs for an issue
agenta issue runs <issue-id>
agenta issue runs <issue-id> --full-id
agenta issue runs <issue-id> --output json

# View messages for a specific execution run
agenta issue run-messages <task-id>
agenta issue run-messages <short-task-id> --issue <issue-id>
agenta issue run-messages <task-id> --output json

# Incremental fetch (only messages after a given sequence number)
agenta issue run-messages <task-id> --since 42 --output json
```

The `runs` command shows all past and current executions for an issue, including running tasks. Table output uses short task UUID prefixes by default; pass `--full-id` to print canonical task UUIDs. The `run-messages` command accepts full task UUIDs directly; copied short task prefixes must be scoped with `--issue <issue-id>` so the CLI only checks that issue's runs. It shows the detailed message log (tool calls, thinking, text, errors) for a single run. Use `--since` for efficient polling of in-progress runs.

## Projects

Projects group related issues (e.g. a sprint, an epic, a workstream). Every project
belongs to a workspace and can optionally have a lead (member or agent).

### List Projects

```bash
agenta project list
agenta project list --status in_progress
agenta project list --output json
```

Available filters: `--status`.

### Get Project

```bash
agenta project get <id>
agenta project get <id> --output json
```

### Create Project

```bash
agenta project create --title "2026 Week 16 Sprint" --icon "🏃" --lead "Lambda"
```

Flags: `--title` (required), `--description`, `--status`, `--icon`, `--lead`.

### Update Project

```bash
agenta project update <id> --title "New title" --status in_progress
agenta project update <id> --lead "Lambda"
```

Flags: `--title`, `--description`, `--status`, `--icon`, `--lead`.

### Change Status

```bash
agenta project status <id> in_progress
```

Valid statuses: `planned`, `in_progress`, `paused`, `completed`, `cancelled`.

### Delete Project

```bash
agenta project delete <id>
```

### Associating Issues with Projects

Use the `--project` flag on `issue create` / `issue update` to attach an issue to a
project, or on `issue list` to filter issues by project:

```bash
agenta issue create --title "Login bug" --project <project-id>
agenta issue update <issue-id> --project <project-id>
agenta issue list --project <project-id>
```

## Setup

```bash
# One-command setup for Agenta Cloud: configure, authenticate, and start the daemon
agenta setup

# For local self-hosted deployments
agenta setup self-host

# Custom ports
agenta setup self-host --port 9090 --frontend-port 4000

# On-premise with custom domains
agenta setup self-host --server-url https://api.example.com --app-url https://app.example.com
```

`agenta setup` configures the CLI, opens your browser for authentication, and starts the daemon — all in one step. Use `agenta setup self-host` to connect to a self-hosted server instead of Agenta Cloud.

## Configuration

### View Config

```bash
agenta config show
```

Shows config file path, server URL, app URL, and default workspace.

### Set Values

```bash
agenta config set server_url https://api.example.com
agenta config set app_url https://app.example.com
agenta config set workspace_id <workspace-id>
```

`config set workspace_id <id>` is the low-level interface — it writes the value verbatim without checking that the workspace exists or that you have access. Prefer `agenta workspace switch <id|slug>` for day-to-day workspace changes; it does both checks before saving.

## Autopilot Commands

Autopilots are scheduled/triggered automations that dispatch agent tasks (either by creating an issue or by running an agent directly).

### List Autopilots

```bash
agenta autopilot list
agenta autopilot list --full-id
agenta autopilot list --status active --output json
```

Autopilot table IDs are short UUID prefixes; follow-up autopilot commands accept copied prefixes when they are unique in the current workspace. Use `--full-id` to print canonical UUIDs.

### Get Autopilot Details

```bash
agenta autopilot get <id>
agenta autopilot get <id> --output json   # includes triggers
```

### Create / Update / Delete

```bash
agenta autopilot create \
  --title "Nightly bug triage" \
  --description "Scan todo issues and prioritize." \
  --agent "Lambda" \
  --mode create_issue

agenta autopilot update <id> --status paused
agenta autopilot update <id> --description "New prompt"
agenta autopilot delete <id>
```

`--mode` accepts `create_issue` (creates a new issue on each run and assigns it to the agent) or `run_only` (enqueues a direct agent task without creating an issue). `--agent` accepts either a name or UUID.

### Manual Trigger

```bash
agenta autopilot trigger <id>            # Fires the autopilot once, returns the run
```

### Run History

```bash
agenta autopilot runs <id>
agenta autopilot runs <id> --limit 50 --output json
```

### Schedule Triggers

```bash
agenta autopilot trigger-add <autopilot-id> --cron "0 9 * * 1-5" --timezone "America/New_York"
agenta autopilot trigger-update <autopilot-id> <trigger-id> --enabled=false
agenta autopilot trigger-delete <autopilot-id> <trigger-id>
```

Only cron-based `schedule` triggers are currently exposed via the CLI. The data model also defines `webhook` and `api` kinds, but there is no server endpoint that fires them yet, so they're not surfaced here.

## Other Commands

```bash
agenta version              # Show CLI version and commit hash
agenta update               # Update to latest version
agenta agent list           # List agents in the current workspace
```

## Output Formats

Most commands support `--output` with two formats:

- `table` — human-readable table (default for list commands)
- `json` — structured JSON (useful for scripting and automation)

```bash
agenta issue list --output json
agenta daemon status --output json
```

## Error Messages

The CLI funnels command errors returned to the top-level handler through a
single user-facing translation layer (`server/internal/cli/errors.go`) so that
what you see on the terminal is a short, actionable sentence rather than a raw
Go error, an HTTP status line, or an internal `resolve issue: ...` chain. (A
few commands print their own output or run deliberate fast probes — for example
`setup`'s short `/health` reachability check — and don't go through this
layer.) The underlying detail is still available on demand (see `--debug`).

### What you see

- **Friendly, single-line message.** Transport failures (timeout, DNS,
  connection refused, TLS) and HTTP status failures (401/403/404/409/400·422/
  429/5xx) are each rendered as one clear sentence with a next step — for
  example a timeout suggests checking the network or raising
  `AGENTA_HTTP_TIMEOUT`, and a 401 tells you to run `agenta login`.
- **Server-provided validation messages are preserved.** For a 400/422 that
  carries a message from the server, that message is shown verbatim
  (`Invalid request: <server message>`); only when there is none do you get the
  generic "check your values / run with --help" hint.
- **No leaked internals by default.** Raw URLs, status lines, JSON bodies, and
  the internal verb chain are hidden unless you ask for them.

### Language

Messages default to **English**, matching the rest of the CLI's help output.
If a Chinese locale is detected in `LC_ALL`, `LC_MESSAGES`, or `LANG` (in that
precedence order), messages switch to **Chinese**. No flag is needed; set the
locale as usual:

```bash
LANG=zh_CN.UTF-8 agenta issue get MUL-9999   # 错误信息显示为中文
```

### Exit codes

The process exit code is tiered so scripts can branch on the failure class:

| Exit code | Meaning |
| --- | --- |
| `0` | success |
| `1` | generic / unclassified error |
| `2` | network error (timeout, DNS, connection refused, TLS, offline) |
| `3` | authentication / authorization (HTTP 401, 403) |
| `4` | not found (HTTP 404) |
| `5` | validation (HTTP 400, 422) |

```bash
agenta issue get MUL-9999
if [ $? -eq 4 ]; then echo "no such issue"; fi
```

### Seeing the full detail (`--debug`)

Pass the global `--debug` flag (or set `AGENTA_DEBUG=1`) to print the complete
original error chain — the internal verb chain, the request method/path/status,
and the raw server body — underneath the friendly message. Use it when you need
to file a bug or understand exactly what the server returned:

```bash
agenta issue list --debug
AGENTA_DEBUG=1 agenta issue update MUL-1234 --title "x"
```

### Request timeout

API requests use a default timeout of 30 seconds. Override it with
`AGENTA_HTTP_TIMEOUT` when you are on a slow network; it accepts a Go duration
(`45s`, `2m`) or a plain number of seconds (`45`). Command-level deadlines are
always at least this value, so raising it takes effect across all commands.

```bash
AGENTA_HTTP_TIMEOUT=60s agenta issue list
```
