package registry

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newGetCmd() *cobra.Command {
	var showManifest bool

	cmd := &cobra.Command{
		Use:   "get [SERVICE_ID]",
		Short: "Get service details",
		Long:  "Get detailed information about a specific service",
		Args:  cobra.ExactArgs(1),
		Example: `  # Get service details
  aphelion registry get service-123

  # Get service manifest for tool discovery
  aphelion registry get service-123 --manifest`,
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceID := args[0]
			client := api.NewClient()

			if showManifest {
				var manifest map[string]interface{}
				endpoint := fmt.Sprintf("/services/%s/manifest", serviceID)
				if err := client.Get(endpoint, &manifest); err != nil {
					return fmt.Errorf("failed to get service manifest: %w", err)
				}

				return utils.PrintOutput(manifest, config.GetOutputFormat())
			}

			var service api.Service
			endpoint := fmt.Sprintf("/services/%s", serviceID)
			if err := client.Get(endpoint, &service); err != nil {
				return fmt.Errorf("failed to get service: %w", err)
			}

			data := map[string]interface{}{
				"ID":          service.ID,
				"Name":        service.Name,
				"Description": service.Description,
				"Created":     service.CreatedAt.Format("2006-01-02 15:04:05"),
				"Updated":     service.UpdatedAt.Format("2006-01-02 15:04:05"),
			}

			if config.GetOutputFormat() == "json" || config.GetOutputFormat() == "yaml" {
				data["Spec"] = service.Spec
			}

			return utils.PrintOutput(data, config.GetOutputFormat())
		},
	}

	cmd.Flags().BoolVarP(&showManifest, "manifest", "m", false, "show service manifest instead of details")

	return cmd
}