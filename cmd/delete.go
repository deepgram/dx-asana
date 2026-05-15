package cmd

import (
	"fmt"
	"os"

	"github.com/deepgram/dx-asana/internal/asana"
	"github.com/deepgram/dx-asana/internal/ui"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [task-id]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		client := asana.NewClient(token)

		err := client.DeleteTask(taskID)
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
				"action":  "deleted",
				"task_id": taskID,
			}
			ui.PrintJSONWithMeta(map[string]string{"status": "deleted"}, meta, nil)
		} else {
			fmt.Printf("✓ Task deleted\n")
		}

		return nil
	},
}
