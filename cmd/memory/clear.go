package memory

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newClearCmd() *cobra.Command {
	var (
		force     bool
		sessionID string
	)

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Delete memories",
		Long:  "Delete all memories or memories from a specific session (WARNING: This action cannot be undone)",
		Example: `  # Clear all memories (with confirmation)
  aphelion memory clear

  # Clear memories for specific session
  aphelion memory clear --session abc123

  # Force clear without confirmation
  aphelion memory clear --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !config.IsAuthenticated() {
				return fmt.Errorf("authentication required. Please run 'aphelion auth login' first")
			}

			var warningMsg, spinnerMsg, successMsg string
			var endpoint string
			
			if sessionID != "" {
				warningMsg = fmt.Sprintf("This will delete all memories for session %s permanently!", sessionID)
				spinnerMsg = fmt.Sprintf("Deleting memories for session %s...", sessionID)
				successMsg = fmt.Sprintf("Memories for session %s have been deleted", sessionID)
				endpoint = fmt.Sprintf("/memory?session_id=%s", sessionID)
			} else {
				warningMsg = "This will delete ALL of your memories permanently!"
				spinnerMsg = "Deleting all memories..."
				successMsg = "All memories have been deleted"
				endpoint = "/memory"
			}

			if !force {
				utils.PrintWarning(warningMsg)
				fmt.Print("Are you sure you want to continue? (y/N): ")
				var response string
				if _, err := fmt.Scanln(&response); err != nil || (response != "y" && response != "Y") {
					utils.PrintInfo("Operation cancelled")
					return nil
				}
			}

			client := api.NewClient()
			
			spinner := utils.NewSpinner(spinnerMsg)
			spinner.Start()

			err := client.Delete(endpoint)
			spinner.Stop()

			if err != nil {
				return fmt.Errorf("failed to clear memories: %w", err)
			}

			utils.PrintSuccess(successMsg)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "skip confirmation prompt")
	cmd.Flags().StringVar(&sessionID, "session", "", "clear memories for specific session ID")

	return cmd
}