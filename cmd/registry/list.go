package registry

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List public services",
		Long:  "List all publicly available services in the registry",
		Example: `  # List all public services
  aphelion registry list

  # List services in JSON format
  aphelion registry list --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient()
			
			var response api.ServicesResponse
			if err := client.Get("/services", &response); err != nil {
				return fmt.Errorf("failed to list services: %w", err)
			}

			if len(response.Services) == 0 {
				utils.PrintInfo("No public services found")
				return nil
			}

			utils.PrintInfo("Found %d public services", len(response.Services))
			
			var data []map[string]interface{}
			for _, service := range response.Services {
				data = append(data, map[string]interface{}{
					"ID":          service.ID,
					"Name":        service.Name,
					"Description": service.Description,
					"Created":     service.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}

			return utils.PrintOutput(data, config.GetOutputFormat())
		},
	}

	return cmd
}