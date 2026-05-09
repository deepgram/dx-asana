package cmd

import (
	"fmt"

	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	searchFields    string
	searchLimit     int
	searchAssignee  string
	searchProject   string
	searchCompleted bool
)

var searchCmd = &cobra.Command{
	Use:   "search [workspace-id] [query]",
	Short: "Search for tasks",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceGID := args[0]
		query := args[1]
		client := asana.NewClient(token)

		opts := []asana.Option{asana.WithText(query)}
		if searchFields != "" {
			opts = append(opts, asana.WithOptFields(splitFields(searchFields)...))
		} else {
			opts = append(opts, asana.WithOptFields(asana.DefaultTaskListFields...))
		}
		if searchLimit > 0 {
			opts = append(opts, asana.WithLimit(searchLimit))
		}
		if searchAssignee != "" {
			opts = append(opts, asana.WithAssignee(searchAssignee))
		}
		if searchProject != "" {
			opts = append(opts, asana.WithProject(searchProject))
		}
		if searchCompleted {
			opts = append(opts, asana.WithRaw("completed", "true"))
		}

		tasks, err := client.SearchTasks(workspaceGID, opts...)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Println("Error:", err)
			}
			return err
		}

		taskPtrs := make([]*asana.Task, len(tasks))
		for i := range tasks {
			taskPtrs[i] = &tasks[i]
		}

		if jsonOutput {
			meta := map[string]interface{}{
				"count": len(tasks),
				"query": query,
			}
			ui.PrintJSONWithMeta(tasks, meta, nil)
		} else {
			ui.StartTUI(taskPtrs, client, "")
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchFields, "fields", "", "Comma-separated opt_fields (overrides default)")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 0, "Maximum results to return (1-100)")
	searchCmd.Flags().StringVar(&searchAssignee, "assignee", "", "Filter by assignee user GID")
	searchCmd.Flags().StringVar(&searchProject, "project", "", "Filter by project GID")
	searchCmd.Flags().BoolVar(&searchCompleted, "completed", false, "Include only completed tasks")
}
