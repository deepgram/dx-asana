package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/config"
)

type TaskListItem struct {
	Task     *asana.Task
	Selected bool
}

// addFormField represents which field is focused in the add form
type addFormField int

const (
	addFieldName addFormField = iota
	addFieldDescription
	addFieldDueDate
	addFieldPriority
	addFieldCount
)

type Model struct {
	items           []TaskListItem
	cursor          int
	filterCompleted bool
	//searchQuery     string
	sortBy          string // name, due_date, priority
	width           int
	height          int
	loading         bool
	message         string
	mode            string // "tasks", "projects", "confirm", or "add"
	projects        []config.ProjectConfig
	projectCursor   int
	currentProject  string
	client          *asana.Client
	projectGID      string
	confirmAction   string // "delete", "complete"
	confirmTaskGID  string
	confirmTaskName string

	// Add form state
	addFields      [addFieldCount]string
	addFocusField  addFormField
}

func NewModel(tasks []*asana.Task, client *asana.Client, projectGID string) Model {
	cfg, _ := config.Load()
	items := make([]TaskListItem, len(tasks))
	for i, t := range tasks {
		items[i] = TaskListItem{Task: t, Selected: false}
	}

	return Model{
		items:          items,
		cursor:         0,
		sortBy:         "name",
		width:          80,
		height:         24,
		mode:           "tasks",
		projects:       cfg.ListProjects(),
		projectCursor:  0,
		currentProject: cfg.CurrentProject,
		client:         client,
		projectGID:     projectGID,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.mode == "confirm" {
			return m.updateConfirmMode(msg)
		} else if m.mode == "projects" {
			return m.updateProjectMode(msg)
		} else if m.mode == "add" {
			return m.updateAddMode(msg)
		}
		return m.updateTaskMode(msg)
	}
	return m, nil
}

func (m Model) updateTaskMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}

	case " ":
		if m.cursor < len(m.items) {
			m.items[m.cursor].Selected = !m.items[m.cursor].Selected
		}

	case "c":
		if m.cursor < len(m.items) {
			m.confirmAction = "complete"
			m.confirmTaskGID = m.items[m.cursor].Task.GID
			m.confirmTaskName = m.items[m.cursor].Task.Name
			m.mode = "confirm"
		}

	case "d":
		if m.cursor < len(m.items) {
			m.confirmAction = "delete"
			m.confirmTaskGID = m.items[m.cursor].Task.GID
			m.confirmTaskName = m.items[m.cursor].Task.Name
			m.mode = "confirm"
		}

	case "a":
		// Enter add mode, reset form fields
		m.addFields = [addFieldCount]string{}
		m.addFocusField = addFieldName
		m.mode = "add"
		m.message = ""

	case "f":
		m.filterCompleted = !m.filterCompleted
		if m.filterCompleted {
			m.message = "Filtering: Hidden completed tasks"
		} else {
			m.message = "Showing: All tasks"
		}

	case "s":
		m.toggleSort()

	case "p":
		m.mode = "projects"
		m.message = "Project selector (↑↓ navigate, enter to switch, q to close)"

	case "e":
		if m.cursor < len(m.items) {
			m.message = fmt.Sprintf("Edit task: %s (not yet implemented in TUI)", m.items[m.cursor].Task.Name)
		}

	case "/":
		m.message = "Enter search query (not implemented in TUI demo)"

	case "enter":
		if m.cursor < len(m.items) {
			m.message = fmt.Sprintf("View details: %s (press 'v' to view full)", m.items[m.cursor].Task.Name)
		}
	}
	return m, nil
}

func (m Model) updateConfirmMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.loading = true
		if m.confirmAction == "complete" {
			_, err := m.client.CompleteTask(m.confirmTaskGID)
			if err != nil {
				m.message = fmt.Sprintf("❌ Error completing task: %v", err)
			} else {
				// Mark as completed in the list
				for i, item := range m.items {
					if item.Task.GID == m.confirmTaskGID {
						m.items[i].Task.Completed = true
						break
					}
				}
				m.message = fmt.Sprintf("✓ Completed: %s", m.confirmTaskName)
			}
		} else if m.confirmAction == "delete" {
			err := m.client.DeleteTask(m.confirmTaskGID)
			if err != nil {
				m.message = fmt.Sprintf("❌ Error deleting task: %v", err)
			} else {
				// Remove from the list
				for i, item := range m.items {
					if item.Task.GID == m.confirmTaskGID {
						m.items = append(m.items[:i], m.items[i+1:]...)
						if m.cursor >= len(m.items) && m.cursor > 0 {
							m.cursor--
						}
						break
					}
				}
				m.message = fmt.Sprintf("✓ Deleted: %s", m.confirmTaskName)
			}
		}
		m.loading = false
		m.mode = "tasks"
		m.confirmAction = ""
		m.confirmTaskGID = ""
		m.confirmTaskName = ""

	case "n", "esc":
		m.mode = "tasks"
		m.confirmAction = ""
		m.confirmTaskGID = ""
		m.confirmTaskName = ""
		m.message = "Cancelled"
	}
	return m, nil
}

func (m Model) updateProjectMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.mode = "tasks"
		m.message = ""

	case "up", "k":
		if m.projectCursor > 0 {
			m.projectCursor--
		}

	case "down", "j":
		if m.projectCursor < len(m.projects)-1 {
			m.projectCursor++
		}

	case "enter":
		if m.projectCursor < len(m.projects) {
			selectedProject := m.projects[m.projectCursor]
			cfg, _ := config.Load()
			cfg.SetCurrentProject(selectedProject.Name)
			m.currentProject = selectedProject.Name
			m.message = fmt.Sprintf("✓ Switched to: %s", selectedProject.Name)
			m.mode = "tasks"
		}
	}
	return m, nil
}

func (m Model) updateAddMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = "tasks"
		m.message = "Cancelled"
		return m, nil

	case "tab", "down":
		m.addFocusField = (m.addFocusField + 1) % addFieldCount
		return m, nil

	case "shift+tab", "up":
		m.addFocusField = (m.addFocusField - 1 + addFieldCount) % addFieldCount
		return m, nil

	case "enter":
		// Submit the form if name is filled in
		if strings.TrimSpace(m.addFields[addFieldName]) == "" {
			m.message = "❌ Task name is required"
			return m, nil
		}

		if m.projectGID == "" {
			m.message = "❌ No project selected — cannot create task"
			return m, nil
		}

		m.loading = true
		req := &asana.TaskCreateRequest{
			Name:     strings.TrimSpace(m.addFields[addFieldName]),
			Notes:    strings.TrimSpace(m.addFields[addFieldDescription]),
			Projects: []string{m.projectGID},
			DueOn:    strings.TrimSpace(m.addFields[addFieldDueDate]),
		}

		task, err := m.client.CreateTask(req)
		m.loading = false

		if err != nil {
			m.message = fmt.Sprintf("❌ Error creating task: %v", err)
			return m, nil
		}

		// Prepend new task to the list
		newItem := TaskListItem{Task: task, Selected: false}
		m.items = append([]TaskListItem{newItem}, m.items...)
		m.cursor = 0
		m.mode = "tasks"
		m.message = fmt.Sprintf("✓ Created: %s", task.Name)
		return m, nil

	case "backspace":
		field := &m.addFields[m.addFocusField]
		if len(*field) > 0 {
			*field = (*field)[:len(*field)-1]
		}
		return m, nil

	default:
		// Only accept printable characters (single rune keys)
		key := msg.String()
		if len(key) == 1 {
			m.addFields[m.addFocusField] += key
		}
	}

	return m, nil
}

func (m Model) toggleSort() {
	switch m.sortBy {
	case "name":
		m.sortBy = "due_date"
		m.message = "Sorting by: Due Date"
	case "due_date":
		m.sortBy = "priority"
		m.message = "Sorting by: Priority"
	default:
		m.sortBy = "name"
		m.message = "Sorting by: Name"
	}
}

func (m Model) View() string {
	if m.mode == "confirm" {
		return m.viewConfirmation()
	} else if m.mode == "projects" {
		return m.viewProjects()
	} else if m.mode == "add" {
		return m.viewAddForm()
	}
	return m.viewTasks()
}

func (m Model) viewConfirmation() string {
	var sb strings.Builder

	sb.WriteString(StyleError.Render("⚠️  CONFIRM ACTION") + "\n\n")

	if m.confirmAction == "complete" {
		sb.WriteString(fmt.Sprintf("Mark task as complete?\n\n  %s\n\n", m.confirmTaskName))
	} else if m.confirmAction == "delete" {
		sb.WriteString(fmt.Sprintf("Delete this task permanently?\n\n  %s\n\n", m.confirmTaskName))
	}

	if m.loading {
		sb.WriteString(StyleWarning.Render("Processing...") + "\n")
	} else {
		sb.WriteString(StyleDim.Render("[y] Confirm  [n] Cancel") + "\n")
	}

	return sb.String()
}

func (m Model) viewAddForm() string {
	var sb strings.Builder

	sb.WriteString(StyleTitle.Render("➕ New Task") + "\n\n")

	type formRow struct {
		label string
		field addFormField
		hint  string
	}

	rows := []formRow{
		{"Name*", addFieldName, "required"},
		{"Description", addFieldDescription, "optional"},
		{"Due Date", addFieldDueDate, "YYYY-MM-DD, optional"},
		{"Priority", addFieldPriority, "high / medium / low, optional"},
	}

	for _, row := range rows {
		focused := m.addFocusField == row.field
		value := m.addFields[row.field]

		labelStyle := StyleDim
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("59")).
			Padding(0, 1).
			Width(40)

		if focused {
			labelStyle = StyleSuccess
			inputStyle = inputStyle.
				BorderForeground(ColorPrimary)
		}

		displayValue := value
		if focused {
			displayValue = value + "█" // cursor
		} else if value == "" {
			displayValue = StyleDim.Render(row.hint)
		}

		sb.WriteString(labelStyle.Render(fmt.Sprintf("%-12s", row.label)))
		sb.WriteString(inputStyle.Render(displayValue))
		sb.WriteString("\n\n")
	}

	sb.WriteString(StyleDim.Render(strings.Repeat("─", m.width)) + "\n")

	if m.loading {
		sb.WriteString(StyleWarning.Render("Creating task...") + "\n")
	} else {
		sb.WriteString(StyleDim.Render("[tab/↑↓] move between fields  [enter] create  [esc] cancel") + "\n")
	}

	if m.message != "" {
		sb.WriteString("\n" + StyleError.Render(m.message) + "\n")
	}

	return sb.String()
}

func (m Model) viewTasks() string {
	if len(m.items) == 0 {
		var sb strings.Builder
		sb.WriteString(StyleTitle.Render(fmt.Sprintf("📋 Asana Tasks - %s", m.currentProject)) + "\n\n")
		sb.WriteString(StyleDim.Render("No tasks found. Press [a] to add one.") + "\n\n")
		sb.WriteString(StyleDim.Render(strings.Repeat("─", m.width)) + "\n")
		sb.WriteString(StyleDim.Render("[a] add task  [q] quit") + "\n")
		if m.message != "" {
			sb.WriteString("\n" + StyleSuccess.Render(m.message) + "\n")
		}
		return sb.String()
	}

	var sb strings.Builder

	// Title with current project
	title := fmt.Sprintf("📋 Asana Tasks - %s", m.currentProject)
	sb.WriteString(StyleTitle.Render(title) + "\n\n")

	// Status line
	statusLine := fmt.Sprintf("[%d/%d tasks]", len(m.items), len(m.items))
	if m.filterCompleted {
		statusLine += " [filtered]"
	}
	statusLine += fmt.Sprintf(" [sort: %s]", m.sortBy)
	sb.WriteString(StyleDim.Render(statusLine) + "\n\n")

	// Task list
	for i, item := range m.items {
		if m.filterCompleted && item.Task.Completed {
			continue
		}

		cursor := "  "
		if m.cursor == i {
			cursor = "→ "
		}

		checkmark := "☐"
		if item.Task.Completed {
			checkmark = "☑"
		}

		selected := " "
		if item.Selected {
			selected = "✓"
		}

		// Task name with style
		taskName := item.Task.Name
		if item.Task.Completed {
			taskName = StyleCompleted.Render(taskName)
		} else if m.cursor == i {
			taskName = StyleSelected.Render(taskName)
		}

		priority := ""

		dueDate := ""
		if item.Task.DueOn != nil && !item.Task.DueOn.IsZero() {
			daysUntil := int(time.Until(item.Task.DueOn.Time).Hours() / 24)
			if daysUntil < 0 {
				dueDate = fmt.Sprintf(" %s", StyleError.Render(fmt.Sprintf("[%d days overdue]", -daysUntil)))
			} else if daysUntil == 0 {
				dueDate = fmt.Sprintf(" %s", StyleWarning.Render("[Due today]"))
			} else if daysUntil <= 3 {
				dueDate = fmt.Sprintf(" %s", StyleWarning.Render(fmt.Sprintf("[%d days left]", daysUntil)))
			}
		}

		line := fmt.Sprintf("%s[%s] %s %s%s%s\n", cursor, selected, checkmark, taskName, priority, dueDate)
		sb.WriteString(line)
	}

	// Footer
	sb.WriteString("\n" + StyleDim.Render(strings.Repeat("─", m.width)) + "\n")
	helpText := "[↑↓] navigate  [space] select  [a] add  [c] complete  [d] delete  [f] filter  [s] sort  [p] projects  [q] quit"
	if m.width < len(helpText) {
		helpText = "[↑↓] nav  [a] add  [c] done  [d] del  [p] proj  [q] quit"
	}
	sb.WriteString(StyleDim.Render(helpText) + "\n")

	if m.message != "" {
		sb.WriteString("\n" + StyleSuccess.Render(m.message) + "\n")
	}

	return sb.String()
}

func (m Model) viewProjects() string {
	if len(m.projects) == 0 {
		return StyleError.Render("No projects configured")
	}

	var sb strings.Builder

	sb.WriteString(StyleTitle.Render("🗂️  Projects") + "\n\n")
	sb.WriteString(StyleDim.Render("Select a project to switch to:") + "\n\n")

	for i, proj := range m.projects {
		cursor := "  "
		if m.projectCursor == i {
			cursor = "→ "
		}

		style := lipgloss.NewStyle()
		if m.projectCursor == i {
			style = StyleSelected
		}

		current := ""
		if proj.Name == m.currentProject {
			current = " ✓"
		}

		line := fmt.Sprintf("%s%s%s\n", cursor, style.Render(proj.Name), current)
		sb.WriteString(line)

		if proj.Description != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", StyleDim.Render(proj.Description)))
		}
	}

	sb.WriteString("\n" + StyleDim.Render(strings.Repeat("─", m.width)) + "\n")
	sb.WriteString(StyleDim.Render("[↑↓] navigate  [enter] switch  [q] close") + "\n")

	if m.message != "" {
		sb.WriteString("\n" + StyleSuccess.Render(m.message) + "\n")
	}

	return sb.String()
}

func StartTUI(tasks []*asana.Task, client *asana.Client, projectGID string) {
	m := NewModel(tasks, client, projectGID)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}