package cmd

import (
	"fmt"
	"os"

	"github.com/deepgram/dx-asana/internal/asana"
	"github.com/deepgram/dx-asana/internal/ui"
	"github.com/spf13/cobra"
)

var viewFields string

var viewCmd = &cobra.Command{
	Use:   "view [task-id]",
	Short: "View task details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskGID := args[0]
		client := asana.NewClient(token)

		opts := []asana.Option{}
		if viewFields != "" {
			opts = append(opts, asana.WithOptFields(splitFields(viewFields)...))
		} else {
			opts = append(opts, asana.WithOptFields(asana.DefaultTaskFields...))
		}

		task, err := client.GetTask(taskGID, opts...)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		if jsonOutput {
			ui.PrintJSON(task, nil)
		} else {
			fmt.Printf("📋 %s\n", task.Name)
			fmt.Printf("   GID: %s\n", task.GID)
			fmt.Printf("   Completed: %v\n", task.Completed)
			if task.DueOn != nil && !task.DueOn.IsZero() {
				fmt.Printf("   Due: %s\n", task.DueOn.Format("2006-01-02"))
			}
			if task.Notes != "" {
				fmt.Printf("   Notes: %s\n", task.Notes)
			}
			if task.Assignee != nil {
				fmt.Printf("   Assigned to: %s\n", task.Assignee.Name)
			}
			if len(task.Tags) > 0 {
				fmt.Printf("   Tags: ")
				for i, tag := range task.Tags {
					if i > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s", tag.Name)
				}
				fmt.Printf("\n")
			}
			if task.PermalinkURL != "" {
				fmt.Printf("   Link: %s\n", task.PermalinkURL)
			}
		}

		return nil
	},
}

func init() {
	viewCmd.Flags().StringVar(&viewFields, "fields", "", "Comma-separated opt_fields (overrides default)")
}
