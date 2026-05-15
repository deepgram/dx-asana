package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/deepgram/dx-asana/internal/asana"
	"github.com/deepgram/dx-asana/internal/ui"
	"github.com/spf13/cobra"
)

var (
	myTasksFields    string
	myTasksCompleted bool
	myTasksWorkspace string
)

var myTasksCmd = &cobra.Command{
	Use:   "my-tasks",
	Short: "Show tasks assigned to the authenticated user",
	Long: `List tasks across all projects in the user's "My Tasks" view.

Resolves the authenticated user's GID via /users/me, picks the first
workspace returned (override with --workspace), then fetches the user's
user_task_list. By default, completed tasks are hidden.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := asana.NewClient(token)

		me, err := client.GetMe(asana.WithOptFields("name", "workspaces.name"))
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		ws := myTasksWorkspace
		if ws == "" {
			if len(me.Workspaces) == 0 {
				err := fmt.Errorf("no workspaces returned for the authenticated user; pass --workspace explicitly")
				if jsonOutput {
					ui.PrintJSON(nil, err)
				} else {
					fmt.Fprintln(os.Stderr, "Error:", err)
				}
				return err
			}
			ws = me.Workspaces[0].GID
		}

		opts := []asana.Option{
			asana.WithAssignee(me.GID),
		}
		if myTasksFields != "" {
			opts = append(opts, asana.WithOptFields(splitFields(myTasksFields)...))
		} else {
			opts = append(opts, asana.WithOptFields(asana.DefaultTaskListFields...))
		}
		if !myTasksCompleted {
			opts = append(opts, asana.WithCompletedSince("now"))
		}

		tasks, err := client.GetTasksByWorkspace(ws, opts...)
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
				"count":      len(tasks),
				"user_gid":   me.GID,
				"user_name":  me.Name,
				"fetched_at": time.Now().Format(time.RFC3339),
			}
			ui.PrintJSONWithMeta(tasks, meta, nil)
		} else {
			fmt.Printf("📋 %s's tasks (%d)\n", me.Name, len(tasks))
			for _, t := range tasks {
				marker := "○"
				if t.Completed {
					marker = "●"
				}
				due := ""
				if t.DueOn != nil && !t.DueOn.IsZero() {
					due = " [due " + t.DueOn.Format("2006-01-02") + "]"
				}
				fmt.Printf("  %s %s%s\n", marker, t.Name, due)
			}
		}
		return nil
	},
}

func init() {
	myTasksCmd.Flags().StringVar(&myTasksFields, "fields", "", "Comma-separated opt_fields (overrides default)")
	myTasksCmd.Flags().BoolVar(&myTasksCompleted, "include-completed", false, "Include completed tasks")
	myTasksCmd.Flags().StringVar(&myTasksWorkspace, "workspace", "", "Workspace GID (defaults to first workspace from /users/me)")
}
