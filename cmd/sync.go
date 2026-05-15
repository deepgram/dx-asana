package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/TheCoolRobot/asana-cli/internal/syncdaemon"
	"github.com/spf13/cobra"
)

var projects string

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Start the sync daemon",
	Long:  "Start the background sync daemon to cache Asana data locally",
	RunE: func(cmd *cobra.Command, args []string) error {
		if projects == "" {
			return fmt.Errorf("--projects flag is required")
		}

		projectIDs := strings.Split(projects, ",")
		for i := range projectIDs {
			projectIDs[i] = strings.TrimSpace(projectIDs[i])
		}

		daemon := syncdaemon.NewDaemon(token, projectIDs)

		// Handle Ctrl+C gracefully
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		go func() {
			<-sigChan
			fmt.Println("\n[sync-daemon] Received interrupt signal")
			daemon.Stop()
		}()

		daemon.Start()
		return nil
	},
}

func init() {
	syncCmd.Flags().StringVar(&projects, "projects", "", "Comma-separated list of project IDs to sync (required)")
	if err := syncCmd.MarkFlagRequired("projects"); err != nil {
		log.Fatal(err)
	}
}