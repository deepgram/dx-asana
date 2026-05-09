package cmd

import (
	"github.com/TheCoolRobot/asana-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	token      string
	workspace  string
	project    string
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "asana-cli",
	Short:   "Asana CLI - Beautiful task management",
	Long:    "A feature-rich CLI for managing Asana tasks with TUI and sync daemon",
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if token == "" {
			token = config.GetAPIToken()
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Asana API token (or set ASANA_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&workspace, "workspace", "", "Default workspace ID")
	rootCmd.PersistentFlags().StringVar(&project, "project", "", "Default project ID")

	// Add all commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(meCmd)
}

func Execute() error {
	rootCmd.Version = Version
	return rootCmd.Execute()
}