package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	taskName        string
	taskDescription string
	taskAssignee    string
	taskDueDate     string
	taskPriority    string
	taskSection     string
	taskParent      string
)

var createCmd = &cobra.Command{
	Use:   "create [project-id]",
	Short: "Create a new task",
	Long: `Create a new task in a project, or as a subtask of an existing task.

  asana-cli create <project-gid> --name "..."
  asana-cli create --parent <task-gid> --name "Subtask"   # project-id not required`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if taskName == "" {
			err := fmt.Errorf("task name is required")
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error: task name is required")
			}
			return err
		}
		if err := validateDueDate(taskDueDate); err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		projectGID := ""
		if len(args) > 0 {
			projectGID = args[0]
		}
		if projectGID == "" && taskParent == "" {
			err := fmt.Errorf("either a project GID positional or --parent <task-gid> is required")
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}
		client := asana.NewClient(token)

		req := &asana.TaskCreateRequest{
			Name:     taskName,
			Notes:    taskDescription,
			Assignee: taskAssignee,
			DueOn:    taskDueDate,
			Parent:   taskParent,
		}
		if taskParent == "" {
			req.Projects = []string{projectGID}
		}

		if taskPriority != "" && !jsonOutput {
			fmt.Println("⚠ --priority is not yet supported (Asana priority is a custom_field; coming in a later release)")
		}
		if taskSection != "" && !jsonOutput {
			fmt.Println("⚠ --section is not yet supported (requires POST /sections/{gid}/addTask; coming in a later release)")
		}

		task, err := client.CreateTask(req)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		if jsonOutput {
			meta := map[string]interface{}{
				"action":     "created",
				"project_id": projectGID,
			}
			ui.PrintJSONWithMeta(task, meta, nil)
		} else {
			fmt.Printf("✓ Task created: %s\n", task.Name)
			fmt.Printf("  GID: %s\n", task.GID)
			if task.PermalinkURL != "" {
				fmt.Printf("  Link: %s\n", task.PermalinkURL)
			}
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringVar(&taskName, "name", "", "Task name (required)")
	createCmd.Flags().StringVar(&taskDescription, "description", "", "Task description (mapped to Asana 'notes')")
	createCmd.Flags().StringVar(&taskAssignee, "assignee", "", "Assignee user GID")
	createCmd.Flags().StringVar(&taskDueDate, "due", "", "Due date (YYYY-MM-DD)")
	createCmd.Flags().StringVar(&taskPriority, "priority", "", "Priority (not yet supported; see warning at runtime)")
	createCmd.Flags().StringVar(&taskSection, "section", "", "Section GID")
	createCmd.Flags().StringVar(&taskParent, "parent", "", "Parent task GID (creates a subtask; project arg is ignored when set)")
	if err := createCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf(err.Error())
	}
}
