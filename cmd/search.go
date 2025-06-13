package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/exmplrai/aphelion-cli/pkg/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Tool represents a tool in search results
type Tool struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	ServiceID   string                 `json:"service_id"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	UsageCount  int                    `json:"usage_count,omitempty"`
	Similarity  float64                `json:"similarity,omitempty"`
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search and discover tools",
	Long:  `Search for tools and services using various discovery mechanisms.`,
}

// searchUniversalCmd represents the search universal command
var searchUniversalCmd = &cobra.Command{
	Use:   "universal QUERY",
	Short: "Search across tools and memories",
	Long:  `Search across tools and memories simultaneously.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		limit, _ := cmd.Flags().GetInt("limit")
		
		endpoint := fmt.Sprintf("/search?q=%s", query)
		if limit > 0 {
			endpoint += fmt.Sprintf("&limit=%d", limit)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to search: %w", err)
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
			// Display tools
			if toolsData, ok := result["tools"].([]interface{}); ok && len(toolsData) > 0 {
				fmt.Printf("Tools matching '%s':\n\n", query)
				
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Service", "Description"})
				table.SetBorder(false)
				table.SetColWidth(40)
				
				for _, toolData := range toolsData {
					if tool, ok := toolData.(map[string]interface{}); ok {
						name := getString(tool, "name")
						serviceID := getString(tool, "service_id")
						description := getString(tool, "description")
						
						if len(description) > 37 {
							description = description[:37] + "..."
						}
						
						table.Append([]string{name, serviceID, description})
					}
				}
				
				table.Render()
				fmt.Println()
			}
			
			// Display memories
			if memoriesData, ok := result["memories"].([]interface{}); ok && len(memoriesData) > 0 {
				fmt.Printf("Memories matching '%s':\n\n", query)
				
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"ID", "Session", "Summary"})
				table.SetBorder(false)
				table.SetColWidth(40)
				
				for _, memoryData := range memoriesData {
					if memory, ok := memoryData.(map[string]interface{}); ok {
						id := getString(memory, "id")
						sessionID := getString(memory, "session_id")
						if sessionID == "" {
							sessionID = "global"
						}
						summary := getString(memory, "summary")
						
						if len(summary) > 37 {
							summary = summary[:37] + "..."
						}
						
						table.Append([]string{id, sessionID, summary})
					}
				}
				
				table.Render()
			}
			
			if len(result) == 0 {
				fmt.Printf("No results found for query: %s\n", query)
			}
		}
		
		return nil
	},
}

// searchToolsCmd represents the search tools command
var searchToolsCmd = &cobra.Command{
	Use:   "tools QUERY",
	Short: "Search for tools",
	Long:  `Search for tools using semantic similarity.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		limit, _ := cmd.Flags().GetInt("limit")
		
		endpoint := fmt.Sprintf("/search/tools?q=%s", query)
		if limit > 0 {
			endpoint += fmt.Sprintf("&limit=%d", limit)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to search tools: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var tools []Tool
		if err := json.NewDecoder(resp.Body).Decode(&tools); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(tools)
		case "yaml":
			return outputYAML(tools)
		default:
			if len(tools) == 0 {
				fmt.Printf("No tools found matching: %s\n", query)
				return nil
			}
			
			fmt.Printf("Found %d tools matching: %s\n\n", len(tools), query)
			
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Service", "Description", "Similarity"})
			table.SetBorder(false)
			table.SetColWidth(35)
			
			for _, tool := range tools {
				description := tool.Description
				if len(description) > 32 {
					description = description[:32] + "..."
				}
				
				similarity := ""
				if tool.Similarity > 0 {
					similarity = fmt.Sprintf("%.2f", tool.Similarity)
				}
				
				table.Append([]string{
					tool.Name,
					tool.ServiceID,
					description,
					similarity,
				})
			}
			
			table.Render()
		}
		
		return nil
	},
}

// searchPopularCmd represents the search popular command
var searchPopularCmd = &cobra.Command{
	Use:   "popular",
	Short: "Get popular tools",
	Long:  `Get most popular tools based on usage frequency.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		limit, _ := cmd.Flags().GetInt("limit")
		
		endpoint := "/search/popular"
		if limit > 0 {
			endpoint += fmt.Sprintf("?limit=%d", limit)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to get popular tools: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var tools []Tool
		if err := json.NewDecoder(resp.Body).Decode(&tools); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(tools)
		case "yaml":
			return outputYAML(tools)
		default:
			if len(tools) == 0 {
				fmt.Println("No popular tools found.")
				return nil
			}
			
			fmt.Printf("Top %d popular tools:\n\n", len(tools))
			
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Rank", "Name", "Service", "Description", "Usage"})
			table.SetBorder(false)
			table.SetColWidth(30)
			
			for i, tool := range tools {
				description := tool.Description
				if len(description) > 27 {
					description = description[:27] + "..."
				}
				
				usage := ""
				if tool.UsageCount > 0 {
					usage = fmt.Sprintf("%d", tool.UsageCount)
				}
				
				table.Append([]string{
					fmt.Sprintf("%d", i+1),
					tool.Name,
					tool.ServiceID,
					description,
					usage,
				})
			}
			
			table.Render()
		}
		
		return nil
	},
}

// searchRecommendationsCmd represents the search recommendations command
var searchRecommendationsCmd = &cobra.Command{
	Use:   "recommendations",
	Short: "Get personalized tool recommendations",
	Long:  `Get personalized tool recommendations based on your usage history.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		limit, _ := cmd.Flags().GetInt("limit")
		
		endpoint := "/search/recommendations"
		if limit > 0 {
			endpoint += fmt.Sprintf("?limit=%d", limit)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to get recommendations: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var tools []Tool
		if err := json.NewDecoder(resp.Body).Decode(&tools); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(tools)
		case "yaml":
			return outputYAML(tools)
		default:
			if len(tools) == 0 {
				fmt.Println("No recommendations available. Use more tools to get personalized recommendations.")
				return nil
			}
			
			fmt.Printf("Recommended tools for you:\n\n")
			
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Service", "Description"})
			table.SetBorder(false)
			table.SetColWidth(40)
			
			for _, tool := range tools {
				description := tool.Description
				if len(description) > 37 {
					description = description[:37] + "..."
				}
				
				table.Append([]string{
					tool.Name,
					tool.ServiceID,
					description,
				})
			}
			
			table.Render()
		}
		
		return nil
	},
}

// searchSuggestCmd represents the search suggest command
var searchSuggestCmd = &cobra.Command{
	Use:   "suggest PARTIAL_QUERY",
	Short: "Get search suggestions",
	Long:  `Get search suggestions for partial queries.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		partial := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		limit, _ := cmd.Flags().GetInt("limit")
		
		endpoint := fmt.Sprintf("/search/suggestions?partial=%s", partial)
		if limit > 0 {
			endpoint += fmt.Sprintf("&limit=%d", limit)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to get suggestions: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var suggestions []string
		if err := json.NewDecoder(resp.Body).Decode(&suggestions); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(suggestions)
		case "yaml":
			return outputYAML(suggestions)
		default:
			if len(suggestions) == 0 {
				fmt.Printf("No suggestions found for: %s\n", partial)
				return nil
			}
			
			fmt.Printf("Suggestions for '%s':\n", partial)
			for _, suggestion := range suggestions {
				fmt.Printf("  %s\n", suggestion)
			}
		}
		
		return nil
	},
}

// Helper function to safely get string from map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(searchCmd)
	
	searchCmd.AddCommand(searchUniversalCmd)
	searchCmd.AddCommand(searchToolsCmd)
	searchCmd.AddCommand(searchPopularCmd)
	searchCmd.AddCommand(searchRecommendationsCmd)
	searchCmd.AddCommand(searchSuggestCmd)
	
	// Flags
	searchUniversalCmd.Flags().Int("limit", 10, "Limit number of results")
	searchToolsCmd.Flags().Int("limit", 10, "Limit number of results")
	searchPopularCmd.Flags().Int("limit", 10, "Limit number of results")
	searchRecommendationsCmd.Flags().Int("limit", 5, "Limit number of results")
	searchSuggestCmd.Flags().Int("limit", 5, "Limit number of suggestions")
}