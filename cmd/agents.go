package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Agent represents an agent session
type Agent struct {
	SessionID          string   `json:"session_id"`
	UserID             string   `json:"user_id"`
	SubscribedServices []string `json:"subscribed_services"`
	CreatedAt          string   `json:"created_at"`
	LastActivity       string   `json:"last_activity,omitempty"`
}

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Manage agent sessions",
	Long:  `Create and manage agent sessions for tool execution and API interaction.`,
}

// agentsCreateCmd represents the agents create command
var agentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new agent session",
	Long: `Create a new agent session with optional service subscriptions.

Examples:
  aphelion agents create
  aphelion agents create --services=service1,service2`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		services, _ := cmd.Flags().GetString("services")
		var serviceList []string
		if services != "" {
			serviceList = strings.Split(services, ",")
		}
		
		payload := map[string]interface{}{
			"subscribed_services": serviceList,
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Post(ctx, "/v1/agents", payload)
		if err != nil {
			return fmt.Errorf("failed to create agent session: %w", err)
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
			if sessionID, ok := result["session_id"].(string); ok {
				fmt.Printf("✓ Created agent session: %s\n", sessionID)
				if len(serviceList) > 0 {
					fmt.Printf("Subscribed services: %s\n", strings.Join(serviceList, ", "))
				}
			} else {
				fmt.Println("✓ Agent session created successfully")
			}
		}
		
		return nil
	},
}

// agentsListCmd represents the agents list command
var agentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List agent sessions",
	Long:  `List all agent sessions for the current user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, "/v1/agents")
		if err != nil {
			return fmt.Errorf("failed to list agent sessions: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var agents []Agent
		if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(agents)
		case "yaml":
			return outputYAML(agents)
		default:
			if len(agents) == 0 {
				fmt.Println("No agent sessions found.")
				return nil
			}
			
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Session ID", "Services", "Created", "Last Activity"})
			table.SetBorder(false)
			
			for _, agent := range agents {
				services := strings.Join(agent.SubscribedServices, ", ")
				if services == "" {
					services = "none"
				}
				
				lastActivity := agent.LastActivity
				if lastActivity == "" {
					lastActivity = "never"
				}
				
				table.Append([]string{
					agent.SessionID,
					services,
					agent.CreatedAt,
					lastActivity,
				})
			}
			
			table.Render()
		}
		
		return nil
	},
}

// agentsDescribeCmd represents the agents describe command
var agentsDescribeCmd = &cobra.Command{
	Use:   "describe SESSION_ID",
	Short: "Get detailed information about an agent session",
	Long:  `Get detailed information about a specific agent session.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, fmt.Sprintf("/v1/agents/%s", sessionID))
		if err != nil {
			return fmt.Errorf("failed to get agent session: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var agent Agent
		if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(agent)
		case "yaml":
			return outputYAML(agent)
		default:
			fmt.Printf("Session ID: %s\n", agent.SessionID)
			fmt.Printf("User ID: %s\n", agent.UserID)
			fmt.Printf("Created: %s\n", agent.CreatedAt)
			
			if agent.LastActivity != "" {
				fmt.Printf("Last Activity: %s\n", agent.LastActivity)
			}
			
			if len(agent.SubscribedServices) > 0 {
				fmt.Printf("Subscribed Services:\n")
				for _, service := range agent.SubscribedServices {
					fmt.Printf("  - %s\n", service)
				}
			} else {
				fmt.Println("Subscribed Services: none")
			}
		}
		
		return nil
	},
}

// agentsDeleteCmd represents the agents delete command
var agentsDeleteCmd = &cobra.Command{
	Use:   "delete SESSION_ID",
	Short: "Delete an agent session",
	Long:  `Delete a specific agent session. This cannot be undone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Delete(ctx, fmt.Sprintf("/v1/agents/%s", sessionID))
		if err != nil {
			return fmt.Errorf("failed to delete agent session: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		fmt.Printf("✓ Deleted agent session: %s\n", sessionID)
		
		return nil
	},
}

// agentsExecuteCmd represents the agents execute command
var agentsExecuteCmd = &cobra.Command{
	Use:   "execute SESSION_ID --tool=TOOL_NAME [--params=JSON]",
	Short: "Execute a tool in an agent session",
	Long: `Execute a tool within a specific agent session.

Examples:
  aphelion agents execute abc123 --tool=echo --params='{"message":"hello"}'
  aphelion agents execute abc123 --tool=get_current_time
  aphelion agents execute abc123 --tool=calculate --params='{"expression":"2+2"}'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		
		toolName, _ := cmd.Flags().GetString("tool")
		if toolName == "" {
			return fmt.Errorf("--tool flag is required")
		}
		
		paramsStr, _ := cmd.Flags().GetString("params")
		var params interface{}
		if paramsStr != "" {
			if err := json.Unmarshal([]byte(paramsStr), &params); err != nil {
				return fmt.Errorf("invalid JSON in --params: %w", err)
			}
		} else {
			params = map[string]interface{}{}
		}
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		payload := map[string]interface{}{
			"tool":       toolName,
			"parameters": params,
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		
		resp, err := client.Post(ctx, fmt.Sprintf("/v1/agents/%s/execute", sessionID), payload)
		if err != nil {
			return fmt.Errorf("failed to execute tool: %w", err)
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
			fmt.Printf("Tool: %s\n", toolName)
			if success, ok := result["success"].(bool); ok && success {
				fmt.Println("Status: ✓ Success")
			} else {
				fmt.Println("Status: ✗ Failed")
			}
			
			if resultData, ok := result["result"]; ok {
				fmt.Println("Result:")
				resultJSON, _ := json.MarshalIndent(resultData, "  ", "  ")
				fmt.Printf("  %s\n", string(resultJSON))
			}
			
			if errorMsg, ok := result["error"].(string); ok && errorMsg != "" {
				fmt.Printf("Error: %s\n", errorMsg)
			}
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentsCmd)
	
	agentsCmd.AddCommand(agentsCreateCmd)
	agentsCmd.AddCommand(agentsListCmd)
	agentsCmd.AddCommand(agentsDescribeCmd)
	agentsCmd.AddCommand(agentsDeleteCmd)
	agentsCmd.AddCommand(agentsExecuteCmd)
	
	// Flags
	agentsCreateCmd.Flags().String("services", "", "Comma-separated list of service IDs to subscribe to")
	agentsExecuteCmd.Flags().String("tool", "", "Tool name to execute (required)")
	agentsExecuteCmd.Flags().String("params", "", "Tool parameters as JSON")
	agentsExecuteCmd.MarkFlagRequired("tool")
}