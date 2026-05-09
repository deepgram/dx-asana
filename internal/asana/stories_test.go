package asana

import (
	"encoding/json"
	"testing"
)

func TestStoryUnmarshal(t *testing.T) {
	body := `{
  "data": [
    {
      "gid": "111",
      "resource_type": "story",
      "resource_subtype": "comment_added",
      "type": "comment",
      "text": "Looks great!",
      "html_text": "<body>Looks great!</body>",
      "is_pinned": false,
      "created_by": {"gid": "222", "resource_type": "user", "name": "Luke"},
      "created_at": "2026-05-09T10:00:00Z"
    },
    {
      "gid": "112",
      "resource_type": "story",
      "resource_subtype": "marked_complete",
      "type": "system",
      "text": "Marked complete",
      "created_at": "2026-05-09T11:00:00Z"
    }
  ]
}`
	var envelope struct {
		Data []Story `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(envelope.Data) != 2 {
		t.Fatalf("len = %d, want 2", len(envelope.Data))
	}
	if envelope.Data[0].Type != "comment" {
		t.Errorf("Data[0].Type = %q", envelope.Data[0].Type)
	}
	if envelope.Data[0].CreatedBy == nil || envelope.Data[0].CreatedBy.Name != "Luke" {
		t.Errorf("Data[0].CreatedBy.Name = %v", envelope.Data[0].CreatedBy)
	}
	if envelope.Data[1].Type != "system" {
		t.Errorf("Data[1].Type = %q", envelope.Data[1].Type)
	}
}

func TestStoryCreateRequestSerializes(t *testing.T) {
	req := &StoryCreateRequest{
		Text:     "Status update",
		IsPinned: true,
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	if got != `{"text":"Status update","is_pinned":true}` {
		t.Errorf("got %s", got)
	}
}
