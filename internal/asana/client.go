package asana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const defaultBaseURL = "https://app.asana.com/api/1.0"

// Client is a thin wrapper around Asana's REST API. Construct via NewClient.
// Methods on Client follow Asana's resource model (User, Task, Project,
// Workspace, Section, Story, Attachment, Tag) and accept variadic Option
// values for query parameters such as opt_fields, limit, and offset.
type Client struct {
	apiToken  string
	baseURL   string
	userAgent string
	http      *http.Client
}

// NewClient builds a Client. If apiToken is empty, ASANA_TOKEN is read from
// the environment as a fallback.
func NewClient(apiToken string) *Client {
	if apiToken == "" {
		apiToken = os.Getenv("ASANA_TOKEN")
	}
	return &Client{
		apiToken:  apiToken,
		baseURL:   defaultBaseURL,
		userAgent: "asana-cli/dev (+https://github.com/TheCoolRobot/asana-cli)",
		http:      &http.Client{},
	}
}

// SetBaseURL overrides the API base URL. Useful for testing against a fake
// server or a non-default Asana deployment.
func (c *Client) SetBaseURL(u string) { c.baseURL = u }

// SetUserAgent overrides the User-Agent header sent with each request.
func (c *Client) SetUserAgent(ua string) { c.userAgent = ua }

func (c *Client) do(method, endpoint string, body interface{}) ([]byte, error) {
	if c.apiToken == "" {
		return nil, fmt.Errorf("ASANA_TOKEN not set")
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}
	return respBody, nil
}

// GetMe retrieves the authenticated user.
func (c *Client) GetMe(opts ...Option) (*User, error) {
	body, err := c.do("GET", "/users/me"+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data *User `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// GetWorkspaces retrieves all workspaces the user has access to.
func (c *Client) GetWorkspaces(opts ...Option) ([]Workspace, error) {
	body, err := c.do("GET", "/workspaces"+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Workspace `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// GetProjects lists projects in a workspace.
func (c *Client) GetProjects(workspaceGID string, opts ...Option) ([]Project, error) {
	opts = append([]Option{WithWorkspace(workspaceGID)}, opts...)
	body, err := c.do("GET", "/projects"+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Project `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// GetTasks lists tasks in a project. By default, only the compact
// representation is returned; pass WithOptFields(...) to widen the response.
func (c *Client) GetTasks(projectGID string, opts ...Option) ([]Task, error) {
	body, err := c.do("GET", fmt.Sprintf("/projects/%s/tasks", projectGID)+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Task `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// GetTask retrieves a single task by GID.
func (c *Client) GetTask(taskGID string, opts ...Option) (*Task, error) {
	body, err := c.do("GET", fmt.Sprintf("/tasks/%s", taskGID)+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data *Task `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// CreateTask creates a new task. POST /tasks.
func (c *Client) CreateTask(req *TaskCreateRequest, opts ...Option) (*Task, error) {
	body, err := c.do("POST", "/tasks"+applyOptions(opts), map[string]interface{}{"data": req})
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data *Task `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// UpdateTask edits an existing task. PUT /tasks/{task_gid}.
func (c *Client) UpdateTask(taskGID string, req *TaskUpdateRequest, opts ...Option) (*Task, error) {
	body, err := c.do("PUT", fmt.Sprintf("/tasks/%s", taskGID)+applyOptions(opts), map[string]interface{}{"data": req})
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data *Task `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// CompleteTask marks a task as completed.
func (c *Client) CompleteTask(taskGID string) (*Task, error) {
	completed := true
	return c.UpdateTask(taskGID, &TaskUpdateRequest{Completed: &completed})
}

// DeleteTask removes a task. DELETE /tasks/{task_gid}.
func (c *Client) DeleteTask(taskGID string) error {
	_, err := c.do("DELETE", fmt.Sprintf("/tasks/%s", taskGID), nil)
	return err
}

// GetSections lists sections in a project.
func (c *Client) GetSections(projectGID string, opts ...Option) ([]Section, error) {
	body, err := c.do("GET", fmt.Sprintf("/projects/%s/sections", projectGID)+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Section `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// SearchTasks runs a typeahead text search across a workspace's tasks.
// Caller must pass WithText("...") for a useful query.
func (c *Client) SearchTasks(workspaceGID string, opts ...Option) ([]Task, error) {
	body, err := c.do("GET", fmt.Sprintf("/workspaces/%s/tasks/search", workspaceGID)+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Task `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// Search is a backward-compatible wrapper around SearchTasks for the simple
// "workspace + text query" case used by older callers.
func (c *Client) Search(workspaceGID, query string) ([]Task, error) {
	return c.SearchTasks(workspaceGID, WithText(query))
}

// GetUserTaskList retrieves a user's "My Tasks" list as tasks.
func (c *Client) GetUserTaskList(userGID string, opts ...Option) ([]Task, error) {
	body, err := c.do("GET", fmt.Sprintf("/users/%s/user_task_list", userGID)+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Task `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// GetUserTeams retrieves teams a user belongs to.
func (c *Client) GetUserTeams(userGID string, opts ...Option) ([]Team, error) {
	body, err := c.do("GET", fmt.Sprintf("/users/%s/teams", userGID)+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Team `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// UpdateWorkspace renames or otherwise updates a workspace.
func (c *Client) UpdateWorkspace(workspaceGID string, req *WorkspaceUpdateRequest, opts ...Option) (*Workspace, error) {
	body, err := c.do("PUT", fmt.Sprintf("/workspaces/%s", workspaceGID)+applyOptions(opts),
		map[string]interface{}{"data": req})
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data *Workspace `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// GetStories lists stories (comments and system events) on a task.
func (c *Client) GetStories(taskGID string, opts ...Option) ([]Story, error) {
	body, err := c.do("GET", fmt.Sprintf("/tasks/%s/stories", taskGID)+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Story `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// CreateStory adds a comment (or system event) to a task.
func (c *Client) CreateStory(taskGID string, req *StoryCreateRequest, opts ...Option) (*Story, error) {
	body, err := c.do("POST", fmt.Sprintf("/tasks/%s/stories", taskGID)+applyOptions(opts),
		map[string]interface{}{"data": req})
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data *Story `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// GetSubtasks lists subtasks for a parent task.
func (c *Client) GetSubtasks(parentGID string, opts ...Option) ([]Task, error) {
	body, err := c.do("GET", fmt.Sprintf("/tasks/%s/subtasks", parentGID)+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Task `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

// AddTaskToSection moves an existing task into a specific section. The Asana
// API requires this as a separate POST after the task has been created, since
// providing "section" in the create body is silently ignored.
func (c *Client) AddTaskToSection(sectionGID, taskGID string) error {
	_, err := c.do("POST", fmt.Sprintf("/sections/%s/addTask", sectionGID),
		map[string]interface{}{"data": SectionAddTaskRequest{Task: taskGID}})
	return err
}

// GetTasksByWorkspace lists tasks across a workspace, accepting filters via
// options (assignee, project, completed_since, etc.).
func (c *Client) GetTasksByWorkspace(workspaceGID string, opts ...Option) ([]Task, error) {
	body, err := c.do("GET", fmt.Sprintf("/workspaces/%s/tasks", workspaceGID)+applyOptions(opts), nil)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []Task `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}
