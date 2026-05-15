package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	// Load should return empty config if file doesn't exist
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load returned nil config")
	}

	if cfg.Projects == nil {
		t.Error("Projects map should be initialized")
	}
}

func TestAddProject(t *testing.T) {
	cfg := &Config{
		Projects: make(map[string]ProjectConfig),
	}

	err := cfg.AddProject("test-project", "proj-123", "ws-456", "Test project")
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	proj, exists := cfg.Projects["test-project"]
	if !exists {
		t.Error("Project not found after adding")
	}

	if proj.ProjectID != "proj-123" {
		t.Errorf("unexpected project ID: %s", proj.ProjectID)
	}
}

func TestAddProjectValidation(t *testing.T) {
	cfg := &Config{
		Projects: make(map[string]ProjectConfig),
	}

	tests := []struct {
		name    string
		projName string
		projID  string
		wantErr bool
	}{
		{"valid", "project1", "proj-123", false},
		{"empty name", "", "proj-123", true},
		{"empty ID", "project2", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cfg.AddProject(tt.projName, tt.projID, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("AddProject error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCurrentProject(t *testing.T) {
	cfg := &Config{
		Projects: make(map[string]ProjectConfig),
	}

	if err := cfg.AddProject("proj1", "id-1", "", ""); err != nil {
		t.Fatalf("AddProject proj1 failed: %v", err)
	}
	if err := cfg.AddProject("proj2", "id-2", "", ""); err != nil {
		t.Fatalf("AddProject proj2 failed: %v", err)
	}

	err := cfg.SetCurrentProject("proj2")
	if err != nil {
		t.Fatalf("SetCurrentProject failed: %v", err)
	}

	if cfg.CurrentProject != "proj2" {
		t.Errorf("CurrentProject not set correctly: %s", cfg.CurrentProject)
	}
}

func TestSetCurrentProjectNotFound(t *testing.T) {
	cfg := &Config{
		Projects: make(map[string]ProjectConfig),
	}

	err := cfg.SetCurrentProject("nonexistent")
	if err == nil {
		t.Error("expected error when project not found")
	}
}

func TestRemoveProject(t *testing.T) {
	cfg := &Config{
		Projects:       make(map[string]ProjectConfig),
		CurrentProject: "proj1",
	}

	if err := cfg.AddProject("proj1", "id-1", "", ""); err != nil {
		t.Fatalf("AddProject proj1 failed: %v", err)
	}
	if err := cfg.AddProject("proj2", "id-2", "", ""); err != nil {
		t.Fatalf("AddProject proj2 failed: %v", err)
	}

	err := cfg.RemoveProject("proj1")
	if err != nil {
		t.Fatalf("RemoveProject failed: %v", err)
	}

	if _, exists := cfg.Projects["proj1"]; exists {
		t.Error("Project still exists after removal")
	}

	// Current project should change if we deleted it
	if cfg.CurrentProject == "proj1" {
		t.Error("CurrentProject not updated after deletion")
	}
}

func TestGetCurrentProject(t *testing.T) {
	cfg := &Config{
		Projects: make(map[string]ProjectConfig),
	}

	if err := cfg.AddProject("active", "proj-123", "", ""); err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}
	if err := cfg.SetCurrentProject("active"); err != nil {
		t.Fatalf("SetCurrentProject failed: %v", err)
	}

	proj := cfg.GetCurrentProject()
	if proj == nil {
		t.Fatal("GetCurrentProject returned nil")
	}

	if proj.Name != "active" {
		t.Errorf("unexpected project: %s", proj.Name)
	}
}

func TestListProjects(t *testing.T) {
	cfg := &Config{
		Projects: make(map[string]ProjectConfig),
	}

	if err := cfg.AddProject("proj1", "id-1", "", ""); err != nil {
		t.Fatalf("AddProject proj1 failed: %v", err)
	}
	if err := cfg.AddProject("proj2", "id-2", "", ""); err != nil {
		t.Fatalf("AddProject proj2 failed: %v", err)
	}

	projects := cfg.ListProjects()
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
}