package asana

import (
	"time"
)

// CustomTime handles both date-only (YYYY-MM-DD) and full datetime (RFC3339)
// formats from Asana, plus null for missing values.
//
// Asana uses date-only for `*_on` fields (due_on, start_on, completed_at on
// projects), and RFC3339 for `*_at` fields (due_at, start_at, completed_at on
// tasks). A pointer to CustomTime is the canonical way to represent these.
type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == `""` {
		return nil
	}
	s = s[1 : len(s)-1]

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		ct.Time = t
		return nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		ct.Time = t
		return nil
	}
	_, err := time.Parse("2006-01-02", s)
	return err
}

// Photo is the multi-size avatar URLs returned by Asana on User responses.
// Asana returns this as `null` when the user has no avatar set.
type Photo struct {
	Image21x21     string `json:"image_21x21,omitempty"`
	Image27x27     string `json:"image_27x27,omitempty"`
	Image36x36     string `json:"image_36x36,omitempty"`
	Image60x60     string `json:"image_60x60,omitempty"`
	Image128x128   string `json:"image_128x128,omitempty"`
	Image1024x1024 string `json:"image_1024x1024,omitempty"`
}

// User is an Asana user resource (compact + default response shape).
//
// Reference: https://developers.asana.com/reference/users
type User struct {
	GID          string        `json:"gid"`
	ResourceType string        `json:"resource_type,omitempty"`
	Name         string        `json:"name,omitempty"`
	Email        string        `json:"email,omitempty"`         // default response only
	Photo        *Photo        `json:"photo,omitempty"`         // nullable object
	Workspaces   []Workspace   `json:"workspaces,omitempty"`    // default response
	CustomFields []CustomField `json:"custom_fields,omitempty"` // default response
}

// Workspace is an Asana workspace resource.
//
// Reference: https://developers.asana.com/reference/workspaces
type Workspace struct {
	GID            string   `json:"gid"`
	ResourceType   string   `json:"resource_type,omitempty"`
	Name           string   `json:"name,omitempty"`
	EmailDomains   []string `json:"email_domains,omitempty"`
	IsOrganization bool     `json:"is_organization,omitempty"`
}

// Team is an Asana team resource (compact shape).
//
// Members are NOT inlined on Team responses — fetch them via
// GET /teams/{team_gid}/users separately.
//
// Reference: https://developers.asana.com/reference/teams
type Team struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type,omitempty"`
	Name         string `json:"name,omitempty"`
}

// StatusUpdate is the current_status_update object on a project.
//
// Reference: https://developers.asana.com/reference/status-updates
type StatusUpdate struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type,omitempty"`
	Title        string `json:"title,omitempty"`
	Text         string `json:"text,omitempty"`
	HTMLText     string `json:"html_text,omitempty"`
	StatusType   string `json:"status_type,omitempty"` // on_track | at_risk | off_track | on_hold | complete
}

// Project is an Asana project resource (compact + base + response).
//
// Reference: https://developers.asana.com/reference/projects
type Project struct {
	GID                 string        `json:"gid"`
	ResourceType        string        `json:"resource_type,omitempty"`
	Name                string        `json:"name,omitempty"`
	Notes               string        `json:"notes,omitempty"`
	HTMLNotes           string        `json:"html_notes,omitempty"`
	Archived            bool          `json:"archived,omitempty"`
	Color               *string       `json:"color,omitempty"`           // nullable enum
	Icon                *string       `json:"icon,omitempty"`            // nullable enum
	DueOn               *CustomTime   `json:"due_on,omitempty"`
	StartOn             *CustomTime   `json:"start_on,omitempty"`
	DefaultView         string        `json:"default_view,omitempty"`    // list | board | calendar | timeline
	PrivacySetting      string        `json:"privacy_setting,omitempty"` // public_to_workspace | private_to_team | private
	CreatedAt           *time.Time    `json:"created_at,omitempty"`
	ModifiedAt          *time.Time    `json:"modified_at,omitempty"`
	Members             []User        `json:"members,omitempty"`
	Followers           []User        `json:"followers,omitempty"`
	Owner               *User         `json:"owner,omitempty"`
	Team                *Team         `json:"team,omitempty"`
	Workspace           *Workspace    `json:"workspace,omitempty"`
	PermalinkURL        string        `json:"permalink_url,omitempty"`
	Completed           bool          `json:"completed,omitempty"`
	CompletedAt         *time.Time    `json:"completed_at,omitempty"`
	CompletedBy         *User         `json:"completed_by,omitempty"`
	CustomFields        []CustomField `json:"custom_fields,omitempty"`
	CurrentStatusUpdate *StatusUpdate `json:"current_status_update,omitempty"`
}

// Section is a section within a project.
//
// Reference: https://developers.asana.com/reference/sections
type Section struct {
	GID          string     `json:"gid"`
	ResourceType string     `json:"resource_type,omitempty"`
	Name         string     `json:"name,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Project      *Project   `json:"project,omitempty"` // ProjectCompact object — NOT a GID string
}

// Tag is an Asana tag resource.
//
// Reference: https://developers.asana.com/reference/tags
type Tag struct {
	GID          string     `json:"gid"`
	ResourceType string     `json:"resource_type,omitempty"`
	Name         string     `json:"name,omitempty"`
	Color        *string    `json:"color,omitempty"`
	Notes        string     `json:"notes,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Followers    []User     `json:"followers,omitempty"`
	Workspace    *Workspace `json:"workspace,omitempty"`
	PermalinkURL string     `json:"permalink_url,omitempty"`
}

// Membership is a {project, section} pair on a Task. A task can belong to
// multiple projects, and within each project may live in a section.
type Membership struct {
	Project *Project `json:"project,omitempty"`
	Section *Section `json:"section,omitempty"`
}

// Task is an Asana task resource (compact + base + response).
//
// Reference: https://developers.asana.com/reference/tasks
type Task struct {
	GID             string        `json:"gid"`
	ResourceType    string        `json:"resource_type,omitempty"`
	ResourceSubtype string        `json:"resource_subtype,omitempty"` // default_task | milestone | approval | custom
	Name            string        `json:"name,omitempty"`
	Notes           string        `json:"notes,omitempty"`
	HTMLNotes       string        `json:"html_notes,omitempty"`
	Completed       bool          `json:"completed,omitempty"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty"`
	CompletedBy     *User         `json:"completed_by,omitempty"`
	DueOn           *CustomTime   `json:"due_on,omitempty"`
	DueAt           *CustomTime   `json:"due_at,omitempty"`
	StartOn         *CustomTime   `json:"start_on,omitempty"`
	StartAt         *CustomTime   `json:"start_at,omitempty"`
	Liked           bool          `json:"liked,omitempty"`
	NumLikes        int           `json:"num_likes,omitempty"`
	NumSubtasks     int           `json:"num_subtasks,omitempty"`
	Assignee        *User         `json:"assignee,omitempty"`
	AssigneeSection *Section      `json:"assignee_section,omitempty"`
	Parent          *Task         `json:"parent,omitempty"`
	Projects        []Project     `json:"projects,omitempty"`
	Tags            []Tag         `json:"tags,omitempty"`
	Workspace       *Workspace    `json:"workspace,omitempty"`
	Memberships     []Membership  `json:"memberships,omitempty"`
	Followers       []User        `json:"followers,omitempty"`
	CustomFields    []CustomField `json:"custom_fields,omitempty"`
	PermalinkURL    string        `json:"permalink_url,omitempty"`
	CreatedAt       *time.Time    `json:"created_at,omitempty"`
	ModifiedAt      *time.Time    `json:"modified_at,omitempty"`
	CreatedBy       *User         `json:"created_by,omitempty"`
}

// CustomField is an entry in the custom_fields[] array of a Task or Project.
// Polymorphic: which `*_value` field is populated depends on resource_subtype.
//
// Reference: https://developers.asana.com/reference/custom-fields
type CustomField struct {
	GID                string       `json:"gid"`
	ResourceType       string       `json:"resource_type,omitempty"`
	Name               string       `json:"name,omitempty"`
	ResourceSubtype    string       `json:"resource_subtype,omitempty"` // text | enum | multi_enum | number | date | people
	Type               string       `json:"type,omitempty"`             // deprecated alias of resource_subtype
	DisplayValue       *string      `json:"display_value,omitempty"`
	Description        string       `json:"description,omitempty"`
	Enabled            bool         `json:"enabled,omitempty"`
	NumberValue        *float64     `json:"number_value,omitempty"`
	TextValue          *string      `json:"text_value,omitempty"`
	HTMLTextValue      *string      `json:"html_text_value,omitempty"`
	EnumValue          *EnumOption  `json:"enum_value,omitempty"`
	EnumOptions        []EnumOption `json:"enum_options,omitempty"`
	MultiEnumValues    []EnumOption `json:"multi_enum_values,omitempty"`
	PeopleValue        []User       `json:"people_value,omitempty"`
	DateValue          *DateValue   `json:"date_value,omitempty"`
	PrivacySetting     string       `json:"privacy_setting,omitempty"`
	DefaultAccessLevel string       `json:"default_access_level,omitempty"`
}

// EnumOption is an option inside an enum or multi_enum custom field.
type EnumOption struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type,omitempty"`
	Name         string `json:"name,omitempty"`
	Color        string `json:"color,omitempty"`
	Enabled      bool   `json:"enabled,omitempty"`
}

// DateValue is the payload of a `date` custom field.
type DateValue struct {
	Date     string `json:"date,omitempty"`      // YYYY-MM-DD
	DateTime string `json:"date_time,omitempty"` // RFC3339
}

// Story is a comment or system event on a task.
//
// Comments live as `Story` resources with resource_subtype="comment_added".
// Use POST /tasks/{task_gid}/stories to add a comment.
//
// Reference: https://developers.asana.com/reference/stories
type Story struct {
	GID             string     `json:"gid"`
	ResourceType    string     `json:"resource_type,omitempty"`
	ResourceSubtype string     `json:"resource_subtype,omitempty"`
	Type            string     `json:"type,omitempty"` // comment | system
	Text            string     `json:"text,omitempty"`
	HTMLText        string     `json:"html_text,omitempty"`
	IsEditable      bool       `json:"is_editable,omitempty"`
	IsEdited        bool       `json:"is_edited,omitempty"`
	IsPinned        bool       `json:"is_pinned,omitempty"`
	StickerName     string     `json:"sticker_name,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	CreatedBy       *User      `json:"created_by,omitempty"`
	Liked           bool       `json:"liked,omitempty"`
	NumLikes        int        `json:"num_likes,omitempty"`
	Source          string     `json:"source,omitempty"` // web | email | mobile | api | unknown
}

// Attachment is a file attached to a task.
//
// Reference: https://developers.asana.com/reference/attachments
type Attachment struct {
	GID             string     `json:"gid"`
	ResourceType    string     `json:"resource_type,omitempty"`
	ResourceSubtype string     `json:"resource_subtype,omitempty"` // asana | dropbox | gdrive | onedrive | box | vimeo | external
	Name            string     `json:"name,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	DownloadURL     *string    `json:"download_url,omitempty"`
	PermanentURL    *string    `json:"permanent_url,omitempty"`
	ViewURL         *string    `json:"view_url,omitempty"`
	Host            string     `json:"host,omitempty"`
	Parent          *Task      `json:"parent,omitempty"`
	Size            int64      `json:"size,omitempty"`
	ConnectedToApp  bool       `json:"connected_to_app,omitempty"`
}

// TaskCreateRequest is the body for POST /tasks.
//
// Section placement requires a follow-up POST /sections/{section_gid}/addTask;
// the API silently ignores a "section" field in this body.
//
// Tag attachment requires a follow-up POST /tasks/{task_gid}/addTag per tag.
type TaskCreateRequest struct {
	Name      string   `json:"name,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	HTMLNotes string   `json:"html_notes,omitempty"`
	Projects  []string `json:"projects,omitempty"`
	Parent    string   `json:"parent,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	DueOn     string   `json:"due_on,omitempty"` // YYYY-MM-DD
	DueAt     string   `json:"due_at,omitempty"` // RFC3339
	StartOn   string   `json:"start_on,omitempty"`
	StartAt   string   `json:"start_at,omitempty"`
	Followers []string `json:"followers,omitempty"`
	Workspace string   `json:"workspace,omitempty"`
}

// TaskUpdateRequest is the body for PUT /tasks/{task_gid}.
//
// All fields use omitempty so callers can selectively update. To distinguish
// "don't change" from "clear to empty", consider using pointers in callers.
type TaskUpdateRequest struct {
	Name      string `json:"name,omitempty"`
	Notes     string `json:"notes,omitempty"`
	HTMLNotes string `json:"html_notes,omitempty"`
	Completed *bool  `json:"completed,omitempty"`
	Assignee  string `json:"assignee,omitempty"`
	DueOn     string `json:"due_on,omitempty"`
	DueAt     string `json:"due_at,omitempty"`
	StartOn   string `json:"start_on,omitempty"`
	StartAt   string `json:"start_at,omitempty"`
	Parent    string `json:"parent,omitempty"`
}

// StoryCreateRequest is the body for POST /tasks/{task_gid}/stories.
type StoryCreateRequest struct {
	Text     string `json:"text,omitempty"`
	HTMLText string `json:"html_text,omitempty"`
	IsPinned bool   `json:"is_pinned,omitempty"`
}

// SectionAddTaskRequest is the body for POST /sections/{section_gid}/addTask.
type SectionAddTaskRequest struct {
	Task         string `json:"task"`
	InsertBefore string `json:"insert_before,omitempty"`
	InsertAfter  string `json:"insert_after,omitempty"`
}

// WorkspaceUpdateRequest is the body for PUT /workspaces/{workspace_gid}.
type WorkspaceUpdateRequest struct {
	Name string `json:"name,omitempty"`
}
