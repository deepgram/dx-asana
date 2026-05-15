package cmd

import (
	"fmt"
	"os"

	"github.com/deepgram/dx-asana/internal/asana"
	"github.com/deepgram/dx-asana/internal/ui"
	"github.com/spf13/cobra"
)

var (
	commentText     string
	commentHTMLText string
	commentPinned   bool
	commentList     bool
)

var commentCmd = &cobra.Command{
	Use:   "comment [task-id]",
	Short: "Add a comment to a task, or list existing comments",
	Long: `Add a comment to a task or list its existing comments.

Examples:
  dx-asana comment <task-gid> --text "Status update: blocked on review"
  dx-asana comment <task-gid> --list           # show all comments
  dx-asana comment <task-gid> --html-text "<body><b>Bold</b> note</body>"
  dx-asana comment <task-gid> --text "..." --pin`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskGID := args[0]
		client := asana.NewClient(token)

		if commentList {
			stories, err := client.GetStories(taskGID, asana.WithOptFields(
				"created_at", "created_by.name", "text", "type", "resource_subtype", "is_pinned"))
			if err != nil {
				if jsonOutput {
					ui.PrintJSON(nil, err)
				} else {
					fmt.Fprintln(os.Stderr, "Error:", err)
				}
				return err
			}
			if jsonOutput {
				ui.PrintJSONWithMeta(stories, map[string]interface{}{"count": len(stories)}, nil)
			} else {
				fmt.Printf("💬 %d stories on task %s\n", len(stories), taskGID)
				for _, s := range stories {
					ts := ""
					if s.CreatedAt != nil {
						ts = s.CreatedAt.Format("2006-01-02 15:04")
					}
					author := ""
					if s.CreatedBy != nil {
						author = s.CreatedBy.Name
					}
					prefix := "·"
					if s.Type == "comment" {
						prefix = "💬"
					}
					fmt.Printf("  %s [%s] %s: %s\n", prefix, ts, author, s.Text)
				}
			}
			return nil
		}

		if commentText == "" && commentHTMLText == "" {
			err := fmt.Errorf("provide --text or --html-text (or use --list to show existing comments)")
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		req := &asana.StoryCreateRequest{
			Text:     commentText,
			HTMLText: commentHTMLText,
			IsPinned: commentPinned,
		}

		story, err := client.CreateStory(taskGID, req)
		if err != nil {
			if jsonOutput {
				ui.PrintJSON(nil, err)
			} else {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			return err
		}

		if jsonOutput {
			meta := map[string]interface{}{"action": "commented", "task_gid": taskGID}
			ui.PrintJSONWithMeta(story, meta, nil)
		} else {
			fmt.Printf("✓ Comment posted (story %s)\n", story.GID)
		}
		return nil
	},
}

func init() {
	commentCmd.Flags().StringVar(&commentText, "text", "", "Plain-text comment body")
	commentCmd.Flags().StringVar(&commentHTMLText, "html-text", "", "HTML-formatted comment body")
	commentCmd.Flags().BoolVar(&commentPinned, "pin", false, "Pin the comment")
	commentCmd.Flags().BoolVar(&commentList, "list", false, "List existing stories on the task instead of posting")
}
