package cmd

import (
	"fmt"
	"os"

	"github.com/deepgram/dx-asana/internal/asana"
	"github.com/deepgram/dx-asana/internal/ui"
	"github.com/spf13/cobra"
)

var (
	updateName        string
	updateDescription string
	updateAssignee    string
	updateDueDate     string
	updatePriority    string
)

var updateCmd = &cobra.Command{
	Use:   "update [task-id]",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskGID := args[0]
		if err := validateDueDate(updateDueDate); err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}
		client := asana.NewClient(token)

		req := &asana.TaskUpdateRequest{
			Name:     updateName,
			Notes:    updateDescription,
			Assignee: updateAssignee,
			DueOn:    updateDueDate,
		}

		if updatePriority != "" && !jsonOutput {
			fmt.Println("⚠ --priority is not yet supported (Asana priority is a custom_field; coming in a later release)")
		}

		task, err := client.UpdateTask(taskGID, req)
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
				"action":   "updated",
				"task_gid": taskGID,
			}
			ui.PrintJSONWithMeta(task, meta, nil)
		} else {
			fmt.Printf("✓ Task updated: %s\n", task.Name)
		}

		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateName, "name", "", "New task name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "New task description (mapped to Asana 'notes')")
	updateCmd.Flags().StringVar(&updateAssignee, "assignee", "", "New assignee user GID")
	updateCmd.Flags().StringVar(&updateDueDate, "due", "", "New due date (YYYY-MM-DD)")
	updateCmd.Flags().StringVar(&updatePriority, "priority", "", "New priority (not yet supported; see warning at runtime)")
}
