package asana

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCustomTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		checkFn func(*CustomTime) bool
	}{
		{
			name:    "null",
			json:    `null`,
			wantErr: false,
			checkFn: func(ct *CustomTime) bool { return ct.IsZero() },
		},
		{
			name:    "empty string",
			json:    `""`,
			wantErr: false,
			checkFn: func(ct *CustomTime) bool { return ct.IsZero() },
		},
		{
			name:    "date only",
			json:    `"2026-02-19"`,
			wantErr: false,
			checkFn: func(ct *CustomTime) bool {
				return ct.Year() == 2026 && ct.Month() == 2 && ct.Day() == 19
			},
		},
		{
			name:    "RFC3339",
			json:    `"2026-02-19T10:30:00Z"`,
			wantErr: false,
			checkFn: func(ct *CustomTime) bool {
				return ct.Year() == 2026
			},
		},
		{
			name:    "invalid format",
			json:    `"not-a-date"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct CustomTime
			err := json.Unmarshal([]byte(tt.json), &ct)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.checkFn != nil && !tt.checkFn(&ct) {
				t.Error("CustomTime value doesn't match expected")
			}
		})
	}
}

func TestCustomTimeIsZero(t *testing.T) {
	var ct CustomTime
	if !ct.IsZero() {
		t.Error("empty CustomTime should be zero")
	}

	ct.Time = time.Now()
	if ct.IsZero() {
		t.Error("CustomTime with time should not be zero")
	}
}

// fixtureUserMe is the actual response shape from GET /users/me, captured from
// the live Asana API. Includes the multi-size photo object that previously
// caused unmarshal failure.
const fixtureUserMe = `{
  "data": {
    "gid": "1201116150354130",
    "resource_type": "user",
    "name": "Luke Oliff",
    "email": "luke.oliff@deepgram.com",
    "photo": {
      "image_21x21": "https://asanausercontent.com/.../21x21.png",
      "image_27x27": "https://asanausercontent.com/.../27x27.png",
      "image_36x36": "https://asanausercontent.com/.../36x36.png",
      "image_60x60": "https://asanausercontent.com/.../60x60.png",
      "image_128x128": "https://asanausercontent.com/.../128x128.png"
    },
    "workspaces": [
      {
        "gid": "411927538413705",
        "resource_type": "workspace",
        "name": "deepgram.com"
      }
    ]
  }
}`

func TestUserUnmarshalsPhotoAsObject(t *testing.T) {
	var envelope struct {
		Data User `json:"data"`
	}
	if err := json.Unmarshal([]byte(fixtureUserMe), &envelope); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	u := envelope.Data
	if u.GID != "1201116150354130" {
		t.Errorf("GID = %q, want 1201116150354130", u.GID)
	}
	if u.ResourceType != "user" {
		t.Errorf("ResourceType = %q, want user", u.ResourceType)
	}
	if u.Name != "Luke Oliff" {
		t.Errorf("Name = %q, want Luke Oliff", u.Name)
	}
	if u.Email != "luke.oliff@deepgram.com" {
		t.Errorf("Email = %q, want luke.oliff@deepgram.com", u.Email)
	}
	if u.Photo == nil {
		t.Fatal("Photo is nil; expected populated multi-size object")
	}
	if u.Photo.Image21x21 == "" {
		t.Error("Photo.Image21x21 is empty")
	}
	if u.Photo.Image128x128 == "" {
		t.Error("Photo.Image128x128 is empty")
	}
	if len(u.Workspaces) != 1 {
		t.Fatalf("Workspaces len = %d, want 1", len(u.Workspaces))
	}
	if u.Workspaces[0].Name != "deepgram.com" {
		t.Errorf("Workspaces[0].Name = %q, want deepgram.com", u.Workspaces[0].Name)
	}
}

func TestUserUnmarshalsNullPhoto(t *testing.T) {
	body := `{
  "data": {
    "gid": "1234",
    "resource_type": "user",
    "name": "No Avatar",
    "photo": null
  }
}`
	var envelope struct {
		Data User `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("unmarshal failed on null photo: %v", err)
	}
	if envelope.Data.Photo != nil {
		t.Errorf("Photo should be nil for null avatar, got %+v", envelope.Data.Photo)
	}
}

func TestTaskUnmarshalsNotesAndNullableTimes(t *testing.T) {
	body := `{
  "data": {
    "gid": "9999",
    "resource_type": "task",
    "name": "Test task",
    "notes": "These are the notes",
    "html_notes": "<body>These are the notes</body>",
    "completed": false,
    "completed_at": null,
    "due_on": "2026-12-31",
    "due_at": null,
    "permalink_url": "https://app.asana.com/0/123/9999",
    "memberships": [
      {
        "project": {"gid": "111", "resource_type": "project", "name": "P1"},
        "section": {"gid": "222", "resource_type": "section", "name": "Backlog"}
      }
    ],
    "assignee": {"gid": "1", "resource_type": "user", "name": "Luke"}
  }
}`
	var envelope struct {
		Data Task `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	tk := envelope.Data
	if tk.Notes != "These are the notes" {
		t.Errorf("Notes = %q", tk.Notes)
	}
	if tk.HTMLNotes != "<body>These are the notes</body>" {
		t.Errorf("HTMLNotes = %q", tk.HTMLNotes)
	}
	if tk.CompletedAt != nil {
		t.Error("CompletedAt should be nil for null")
	}
	if tk.DueOn == nil || tk.DueOn.IsZero() {
		t.Error("DueOn should be populated from 2026-12-31")
	}
	if tk.DueAt != nil && !tk.DueAt.IsZero() {
		t.Error("DueAt should be nil/zero for null")
	}
	if tk.PermalinkURL == "" {
		t.Error("PermalinkURL should be populated")
	}
	if len(tk.Memberships) != 1 {
		t.Fatalf("Memberships len = %d", len(tk.Memberships))
	}
	if tk.Memberships[0].Project.Name != "P1" {
		t.Errorf("Memberships[0].Project.Name = %q", tk.Memberships[0].Project.Name)
	}
	if tk.Memberships[0].Section.Name != "Backlog" {
		t.Errorf("Memberships[0].Section.Name = %q", tk.Memberships[0].Section.Name)
	}
	if tk.Assignee == nil || tk.Assignee.Name != "Luke" {
		t.Error("Assignee not populated")
	}
}

func TestSectionUnmarshalsProjectAsObject(t *testing.T) {
	body := `{
  "data": {
    "gid": "555",
    "resource_type": "section",
    "name": "Done",
    "project": {"gid": "111", "resource_type": "project", "name": "Roadmap"}
  }
}`
	var envelope struct {
		Data Section `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if envelope.Data.Project == nil {
		t.Fatal("Project should be populated as object")
	}
	if envelope.Data.Project.Name != "Roadmap" {
		t.Errorf("Project.Name = %q", envelope.Data.Project.Name)
	}
}

func TestProjectUnmarshalsCurrentStatusUpdate(t *testing.T) {
	body := `{
  "data": {
    "gid": "111",
    "resource_type": "project",
    "name": "P1",
    "color": "light-green",
    "current_status_update": {
      "gid": "999",
      "resource_type": "status_update",
      "title": "Status: On Track",
      "text": "Things are good",
      "status_type": "on_track"
    }
  }
}`
	var envelope struct {
		Data Project `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if envelope.Data.Color == nil || *envelope.Data.Color != "light-green" {
		t.Errorf("Color = %v", envelope.Data.Color)
	}
	if envelope.Data.CurrentStatusUpdate == nil {
		t.Fatal("CurrentStatusUpdate should be populated")
	}
	if envelope.Data.CurrentStatusUpdate.StatusType != "on_track" {
		t.Errorf("StatusType = %q", envelope.Data.CurrentStatusUpdate.StatusType)
	}
}

func TestProjectUnmarshalsNullColor(t *testing.T) {
	body := `{
  "data": {
    "gid": "111",
    "resource_type": "project",
    "name": "P1",
    "color": null
  }
}`
	var envelope struct {
		Data Project `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if envelope.Data.Color != nil {
		t.Errorf("Color should be nil for null, got %v", envelope.Data.Color)
	}
}

func TestTaskCreateRequestSerializesNotesNotDescription(t *testing.T) {
	req := &TaskCreateRequest{
		Name:  "New task",
		Notes: "Description body",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	got := string(b)
	if got != `{"name":"New task","notes":"Description body"}` {
		t.Errorf("serialized = %s", got)
	}
}
