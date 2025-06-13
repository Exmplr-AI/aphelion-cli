package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Service represents an API service
type Service struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	OwnerID     string                 `json:"owner_id"`
	Spec        map[string]interface{} `json:"spec"`
	Pricing     map[string]interface{} `json:"pricing,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Manage API services",
	Long:  `Register and manage API services in the Aphelion Gateway registry.`,
}

// servicesListCmd represents the services list command
var servicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API services",
	Long:  `List all publicly available API services or your own services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		mine, _ := cmd.Flags().GetBool("mine")
		
		var endpoint string
		if mine {
			endpoint = "/owner/services"
		} else {
			endpoint = "/services"
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var services []Service
		if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(services)
		case "yaml":
			return outputYAML(services)
		default:
			if len(services) == 0 {
				if mine {
					fmt.Println("No services found. Register a service with 'aphelion services register'.")
				} else {
					fmt.Println("No public services available.")
				}
				return nil
			}
			
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Description", "Created"})
			table.SetBorder(false)
			table.SetColWidth(50)
			table.SetRowLine(true)
			
			for _, service := range services {
				description := service.Description
				if len(description) > 47 {
					description = description[:47] + "..."
				}
				
				table.Append([]string{
					service.ID,
					service.Name,
					description,
					service.CreatedAt,
				})
			}
			
			table.Render()
		}
		
		return nil
	},
}

// servicesDescribeCmd represents the services describe command
var servicesDescribeCmd = &cobra.Command{
	Use:   "describe SERVICE_ID",
	Short: "Get detailed information about a service",
	Long:  `Get detailed information about a specific API service.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceID := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, fmt.Sprintf("/owner/services/%s", serviceID))
		if err != nil {
			return fmt.Errorf("failed to get service: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var service Service
		if err := json.NewDecoder(resp.Body).Decode(&service); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(service)
		case "yaml":
			return outputYAML(service)
		default:
			fmt.Printf("ID: %s\n", service.ID)
			fmt.Printf("Name: %s\n", service.Name)
			fmt.Printf("Description: %s\n", service.Description)
			fmt.Printf("Owner ID: %s\n", service.OwnerID)
			fmt.Printf("Created: %s\n", service.CreatedAt)
			fmt.Printf("Updated: %s\n", service.UpdatedAt)
			
			if service.Pricing != nil {
				fmt.Println("Pricing:")
				pricingJSON, _ := json.MarshalIndent(service.Pricing, "  ", "  ")
				fmt.Printf("  %s\n", string(pricingJSON))
			}
			
			fmt.Println("OpenAPI Specification:")
			specJSON, _ := json.MarshalIndent(service.Spec, "  ", "  ")
			fmt.Printf("  %s\n", string(specJSON))
		}
		
		return nil
	},
}

// servicesRegisterCmd represents the services register command
var servicesRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new API service",
	Long: `Register a new API service with OpenAPI specification.

Examples:
  aphelion services register --spec=api.yaml
  aphelion services register --spec=api.json --name="My API"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		specFile, _ := cmd.Flags().GetString("spec")
		if specFile == "" {
			return fmt.Errorf("--spec flag is required")
		}
		
		// Read spec file
		specData, err := os.ReadFile(specFile)
		if err != nil {
			return fmt.Errorf("failed to read spec file: %w", err)
		}
		
		var spec map[string]interface{}
		if err := yaml.Unmarshal(specData, &spec); err != nil {
			// Try JSON if YAML fails
			if err := json.Unmarshal(specData, &spec); err != nil {
				return fmt.Errorf("failed to parse spec file as YAML or JSON: %w", err)
			}
		}
		
		// Extract name and description from spec or flags
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		
		if name == "" {
			if info, ok := spec["info"].(map[string]interface{}); ok {
				if title, ok := info["title"].(string); ok {
					name = title
				}
			}
		}
		
		if description == "" {
			if info, ok := spec["info"].(map[string]interface{}); ok {
				if desc, ok := info["description"].(string); ok {
					description = desc
				}
			}
		}
		
		if name == "" {
			return fmt.Errorf("service name is required (use --name or set info.title in spec)")
		}
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		payload := map[string]interface{}{
			"name":        name,
			"description": description,
			"spec":        spec,
		}
		
		// Add pricing if provided
		if pricingFile, _ := cmd.Flags().GetString("pricing"); pricingFile != "" {
			pricingData, err := os.ReadFile(pricingFile)
			if err != nil {
				return fmt.Errorf("failed to read pricing file: %w", err)
			}
			
			var pricing map[string]interface{}
			if err := json.Unmarshal(pricingData, &pricing); err != nil {
				return fmt.Errorf("failed to parse pricing file: %w", err)
			}
			
			payload["pricing"] = pricing
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		
		resp, err := client.Post(ctx, "/owner/services", payload)
		if err != nil {
			return fmt.Errorf("failed to register service: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(result)
		case "yaml":
			return outputYAML(result)
		default:
			if service, ok := result["service"].(map[string]interface{}); ok {
				if id, ok := service["id"].(string); ok {
					fmt.Printf("✓ Registered service: %s (%s)\n", name, id)
				}
			} else {
				fmt.Printf("✓ Service '%s' registered successfully\n", name)
			}
		}
		
		return nil
	},
}

// servicesManifestCmd represents the services manifest command
var servicesManifestCmd = &cobra.Command{
	Use:   "manifest SERVICE_ID",
	Short: "Get service manifest for tool discovery",
	Long:  `Get the service manifest that contains tool information for agent discovery.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceID := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, fmt.Sprintf("/services/%s/manifest", serviceID))
		if err != nil {
			return fmt.Errorf("failed to get service manifest: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var manifest map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(manifest)
		case "yaml":
			return outputYAML(manifest)
		default:
			manifestJSON, _ := json.MarshalIndent(manifest, "", "  ")
			fmt.Println(string(manifestJSON))
		}
		
		return nil
	},
}

// servicesDeleteCmd represents the services delete command
var servicesDeleteCmd = &cobra.Command{
	Use:   "delete SERVICE_ID",
	Short: "Delete a service",
	Long:  `Delete a service that you own. This cannot be undone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceID := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Delete(ctx, fmt.Sprintf("/owner/services/%s", serviceID))
		if err != nil {
			return fmt.Errorf("failed to delete service: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		fmt.Printf("✓ Deleted service: %s\n", serviceID)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)
	
	servicesCmd.AddCommand(servicesListCmd)
	servicesCmd.AddCommand(servicesDescribeCmd)
	servicesCmd.AddCommand(servicesRegisterCmd)
	servicesCmd.AddCommand(servicesManifestCmd)
	servicesCmd.AddCommand(servicesDeleteCmd)
	
	// Flags
	servicesListCmd.Flags().Bool("mine", false, "Show only your own services")
	
	servicesRegisterCmd.Flags().String("spec", "", "Path to OpenAPI specification file (required)")
	servicesRegisterCmd.Flags().String("name", "", "Service name (overrides spec)")
	servicesRegisterCmd.Flags().String("description", "", "Service description (overrides spec)")
	servicesRegisterCmd.Flags().String("pricing", "", "Path to pricing configuration JSON file")
	servicesRegisterCmd.MarkFlagRequired("spec")
}