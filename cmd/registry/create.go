package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newCreateCmd() *cobra.Command {
	var name, description, specFile string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Register a new API service",
		Long:  "Register a new API service with OpenAPI specification",
		Example: `  # Register service with OpenAPI spec file
  aphelion registry create --name "My API" --description "My API service" --spec-file openapi.json

  # Register service with flags
  aphelion registry create --name "Weather API" --description "Weather service"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !config.IsAuthenticated() {
				return fmt.Errorf("authentication required. Please run 'aphelion auth login' first")
			}

			if name == "" {
				return fmt.Errorf("service name is required")
			}

			if description == "" {
				return fmt.Errorf("service description is required")
			}

			client := api.NewClient()
			
			createReq := map[string]interface{}{
				"name":        name,
				"description": description,
			}

			if specFile != "" {
				specData, err := os.ReadFile(specFile)
				if err != nil {
					return fmt.Errorf("failed to read spec file: %w", err)
				}

				var spec map[string]interface{}
				if err := json.Unmarshal(specData, &spec); err != nil {
					return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
				}

				createReq["spec"] = spec
			}

			spinner := utils.NewSpinner("Registering service...")
			spinner.Start()

			var response map[string]interface{}
			err := client.Post("/owner/services", createReq, &response)
			spinner.Stop()

			if err != nil {
				return fmt.Errorf("failed to register service: %w", err)
			}

			utils.PrintSuccess("Service registered successfully")
			
			if service, ok := response["service"].(map[string]interface{}); ok {
				return utils.PrintOutput(service, config.GetOutputFormat())
			}

			return utils.PrintOutput(response, config.GetOutputFormat())
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "service name (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "service description (required)")
	cmd.Flags().StringVarP(&specFile, "spec-file", "f", "", "OpenAPI specification file path")

	return cmd
}