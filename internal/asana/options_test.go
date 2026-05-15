package asana

import (
	"net/url"
	"strings"
	"testing"
)

func TestApplyOptionsEmpty(t *testing.T) {
	if got := applyOptions(nil); got != "" {
		t.Errorf("nil opts: got %q want empty", got)
	}
	if got := applyOptions([]Option{}); got != "" {
		t.Errorf("empty opts: got %q want empty", got)
	}
}

func TestApplyOptionsSingle(t *testing.T) {
	got := applyOptions([]Option{WithLimit(50)})
	if got != "?limit=50" {
		t.Errorf("got %q want ?limit=50", got)
	}
}

func TestApplyOptionsCombines(t *testing.T) {
	got := applyOptions([]Option{
		WithText("bug fix"),
		WithLimit(25),
		WithAssignee("1234"),
	})
	if !strings.HasPrefix(got, "?") {
		t.Errorf("expected leading ?, got %q", got)
	}
	q, err := url.ParseQuery(got[1:])
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if q.Get("text") != "bug fix" {
		t.Errorf("text = %q", q.Get("text"))
	}
	if q.Get("limit") != "25" {
		t.Errorf("limit = %q", q.Get("limit"))
	}
	if q.Get("assignee") != "1234" {
		t.Errorf("assignee = %q", q.Get("assignee"))
	}
}

func TestWithOptFieldsJoinsCommaSeparated(t *testing.T) {
	got := applyOptions([]Option{WithOptFields("name", "notes", "assignee.name")})
	q, _ := url.ParseQuery(got[1:])
	if q.Get("opt_fields") != "name,notes,assignee.name" {
		t.Errorf("opt_fields = %q", q.Get("opt_fields"))
	}
}

func TestWithOptFieldsEmptyOmits(t *testing.T) {
	got := applyOptions([]Option{WithOptFields()})
	if got != "" {
		t.Errorf("empty fields should produce no query: got %q", got)
	}
}

func TestWithEmptyValuesAreOmitted(t *testing.T) {
	got := applyOptions([]Option{
		WithAssignee(""),
		WithProject(""),
		WithText(""),
		WithOffset(""),
	})
	if got != "" {
		t.Errorf("all-empty opts should produce no query: got %q", got)
	}
}

func TestWithLimitZeroOmits(t *testing.T) {
	got := applyOptions([]Option{WithLimit(0), WithLimit(-5)})
	if got != "" {
		t.Errorf("non-positive limit should be omitted: got %q", got)
	}
}

func TestWithRawArbitraryParam(t *testing.T) {
	got := applyOptions([]Option{WithRaw("custom_fields.123.value", "foo")})
	q, _ := url.ParseQuery(got[1:])
	if q.Get("custom_fields.123.value") != "foo" {
		t.Errorf("raw param: got %q", q.Get("custom_fields.123.value"))
	}
}

func TestDefaultTaskListFieldsIncludesEssentials(t *testing.T) {
	want := []string{"name", "completed", "due_on", "assignee.name", "permalink_url"}
	have := make(map[string]bool, len(DefaultTaskListFields))
	for _, f := range DefaultTaskListFields {
		have[f] = true
	}
	for _, w := range want {
		if !have[w] {
			t.Errorf("DefaultTaskListFields missing %q", w)
		}
	}
}
