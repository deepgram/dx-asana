# Changelog

All notable changes to this project will be documented in this file.


## [Unreleased]

### Fixed
- Schema mismatches with Asana's actual API. `User.photo` is now an
  object (multi-size avatar URLs) instead of a string, fixing the
  `me` command. `Task.notes`/`html_notes` replace the fictional
  `description` field. Removed `Task.priority`/`Task.status` and
  `User.timezone` (none of these exist on Asana resources).
- `Project.current_status_update` replaces the broken `status`
  string. `Section.project` is parsed as a `ProjectCompact` object
  rather than a GID string.
- Nullable timestamps (`completed_at`, `created_at`, etc.) are
  pointers so `null` no longer breaks unmarshalling.
- `--version` now reports the actual `git describe` output.
  Previously it always printed `dev` because LDFLAGS targeted
  `main.*` while cobra read `cmd.*`.
- `--json` output is no longer contaminated by Cobra usage text on
  errors. Errors now route to stderr; stdout stays clean for
  scripting.
- `config get` shows `(not set)` for an empty token instead of
  reusing the same `***` mask used for short tokens.
- `config set --token ""` was a silent no-op; new `config unset`
  subcommand handles clearing.
- Client-side validation rejects malformed `--due` values before
  hitting the API.

### Added
- `opt_fields` plumbing across all client methods plus default
  field sets. Previously list/search returned only the compact
  representation (gid, name); now they return notes, assignee, due,
  memberships, permalink, etc. by default. Override per call with
  `--fields name,due_on,...`.
- Structured API errors. The `{"errors":[{"message","help"}]}`
  envelope is parsed into `*asana.APIError` and surfaced in JSON
  output as `errors[]` and `status` alongside the legacy `error`
  string.
- New commands:
  - `my-tasks` — tasks across all projects assigned to the user
  - `comment` — post or list comments (Asana stories) on a task
  - `workspaces` — list available workspaces
  - `projects <ws>` — list projects in a workspace
  - `sections <project>` — list sections in a project
  - `subtasks <task>` — list subtasks of a parent task
  - `config unset` — clear stored token/workspace/projects
- `--parent <task-gid>` flag on `create` to make subtasks. Project
  arg becomes optional in that mode.
- `--fields` / `--limit` / `--assignee` / `--project` /
  `--completed` flags on `search`.
- `User-Agent: asana-cli/<version>` header on every API request.
- `make smoke` target and `scripts/smoke.sh` for read-only
  end-to-end verification.

### Changed
- `client.GetTasks(projectGID, filters)` is now
  `client.GetTasks(projectGID, opts ...Option)`. Old `filters`
  callers must move to `WithAssignee`/`WithCompletedSince`/etc.
- `client.Search` retained as a thin wrapper around the new
  `client.SearchTasks(workspaceGID, opts...)` signature.
- `Task.DueDate` renamed to `Task.DueOn` to match Asana naming.
- `Task.Description` renamed to `Task.Notes`.
- Sync daemon now caches full task data via `DefaultTaskFields`
  rather than the compact response.

### Earlier work
- TUI with Bubble Tea
- JSON output format
- Background sync daemon
- Multiple project support
- Complete CRUD operations
- Time parsing for Asana date formats
- Updated to use Asana API GID terminology

## [0.1.0] - 2026-02-17

### Added
- Initial release
- Task listing and management
- Project switching
- Configuration management
- Homebrew support