package registry

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newMyServicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "my-services",
		Short: "List your registered services",
		Long:  "List all services that you have registered",
		Example: `  # List your services
  aphelion registry my-services

  # List your services in JSON format
  aphelion registry my-services --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !config.IsAuthenticated() {
				return fmt.Errorf("authentication required. Please run 'aphelion auth login' first")
			}

			client := api.NewClient()
			
			var response api.ServicesResponse
			if err := client.Get("/owner/services", &response); err != nil {
				return fmt.Errorf("failed to list your services: %w", err)
			}

			if len(response.Services) == 0 {
				utils.PrintInfo("No services found")
				return nil
			}

			utils.PrintInfo("Found %d services", len(response.Services))
			
			var data []map[string]interface{}
			for _, service := range response.Services {
				data = append(data, map[string]interface{}{
					"ID":          service.ID,
					"Name":        service.Name,
					"Description": service.Description,
					"Created":     service.CreatedAt.Format("2006-01-02 15:04:05"),
					"Updated":     service.UpdatedAt.Format("2006-01-02 15:04:05"),
				})
			}

			return utils.PrintOutput(data, config.GetOutputFormat())
		},
	}

	return cmd
}