package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ProjectConfig struct {
	Name        string `json:"name"`
	ProjectID   string `json:"project_id"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Description string `json:"description,omitempty"`
}

type Config struct {
	APIToken         string                    `json:"api_token"`
	CurrentProject   string                    `json:"current_project"` // Name of active project
	Projects         map[string]ProjectConfig  `json:"projects"`
	DefaultWorkspace string                    `json:"default_workspace"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".dx-asana", "config.json")
}

func Load() (*Config, error) {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Projects: make(map[string]ProjectConfig),
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Ensure Projects map is initialized
	if cfg.Projects == nil {
		cfg.Projects = make(map[string]ProjectConfig)
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path := GetConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func (c *Config) AddProject(name, projectID, workspaceID, description string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if projectID == "" {
		return fmt.Errorf("project ID cannot be empty")
	}

	c.Projects[name] = ProjectConfig{
		Name:        name,
		ProjectID:   projectID,
		WorkspaceID: workspaceID,
		Description: description,
	}

	// Set as current if it's the first project
	if c.CurrentProject == "" {
		c.CurrentProject = name
	}

	return c.Save()
}

func (c *Config) RemoveProject(name string) error {
	if _, exists := c.Projects[name]; !exists {
		return fmt.Errorf("project '%s' not found", name)
	}

	delete(c.Projects, name)

	// If we deleted the current project, switch to another one
	if c.CurrentProject == name {
		if len(c.Projects) > 0 {
			for name := range c.Projects {
				c.CurrentProject = name
				break
			}
		} else {
			c.CurrentProject = ""
		}
	}

	return c.Save()
}

func (c *Config) SetCurrentProject(name string) error {
	if _, exists := c.Projects[name]; !exists {
		return fmt.Errorf("project '%s' not found", name)
	}

	c.CurrentProject = name
	return c.Save()
}

func (c *Config) GetCurrentProject() *ProjectConfig {
	if c.CurrentProject == "" {
		return nil
	}

	if proj, exists := c.Projects[c.CurrentProject]; exists {
		return &proj
	}

	return nil
}

func (c *Config) ListProjects() []ProjectConfig {
	projects := make([]ProjectConfig, 0, len(c.Projects))
	for _, proj := range c.Projects {
		projects = append(projects, proj)
	}
	return projects
}

func GetAPIToken() string {
	if token := os.Getenv("ASANA_TOKEN"); token != "" {
		return token
	}

	cfg, _ := Load()
	if cfg != nil && cfg.APIToken != "" {
		return cfg.APIToken
	}

	return ""
}