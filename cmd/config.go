package cmd

import (
	"fmt"
	"os"

	"github.com/TheCoolRobot/asana-cli/internal/config"
	"github.com/TheCoolRobot/asana-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	setToken     string
	setWorkspace string
	setName      string
	setDesc      string
)

var (
	clearAll       bool
	clearToken     bool
	clearWorkspace bool
	clearCurrent   bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Get and set configuration values",
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Clear stored configuration values",
	Long:  "Clear stored config values. Use --all to wipe everything, or pick individual flags.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !clearAll && !clearToken && !clearWorkspace && !clearCurrent {
			return fmt.Errorf("specify at least one of --all, --token, --workspace, --current-project")
		}
		cfg, _ := config.Load()
		if clearAll {
			cfg.APIToken = ""
			cfg.DefaultWorkspace = ""
			cfg.CurrentProject = ""
			cfg.Projects = map[string]config.ProjectConfig{}
		} else {
			if clearToken {
				cfg.APIToken = ""
			}
			if clearWorkspace {
				cfg.DefaultWorkspace = ""
			}
			if clearCurrent {
				cfg.CurrentProject = ""
			}
		}
		if err := cfg.Save(); err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}
		if jsonOutput {
			ui.PrintJSON(cfg, nil)
		} else {
			fmt.Println("✓ Configuration cleared")
		}
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		if jsonOutput {
			ui.PrintJSON(cfg, nil)
		} else {
			fmt.Printf("Configuration:\n")
			fmt.Printf("  API Token: %s\n", maskToken(cfg.APIToken))
			fmt.Printf("  Default Workspace: %s\n", cfg.DefaultWorkspace)
			fmt.Printf("  Current Project: %s\n", cfg.CurrentProject)
			fmt.Printf("\n  Projects:\n")
			for name, proj := range cfg.Projects {
				marker := " "
				if name == cfg.CurrentProject {
					marker = "→"
				}
				fmt.Printf("    %s %s (ID: %s)\n", marker, name, proj.ProjectID)
				if proj.Description != "" {
					fmt.Printf("      %s\n", proj.Description)
				}
			}
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()

		if setToken != "" {
			cfg.APIToken = setToken
		}

		if setWorkspace != "" {
			cfg.DefaultWorkspace = setWorkspace
		}

		err := cfg.Save()
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		if jsonOutput {
			ui.PrintJSON(cfg, nil)
		} else {
			fmt.Println("✓ Configuration saved")
		}

		return nil
	},
}

var configProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  "Add, remove, list, or switch between saved projects",
}

var projectAddCmd = &cobra.Command{
	Use:   "add [name] [project-id]",
	Short: "Add a new project",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		projectID := args[1]

		cfg, _ := config.Load()
		err := cfg.AddProject(name, projectID, workspace, setDesc)

		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Printf("Error: %v\n", err)
			}
			return err
		}

		if jsonOutput {
			ui.PrintJSON(cfg.GetCurrentProject(), nil)
		} else {
			fmt.Printf("✓ Project added: %s\n", name)
			fmt.Printf("  ID: %s\n", projectID)
			fmt.Printf("  Set as current\n")
		}

		return nil
	},
}

var projectRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, _ := config.Load()
		err := cfg.RemoveProject(name)

		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Printf("Error: %v\n", err)
			}
			return err
		}

		if jsonOutput {
			meta := map[string]interface{}{"action": "removed", "project": name}
			ui.PrintJSONWithMeta(map[string]string{"status": "deleted"}, meta, nil)
		} else {
			fmt.Printf("✓ Project removed: %s\n", name)
		}

		return nil
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		projects := cfg.ListProjects()

		if jsonOutput {
			meta := map[string]interface{}{
				"count":           len(projects),
				"current_project": cfg.CurrentProject,
			}
			ui.PrintJSONWithMeta(projects, meta, nil)
		} else {
			if len(projects) == 0 {
				fmt.Println("No projects configured. Add one with: asana-cli config project add <name> <project-id>")
				return nil
			}

			fmt.Println("Projects:")
			for _, proj := range projects {
				marker := " "
				if proj.Name == cfg.CurrentProject {
					marker = "→"
				}
				fmt.Printf("  %s %s\n", marker, proj.Name)
				fmt.Printf("    Project ID: %s\n", proj.ProjectID)
				if proj.Description != "" {
					fmt.Printf("    Description: %s\n", proj.Description)
				}
			}
		}

		return nil
	},
}

var projectSwitchCmd = &cobra.Command{
	Use:   "switch [name]",
	Short: "Switch to a different project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, _ := config.Load()
		err := cfg.SetCurrentProject(name)

		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Printf("Error: %v\n", err)
			}
			return err
		}

		if jsonOutput {
			ui.PrintJSON(cfg.GetCurrentProject(), nil)
		} else {
			fmt.Printf("✓ Switched to project: %s\n", name)
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
	configCmd.AddCommand(configProjectCmd)

	configSetCmd.Flags().StringVar(&setToken, "token", "", "API token")
	configSetCmd.Flags().StringVar(&setWorkspace, "workspace", "", "Default workspace ID")
	configSetCmd.Flags().StringVar(&setName, "name", "", "Default name")

	configUnsetCmd.Flags().BoolVar(&clearAll, "all", false, "Wipe the entire config (token, workspace, projects)")
	configUnsetCmd.Flags().BoolVar(&clearToken, "token", false, "Clear stored API token")
	configUnsetCmd.Flags().BoolVar(&clearWorkspace, "workspace", false, "Clear default workspace")
	configUnsetCmd.Flags().BoolVar(&clearCurrent, "current-project", false, "Clear current project selection (does not delete saved projects)")

	configProjectCmd.AddCommand(projectAddCmd)
	configProjectCmd.AddCommand(projectRemoveCmd)
	configProjectCmd.AddCommand(projectListCmd)
	configProjectCmd.AddCommand(projectSwitchCmd)

	projectAddCmd.Flags().StringVar(&workspace, "workspace", "", "Workspace ID for this project")
	projectAddCmd.Flags().StringVar(&setDesc, "description", "", "Project description")
}

func maskToken(token string) string {
	if token == "" {
		return "(not set)"
	}
	if len(token) < 8 {
		return "***" + token[len(token)-1:]
	}
	return token[:4] + "..." + token[len(token)-4:]
}
