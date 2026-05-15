# dx-asana

Deepgram DX's command-line interface for Asana. Interactive TUI (Bubble Tea), JSON output for scripting, and a background sync daemon for fast local lookups.

This is a friendly fork of [`Pkill-MyDaemons/asana-cli`](https://github.com/Pkill-MyDaemons/asana-cli) (originally by [@TheCoolRobot](https://github.com/TheCoolRobot)). Upstream has been quiet since early 2026, so the Deepgram DX team has picked it up so we can ship schema fixes and the daily-driver commands we need. See [Acknowledgements](#acknowledgements).

## Features

- **Interactive TUI** — navigate, filter, and update tasks via Bubble Tea
- **JSON output** — every command supports `--json` with a stable `{success, data, meta, errors}` envelope
- **Sync daemon** — background process that caches project data locally on a 5-minute loop
- **Full task CRUD** — list, view, create, update, complete, delete, search
- **Daily-driver commands** — `my-tasks`, `comment`, `subtasks`, `workspaces`, `projects`, `sections`, plus `--parent` for subtasks
- **Discoverable** — `opt_fields` plumbed through every list/search call with sensible defaults

## Install

Homebrew distribution is paused while the fork settles. Install from source until that comes back:

```bash
git clone https://github.com/deepgram/dx-asana.git
cd dx-asana

# Build a local ./dx-asana binary
make build

# Or install into $GOPATH/bin (or $GOBIN)
make install
```

Requires Go 1.22 or newer. The build embeds the current commit and `git describe` output into `--version`.

## Authenticate

Asana issues personal access tokens at [app.asana.com/0/my-apps](https://app.asana.com/0/my-apps). Pass it via env var or save it to the config file:

```bash
export ASANA_TOKEN=your-token-here

# or persist it
dx-asana config set --token your-token-here
```

Config lives at `~/.dx-asana/config.json` with `0600` permissions. The file holds your token, your saved projects, and your default workspace.

## Basic usage

```bash
# Tasks assigned to you across every workspace
dx-asana my-tasks

# List tasks in a project (TUI by default)
dx-asana list <project-gid>
dx-asana list <project-gid> --json

# Save a project under a friendly name and use it as the default
dx-asana config project add backlog 1234567890 --description "DX backlog"
dx-asana config project switch backlog
dx-asana list             # uses the active project
dx-asana list --json

# Create
dx-asana create <project-gid> --name "Write README"
dx-asana create --parent <task-gid> --name "Subtask of the task above"

# Update + complete + delete
dx-asana update   <task-gid> --name "Renamed"
dx-asana complete <task-gid>
dx-asana delete   <task-gid>

# Search
dx-asana search <workspace-gid> "bug fix"

# Comments (Asana stories)
dx-asana comment <task-gid> --list
dx-asana comment <task-gid> --text "Shipping today"

# Discovery
dx-asana workspaces
dx-asana projects <workspace-gid>
dx-asana sections <project-gid>
dx-asana subtasks <task-gid>

# Sync daemon — caches projects locally every 5 minutes
dx-asana sync --projects 12345,67890
```

## JSON output

Every command accepts `--json`. The envelope is stable across commands and includes structured Asana errors:

```bash
$ dx-asana list <project-gid> --json
{
  "success": true,
  "data": [
    {
      "gid": "1234567890",
      "name": "Build feature",
      "completed": false,
      "due_on": "2026-03-01",
      "notes": "Long-form notes",
      "permalink_url": "https://app.asana.com/0/.../1234567890"
    }
  ],
  "meta": {
    "count": 1,
    "project_id": "<project-gid>"
  }
}

$ dx-asana view bogus --json
{
  "success": false,
  "status": 404,
  "errors": [{ "message": "Not found", "help": "..." }],
  "error": "Not found"
}
```

`--json` writes only JSON to stdout. Human-readable errors and Cobra usage text go to stderr, so `dx-asana list <p> --json | jq` is always safe.

## TUI controls

```
[↑↓]     Navigate tasks
[space]  Select / deselect
[c]      Mark complete
[f]      Toggle show completed
[s]      Change sort (name / due / priority)
[/]      Search
[a]      Add
[q]      Quit
```

## Commands

### Tasks
- `list` — list tasks in a project
- `view` — view a single task
- `create` — create a task (top-level or subtask via `--parent`)
- `update` — update task fields
- `complete` — mark complete
- `delete` — delete a task
- `search` — full-text typeahead across a workspace
- `my-tasks` — tasks assigned to the authenticated user across all projects
- `comment` — list or post comments on a task
- `subtasks` — list subtasks of a parent task

### Discovery
- `workspaces` — list workspaces you can access
- `projects <workspace-gid>` — list projects in a workspace
- `sections <project-gid>` — list sections in a project
- `me` — show the authenticated user

### System
- `config` — get / set / unset stored values; manage saved projects
- `sync` — start the background cache daemon

Run `dx-asana <command> --help` for flags and `--fields` overrides.

## Security

- Tokens are stored in `~/.dx-asana/config.json` with `0600` permissions
- Never commit `.env` files or raw tokens
- Use `ASANA_TOKEN` in CI and treat it as a secret

## Develop

```bash
make build          # local binary
make test           # unit + integration tests
make smoke          # read-only live smoke check (requires ASANA_TOKEN)
make lint           # golangci-lint
make fmt            # go fmt
make release        # cross-compile to dist/
```

## Contributing

Issues and PRs are welcome at [`deepgram/dx-asana`](https://github.com/deepgram/dx-asana). The standard flow:

1. Fork `deepgram/dx-asana`
2. Branch (`git checkout -b feature/thing`)
3. Commit with a clear subject line
4. Push and open a PR against `main`

Tests and lint must pass. For larger changes (new command surface, behavioural changes to the JSON envelope, etc.), open an issue first to align on scope.

## Acknowledgements

`dx-asana` is a fork of [`Pkill-MyDaemons/asana-cli`](https://github.com/Pkill-MyDaemons/asana-cli), itself a continuation of [`TheCoolRobot/asana-cli`](https://github.com/TheCoolRobot/asana-cli). Huge thanks to the original author and maintainers — without their work as a starting point this would have been a lot more typing. Upstream is dormant, so this fork has diverged with schema fixes, new commands, structured errors, and other changes; see [`CHANGELOG.md`](CHANGELOG.md) for the full delta.

## License

MIT — see [`LICENSE`](LICENSE).
