#!/usr/bin/env bash
set -euo pipefail

# Read-only end-to-end smoke check for asana-cli. Does not mutate any Asana
# data. Requires ASANA_TOKEN to be set.
#
# Usage: scripts/smoke.sh [path-to-binary]
#   defaults to ./asana-cli

BIN="${1:-./asana-cli}"

if [[ -z "${ASANA_TOKEN:-}" ]]; then
  echo "ASANA_TOKEN not set"
  exit 1
fi

if [[ ! -x "$BIN" ]]; then
  echo "binary not executable: $BIN"
  exit 1
fi

pass() { printf "  \033[32m✓\033[0m %s\n" "$1"; }
fail() { printf "  \033[31m✗\033[0m %s — %s\n" "$1" "$2"; exit 1; }

assert_json_success() {
  local label="$1" output="$2"
  if ! printf "%s" "$output" | jq -e '.success == true' >/dev/null 2>&1; then
    fail "$label" "JSON envelope did not report success: $output"
  fi
  pass "$label"
}

echo "asana-cli smoke (read-only)"

OUT=$("$BIN" --version)
[[ "$OUT" =~ ^asana-cli\ version ]] && pass "--version" || fail "--version" "$OUT"

OUT=$("$BIN" me --json 2>/dev/null)
assert_json_success "me --json" "$OUT"

OUT=$("$BIN" workspaces --json 2>/dev/null)
assert_json_success "workspaces --json" "$OUT"

OUT=$("$BIN" my-tasks --json 2>/dev/null)
assert_json_success "my-tasks --json" "$OUT"

WS=$(printf "%s" "$OUT" | jq -r '.meta.user_gid // empty')
[[ -n "$WS" ]] && pass "my-tasks meta.user_gid populated" || fail "my-tasks" "no user_gid in meta"

OUT=$("$BIN" me --json 2>/dev/null | jq -r '.data.workspaces[0].gid // empty')
[[ -n "$OUT" ]] || fail "me workspace lookup" "no workspaces[]"
WORKSPACE_GID="$OUT"
pass "workspace GID resolved: $WORKSPACE_GID"

OUT=$("$BIN" projects "$WORKSPACE_GID" --json 2>/dev/null)
assert_json_success "projects --json" "$OUT"

PROJ=$(printf "%s" "$OUT" | jq -r '.data[0].gid // empty')
[[ -n "$PROJ" ]] || fail "project sample" "no projects[]"
pass "sample project GID: $PROJ"

OUT=$("$BIN" sections "$PROJ" --json 2>/dev/null)
assert_json_success "sections --json" "$OUT"

OUT=$("$BIN" list "$PROJ" --json 2>/dev/null)
assert_json_success "list --json" "$OUT"

TASK=$(printf "%s" "$OUT" | jq -r '.data[0].gid // empty')
[[ -n "$TASK" ]] && pass "sample task GID: $TASK" || pass "list returned 0 tasks (skipping view)"

if [[ -n "$TASK" ]]; then
  OUT=$("$BIN" view "$TASK" --json 2>/dev/null)
  assert_json_success "view --json" "$OUT"

  OUT=$("$BIN" comment "$TASK" --list --json 2>/dev/null)
  assert_json_success "comment --list --json" "$OUT"

  OUT=$("$BIN" subtasks "$TASK" --json 2>/dev/null)
  assert_json_success "subtasks --json" "$OUT"
fi

echo ""
echo "all smoke checks passed"
