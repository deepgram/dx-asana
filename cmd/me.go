package cmd

import (
	"fmt"
	"os"

	"github.com/deepgram/dx-asana/internal/asana"
	"github.com/deepgram/dx-asana/internal/ui"
	"github.com/spf13/cobra"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := asana.NewClient(token)

		user, err := client.GetMe()
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		if jsonOutput {
			ui.PrintJSON(user, nil)
		} else {
			fmt.Printf("👤 %s\n", user.Name)
			fmt.Printf("   Email: %s\n", user.Email)
			fmt.Printf("   GID: %s\n", user.GID)
		}

		return nil
	},
}
