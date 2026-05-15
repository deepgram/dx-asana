package cmd

import (
	"fmt"
	"os"

	"github.com/deepgram/dx-asana/internal/asana"
	"github.com/deepgram/dx-asana/internal/ui"
	"github.com/spf13/cobra"
)

var discoverFields string

var workspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "List workspaces the authenticated user belongs to",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := asana.NewClient(token)

		opts := []asana.Option{}
		if discoverFields != "" {
			opts = append(opts, asana.WithOptFields(splitFields(discoverFields)...))
		} else {
			opts = append(opts, asana.WithOptFields("name", "is_organization", "email_domains"))
		}

		ws, err := client.GetWorkspaces(opts...)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		if jsonOutput {
			ui.PrintJSONWithMeta(ws, map[string]interface{}{"count": len(ws)}, nil)
		} else {
			fmt.Printf("🏢 Workspaces (%d)\n", len(ws))
			for _, w := range ws {
				kind := "workspace"
				if w.IsOrganization {
					kind = "organization"
				}
				fmt.Printf("  %s  %s  (%s)\n", w.GID, w.Name, kind)
			}
		}
		return nil
	},
}

var projectsCmd = &cobra.Command{
	Use:   "projects [workspace-id]",
	Short: "List projects in a workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceGID := args[0]
		client := asana.NewClient(token)

		opts := []asana.Option{}
		if discoverFields != "" {
			opts = append(opts, asana.WithOptFields(splitFields(discoverFields)...))
		} else {
			opts = append(opts, asana.WithOptFields(asana.DefaultProjectListFields...))
		}

		projects, err := client.GetProjects(workspaceGID, opts...)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		if jsonOutput {
			meta := map[string]interface{}{"count": len(projects), "workspace_gid": workspaceGID}
			ui.PrintJSONWithMeta(projects, meta, nil)
		} else {
			fmt.Printf("📂 Projects in workspace %s (%d)\n", workspaceGID, len(projects))
			for _, p := range projects {
				archived := ""
				if p.Archived {
					archived = " [archived]"
				}
				fmt.Printf("  %s  %s%s\n", p.GID, p.Name, archived)
			}
		}
		return nil
	},
}

var sectionsCmd = &cobra.Command{
	Use:   "sections [project-id]",
	Short: "List sections in a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectGID := args[0]
		client := asana.NewClient(token)

		sections, err := client.GetSections(projectGID, asana.WithOptFields("name", "created_at"))
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		if jsonOutput {
			meta := map[string]interface{}{"count": len(sections), "project_gid": projectGID}
			ui.PrintJSONWithMeta(sections, meta, nil)
		} else {
			fmt.Printf("🗂  Sections in project %s (%d)\n", projectGID, len(sections))
			for _, s := range sections {
				fmt.Printf("  %s  %s\n", s.GID, s.Name)
			}
		}
		return nil
	},
}

var subtasksCmd = &cobra.Command{
	Use:   "subtasks [parent-task-id]",
	Short: "List subtasks of a parent task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parentGID := args[0]
		client := asana.NewClient(token)

		opts := []asana.Option{}
		if discoverFields != "" {
			opts = append(opts, asana.WithOptFields(splitFields(discoverFields)...))
		} else {
			opts = append(opts, asana.WithOptFields(asana.DefaultTaskListFields...))
		}

		tasks, err := client.GetSubtasks(parentGID, opts...)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		if jsonOutput {
			meta := map[string]interface{}{"count": len(tasks), "parent_gid": parentGID}
			ui.PrintJSONWithMeta(tasks, meta, nil)
		} else {
			fmt.Printf("📑 Subtasks of %s (%d)\n", parentGID, len(tasks))
			for _, t := range tasks {
				marker := "○"
				if t.Completed {
					marker = "●"
				}
				fmt.Printf("  %s %s\n", marker, t.Name)
			}
		}
		return nil
	},
}

func init() {
	workspacesCmd.Flags().StringVar(&discoverFields, "fields", "", "Comma-separated opt_fields (overrides default)")
	projectsCmd.Flags().StringVar(&discoverFields, "fields", "", "Comma-separated opt_fields (overrides default)")
	subtasksCmd.Flags().StringVar(&discoverFields, "fields", "", "Comma-separated opt_fields (overrides default)")
}
