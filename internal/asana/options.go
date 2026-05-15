package asana

import (
	"net/url"
	"strconv"
	"strings"
)

// Option configures a query string for an Asana request. Use the With* helpers
// (WithOptFields, WithLimit, WithCompletedSince, etc.) at call sites instead
// of constructing url.Values directly.
type Option func(*url.Values)

// WithOptFields requests specific response fields from Asana via the
// `opt_fields` query parameter. Without this, list and search endpoints
// return only the compact representation (gid, resource_type, name).
//
// Field paths follow Asana's dot notation, e.g. "assignee.name" or
// "memberships.section.name". Reference:
// https://developers.asana.com/docs/inputoutput-options
func WithOptFields(fields ...string) Option {
	return func(v *url.Values) {
		if len(fields) == 0 {
			return
		}
		v.Set("opt_fields", strings.Join(fields, ","))
	}
}

// WithLimit caps the number of results returned by a list/search endpoint.
// Asana's API maximum is 100; pass higher values at your own risk.
func WithLimit(n int) Option {
	return func(v *url.Values) {
		if n > 0 {
			v.Set("limit", strconv.Itoa(n))
		}
	}
}

// WithOffset paginates a list/search endpoint using the cursor returned in
// `next_page.offset` from the previous page.
func WithOffset(cursor string) Option {
	return func(v *url.Values) {
		if cursor != "" {
			v.Set("offset", cursor)
		}
	}
}

// WithAssignee filters task list/search to a specific assignee GID.
func WithAssignee(userGID string) Option {
	return func(v *url.Values) {
		if userGID != "" {
			v.Set("assignee", userGID)
		}
	}
}

// WithCompletedSince filters tasks completed after a given timestamp. Pass
// "now" to exclude completed tasks entirely.
func WithCompletedSince(since string) Option {
	return func(v *url.Values) {
		if since != "" {
			v.Set("completed_since", since)
		}
	}
}

// WithModifiedSince filters tasks modified after a given timestamp.
func WithModifiedSince(since string) Option {
	return func(v *url.Values) {
		if since != "" {
			v.Set("modified_since", since)
		}
	}
}

// WithWorkspace pins a request to a workspace GID. Required by some
// endpoints (e.g. task search) and accepted as a filter by others.
func WithWorkspace(workspaceGID string) Option {
	return func(v *url.Values) {
		if workspaceGID != "" {
			v.Set("workspace", workspaceGID)
		}
	}
}

// WithProject filters task list/search to a specific project GID.
func WithProject(projectGID string) Option {
	return func(v *url.Values) {
		if projectGID != "" {
			v.Set("projects", projectGID)
		}
	}
}

// WithSection filters task list to a specific section GID.
func WithSection(sectionGID string) Option {
	return func(v *url.Values) {
		if sectionGID != "" {
			v.Set("section", sectionGID)
		}
	}
}

// WithText is a free-text search query for /workspaces/{gid}/tasks/search.
func WithText(query string) Option {
	return func(v *url.Values) {
		if query != "" {
			v.Set("text", query)
		}
	}
}

// WithRaw lets callers add an arbitrary key/value query parameter for cases
// not covered by the named helpers (e.g. custom_fields.X.value).
func WithRaw(key, value string) Option {
	return func(v *url.Values) {
		if key != "" && value != "" {
			v.Set(key, value)
		}
	}
}

// applyOptions builds a query string from a set of options.
func applyOptions(opts []Option) string {
	if len(opts) == 0 {
		return ""
	}
	q := url.Values{}
	for _, opt := range opts {
		opt(&q)
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// DefaultTaskListFields is a sensible opt_fields default for task list
// commands so the CLI shows useful information without forcing every caller
// to specify field paths.
var DefaultTaskListFields = []string{
	"gid",
	"name",
	"resource_subtype",
	"completed",
	"completed_at",
	"due_on",
	"due_at",
	"start_on",
	"notes",
	"assignee.name",
	"assignee.gid",
	"permalink_url",
	"memberships.project.name",
	"memberships.section.name",
	"num_subtasks",
	"parent.name",
	"parent.gid",
	"tags.name",
}

// DefaultTaskFields is the field set for single-task fetches when more
// detail than the API default is desired.
var DefaultTaskFields = append(
	DefaultTaskListFields,
	"html_notes",
	"custom_fields.gid",
	"custom_fields.name",
	"custom_fields.display_value",
	"custom_fields.resource_subtype",
	"custom_fields.enum_value.name",
	"custom_fields.number_value",
	"custom_fields.text_value",
	"workspace.name",
	"followers.name",
	"created_at",
	"modified_at",
	"created_by.name",
)

// DefaultProjectListFields is a sensible default for project list commands.
var DefaultProjectListFields = []string{
	"gid",
	"name",
	"archived",
	"color",
	"icon",
	"due_on",
	"start_on",
	"owner.name",
	"team.name",
	"permalink_url",
}
