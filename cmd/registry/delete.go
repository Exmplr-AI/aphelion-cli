package registry

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete [SERVICE_ID]",
		Short: "Delete a service",
		Long:  "Delete a service that you own",
		Args:  cobra.ExactArgs(1),
		Example: `  # Delete a service (with confirmation)
  aphelion registry delete service-123

  # Force delete without confirmation
  aphelion registry delete service-123 --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !config.IsAuthenticated() {
				return fmt.Errorf("authentication required. Please run 'aphelion auth login' first")
			}

			serviceID := args[0]

			if !force {
				fmt.Printf("Are you sure you want to delete service %s? (y/N): ", serviceID)
				var response string
				if _, err := fmt.Scanln(&response); err != nil || (response != "y" && response != "Y") {
					utils.PrintInfo("Operation cancelled")
					return nil
				}
			}

			client := api.NewClient()
			
			spinner := utils.NewSpinner("Deleting service...")
			spinner.Start()

			endpoint := fmt.Sprintf("/owner/services/%s", serviceID)
			err := client.Delete(endpoint)
			spinner.Stop()

			if err != nil {
				return fmt.Errorf("failed to delete service: %w", err)
			}

			utils.PrintSuccess("Service deleted successfully")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "skip confirmation prompt")

	return cmd
}