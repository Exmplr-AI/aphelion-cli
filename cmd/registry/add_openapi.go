package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
)

type OpenAPISpec struct {
	OpenAPI string                 `json:"openapi"`
	Info    OpenAPIInfo            `json:"info"`
	Paths   map[string]interface{} `json:"paths"`
	Servers []OpenAPIServer        `json:"servers,omitempty"`
}

type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type OpenAPIServer struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type STELLAManifest struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Version     string                   `json:"version"`
	Tools       []STELLATool             `json:"tools"`
	Metadata    map[string]interface{}   `json:"metadata"`
}

type STELLATool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Method      string                 `json:"method"`
	Endpoint    string                 `json:"endpoint"`
	Parameters  map[string]interface{} `json:"parameters"`
}

func newAddOpenAPICmd() *cobra.Command {
	var (
		file     string
		name     string
		desc     string
		baseURL  string
	)
	
	cmd := &cobra.Command{
		Use:   "add-openapi",
		Short: "Register service from OpenAPI specification",
		Long:  "Parse OpenAPI specification and generate STELLA manifest for service registration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !config.IsAuthenticated() {
				return fmt.Errorf("authentication required. Please run 'aphelion auth login' first")
			}

			if file == "" {
				return fmt.Errorf("OpenAPI file is required")
			}

			return processOpenAPIFile(file, name, desc, baseURL)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "OpenAPI specification file path (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Override service name")
	cmd.Flags().StringVarP(&desc, "description", "d", "", "Override service description")
	cmd.Flags().StringVar(&baseURL, "base-url", "", "Override base URL")
	
	cmd.MarkFlagRequired("file")

	return cmd
}

func processOpenAPIFile(file, name, desc, baseURL string) error {
	// Read OpenAPI file
	specData, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI file: %w", err)
	}

	var openAPISpec OpenAPISpec
	if err := json.Unmarshal(specData, &openAPISpec); err != nil {
		return fmt.Errorf("failed to parse OpenAPI specification: %w", err)
	}

	// Generate STELLA manifest
	stella := generateSTELLAManifest(openAPISpec)
	
	// Override fields if provided
	if name != "" {
		stella.Name = name
	}
	if desc != "" {
		stella.Description = desc
	}
	if baseURL != "" {
		stella.Metadata["base_url"] = baseURL
	} else if len(openAPISpec.Servers) > 0 {
		stella.Metadata["base_url"] = openAPISpec.Servers[0].URL
	}

	if viper.GetBool("verbose") {
		fmt.Printf("Generated STELLA manifest:\n")
		utils.OutputJSON(stella)
		fmt.Println()
	}

	// Register service with STELLA manifest
	return registerServiceWithSTELLA(stella)
}

func generateSTELLAManifest(spec OpenAPISpec) STELLAManifest {
	stella := STELLAManifest{
		Name:        spec.Info.Title,
		Description: spec.Info.Description,
		Version:     spec.Info.Version,
		Tools:       []STELLATool{},
		Metadata: map[string]interface{}{
			"openapi_version": spec.OpenAPI,
			"generated_by":    "aphelion-cli",
		},
	}

	// Convert OpenAPI paths to STELLA tools
	for path, pathItem := range spec.Paths {
		if pathMap, ok := pathItem.(map[string]interface{}); ok {
			for method, operation := range pathMap {
				if opMap, ok := operation.(map[string]interface{}); ok {
					tool := convertOperationToTool(method, path, opMap)
					if tool != nil {
						stella.Tools = append(stella.Tools, *tool)
					}
				}
			}
		}
	}

	return stella
}

func convertOperationToTool(method, path string, operation map[string]interface{}) *STELLATool {
	// Skip non-HTTP methods
	validMethods := map[string]bool{
		"get": true, "post": true, "put": true, "delete": true, 
		"patch": true, "head": true, "options": true,
	}
	
	if !validMethods[method] {
		return nil
	}

	tool := &STELLATool{
		Method:     method,
		Endpoint:   path,
		Parameters: make(map[string]interface{}),
	}

	// Extract operation details
	if summary, ok := operation["summary"].(string); ok {
		tool.Name = summary
	} else {
		tool.Name = fmt.Sprintf("%s %s", method, path)
	}

	if description, ok := operation["description"].(string); ok {
		tool.Description = description
	} else {
		tool.Description = tool.Name
	}

	// Extract parameters
	if params, ok := operation["parameters"].([]interface{}); ok {
		properties := make(map[string]interface{})
		required := []string{}

		for _, param := range params {
			if paramMap, ok := param.(map[string]interface{}); ok {
				if name, ok := paramMap["name"].(string); ok {
					paramSchema := map[string]interface{}{
						"type": "string", // default type
					}
					
					if desc, ok := paramMap["description"].(string); ok {
						paramSchema["description"] = desc
					}
					
					if schema, ok := paramMap["schema"].(map[string]interface{}); ok {
						if paramType, ok := schema["type"].(string); ok {
							paramSchema["type"] = paramType
						}
					}
					
					if req, ok := paramMap["required"].(bool); ok && req {
						required = append(required, name)
					}
					
					properties[name] = paramSchema
				}
			}
		}

		if len(properties) > 0 {
			tool.Parameters = map[string]interface{}{
				"type":       "object",
				"properties": properties,
			}
			
			if len(required) > 0 {
				tool.Parameters["required"] = required
			}
		}
	}

	return tool
}

func registerServiceWithSTELLA(stella STELLAManifest) error {
	client := api.NewClient()
	
	spinner := utils.NewSpinner("Registering service with STELLA manifest...")
	spinner.Start()

	var response map[string]interface{}
	err := client.Post("/owner/services", stella, &response)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	utils.PrintSuccess("Service registered successfully from OpenAPI specification")
	
	if service, ok := response["service"].(map[string]interface{}); ok {
		return utils.OutputTable([]map[string]interface{}{service})
	}

	return utils.OutputJSON(response)
}