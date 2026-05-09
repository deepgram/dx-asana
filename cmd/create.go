package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
)

var (
	taskName        string
	taskDescription string
	taskAssignee    string
	taskDueDate     string
	taskPriority    string
	taskSection     string
)

var createCmd = &cobra.Command{
	Use:   "create [project-id]",
	Short: "Create a new task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if taskName == "" {
			if jsonOutput {
				ui.PrintJSON(nil, fmt.Errorf("task name is required"))
			} else {
				fmt.Println("Error: task name is required")
			}
			return fmt.Errorf("task name required")
		}

		projectGID := args[0]
		client := asana.NewClient(token)

		req := &asana.TaskCreateRequest{
			Name:     taskName,
			Notes:    taskDescription,
			Projects: []string{projectGID},
			Assignee: taskAssignee,
			DueOn:    taskDueDate,
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
				fmt.Println("Error:", err)
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
	if err := createCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf(err.Error())
	}
}