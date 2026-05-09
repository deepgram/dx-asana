package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/TheCoolRobot/asana-cli/internal/asana"
	"github.com/TheCoolRobot/asana-cli/internal/config"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
)

var (
	filterCompleted bool
	filterAssignee  string
	filterTag       string
	listFields      string
)

var listCmd = &cobra.Command{
	Use:   "list [project-id]",
	Short: "List tasks from a project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectGID := ""

		if len(args) > 0 {
			projectGID = args[0]
		} else {
			cfg, _ := config.Load()
			currentProj := cfg.GetCurrentProject()
			if currentProj == nil {
				return fmt.Errorf("no project ID provided and no current project set. Use: asana-cli list <project-id> or asana-cli config project switch <name>")
			}
			projectGID = currentProj.ProjectID
		}

		client := asana.NewClient(token)

		opts := []asana.Option{}
		if listFields != "" {
			opts = append(opts, asana.WithOptFields(splitFields(listFields)...))
		} else {
			opts = append(opts, asana.WithOptFields(asana.DefaultTaskListFields...))
		}
		if filterCompleted {
			opts = append(opts, asana.WithCompletedSince("now"))
		}
		if filterAssignee != "" {
			opts = append(opts, asana.WithAssignee(filterAssignee))
		}

		tasks, err := client.GetTasks(projectGID, opts...)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Println("Error:", err)
			}
			return err
		}

		// Convert to pointers
		taskPtrs := make([]*asana.Task, len(tasks))
		for i := range tasks {
			taskPtrs[i] = &tasks[i]
		}

		if jsonOutput {
			meta := map[string]interface{}{
				"count":      len(tasks),
				"project_id": projectGID,
				"fetched_at": time.Now().Format(time.RFC3339),
			}
			if filterCompleted {
				meta["filter_completed"] = true
			}
			if filterAssignee != "" {
				meta["filter_assignee"] = filterAssignee
			}
			ui.PrintJSONWithMeta(tasks, meta, nil)
		} else {
			ui.StartTUI(taskPtrs, client, projectGID)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVar(&filterCompleted, "completed", false, "Show only completed tasks")
	listCmd.Flags().StringVar(&filterAssignee, "assignee", "", "Filter by assignee ID")
	listCmd.Flags().StringVar(&filterTag, "tag", "", "Filter by tag")
	listCmd.Flags().StringVar(&listFields, "fields", "", "Comma-separated opt_fields (overrides default)")
}

func splitFields(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}