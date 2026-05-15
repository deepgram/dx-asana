package cmd

import (
	"fmt"
	"os"

	"github.com/deepgram/dx-asana/internal/asana"
	"github.com/deepgram/dx-asana/internal/ui"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete [task-id]",
	Short: "Mark a task as complete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		client := asana.NewClient(token)

		task, err := client.CompleteTask(taskID)
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
				"action":  "completed",
				"task_id": taskID,
			}
			ui.PrintJSONWithMeta(task, meta, nil)
		} else {
			fmt.Printf("✓ Task completed: %s\n", task.Name)
		}

		return nil
	},
}
