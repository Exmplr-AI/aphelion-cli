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
)

// Memory represents a memory entry
type Memory struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	SessionID string                 `json:"session_id,omitempty"`
	Summary   string                 `json:"summary"`
	Content   string                 `json:"content,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

// MemoryStats represents memory statistics
type MemoryStats struct {
	TotalMemories int `json:"total_memories"`
	TotalSessions int `json:"total_sessions"`
	OldestMemory  string `json:"oldest_memory"`
	NewestMemory  string `json:"newest_memory"`
}

// memoryCmd represents the memory command
var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Manage memory operations",
	Long:  `Manage memory operations including viewing, searching, and managing session memories.`,
}

// memoryListCmd represents the memory list command
var memoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List memory summaries",
	Long:  `List all memory summaries for the current user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		limit, _ := cmd.Flags().GetInt("limit")
		
		endpoint := "/memory"
		if limit > 0 {
			endpoint = fmt.Sprintf("/memory/paginated?limit=%d", limit)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to list memories: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var memories []Memory
		if limit > 0 {
			var result struct {
				Memories []Memory `json:"memories"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
			memories = result.Memories
		} else {
			if err := json.NewDecoder(resp.Body).Decode(&memories); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
		}
		
		switch output {
		case "json":
			return outputJSON(memories)
		case "yaml":
			return outputYAML(memories)
		default:
			if len(memories) == 0 {
				fmt.Println("No memories found.")
				return nil
			}
			
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Session", "Summary", "Created"})
			table.SetBorder(false)
			table.SetColWidth(50)
			table.SetRowLine(true)
			
			for _, memory := range memories {
				sessionID := memory.SessionID
				if sessionID == "" {
					sessionID = "global"
				}
				
				summary := memory.Summary
				if len(summary) > 47 {
					summary = summary[:47] + "..."
				}
				
				table.Append([]string{
					memory.ID,
					sessionID,
					summary,
					memory.CreatedAt,
				})
			}
			
			table.Render()
		}
		
		return nil
	},
}

// memorySearchCmd represents the memory search command
var memorySearchCmd = &cobra.Command{
	Use:   "search QUERY",
	Short: "Search memories using semantic similarity",
	Long:  `Search through user memories using semantic similarity.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		limit, _ := cmd.Flags().GetInt("limit")
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		
		endpoint := fmt.Sprintf("/memory/search?q=%s", query)
		if limit > 0 {
			endpoint += fmt.Sprintf("&limit=%d", limit)
		}
		if threshold > 0 {
			endpoint += fmt.Sprintf("&threshold=%.2f", threshold)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to search memories: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var memories []Memory
		if err := json.NewDecoder(resp.Body).Decode(&memories); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(memories)
		case "yaml":
			return outputYAML(memories)
		default:
			if len(memories) == 0 {
				fmt.Printf("No memories found matching query: %s\n", query)
				return nil
			}
			
			fmt.Printf("Found %d memories matching: %s\n\n", len(memories), query)
			
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Session", "Summary", "Created"})
			table.SetBorder(false)
			table.SetColWidth(50)
			table.SetRowLine(true)
			
			for _, memory := range memories {
				sessionID := memory.SessionID
				if sessionID == "" {
					sessionID = "global"
				}
				
				summary := memory.Summary
				if len(summary) > 47 {
					summary = summary[:47] + "..."
				}
				
				table.Append([]string{
					memory.ID,
					sessionID,
					summary,
					memory.CreatedAt,
				})
			}
			
			table.Render()
		}
		
		return nil
	},
}

// memorySessionsCmd represents the memory sessions command
var memorySessionsCmd = &cobra.Command{
	Use:   "sessions SESSION_ID",
	Short: "Get memory summary for a session",
	Long:  `Get memory summary for a specific session.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, fmt.Sprintf("/memory/sessions/%s", sessionID))
		if err != nil {
			return fmt.Errorf("failed to get session memory: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var memory Memory
		if err := json.NewDecoder(resp.Body).Decode(&memory); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(memory)
		case "yaml":
			return outputYAML(memory)
		default:
			fmt.Printf("Session: %s\n", memory.SessionID)
			fmt.Printf("Created: %s\n", memory.CreatedAt)
			fmt.Printf("Updated: %s\n", memory.UpdatedAt)
			fmt.Printf("Summary:\n%s\n", memory.Summary)
			
			if memory.Content != "" {
				fmt.Printf("\nContent:\n%s\n", memory.Content)
			}
			
			if memory.Metadata != nil && len(memory.Metadata) > 0 {
				fmt.Println("\nMetadata:")
				metadataJSON, _ := json.MarshalIndent(memory.Metadata, "  ", "  ")
				fmt.Printf("  %s\n", string(metadataJSON))
			}
		}
		
		return nil
	},
}

// memorySummarizeCmd represents the memory summarize command
var memorySummarizeCmd = &cobra.Command{
	Use:   "summarize SESSION_ID",
	Short: "Create memory summary for a session",
	Long:  `Create a memory summary for session activities.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		force, _ := cmd.Flags().GetBool("force")
		
		payload := map[string]interface{}{
			"force": force,
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		
		resp, err := client.Post(ctx, fmt.Sprintf("/memory/sessions/%s/summarize", sessionID), payload)
		if err != nil {
			return fmt.Errorf("failed to create session summary: %w", err)
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
			if success, ok := result["success"].(bool); ok && success {
				fmt.Printf("✓ Created memory summary for session: %s\n", sessionID)
			} else {
				fmt.Printf("Failed to create memory summary for session: %s\n", sessionID)
			}
			
			if message, ok := result["message"].(string); ok {
				fmt.Println(message)
			}
		}
		
		return nil
	},
}

// memoryStatsCmd represents the memory stats command
var memoryStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get memory statistics",
	Long:  `Get user memory statistics and usage information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, "/memory/stats")
		if err != nil {
			return fmt.Errorf("failed to get memory stats: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var stats MemoryStats
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(stats)
		case "yaml":
			return outputYAML(stats)
		default:
			fmt.Printf("Total Memories: %d\n", stats.TotalMemories)
			fmt.Printf("Total Sessions: %d\n", stats.TotalSessions)
			fmt.Printf("Oldest Memory: %s\n", stats.OldestMemory)
			fmt.Printf("Newest Memory: %s\n", stats.NewestMemory)
		}
		
		return nil
	},
}

// memoryDeleteCmd represents the memory delete command
var memoryDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete all memories",
	Long:  `Delete all memories for the current user. This cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Delete(ctx, "/memory")
		if err != nil {
			return fmt.Errorf("failed to delete memories: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		fmt.Println("✓ All memories deleted successfully")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(memoryCmd)
	
	memoryCmd.AddCommand(memoryListCmd)
	memoryCmd.AddCommand(memorySearchCmd)
	memoryCmd.AddCommand(memorySessionsCmd)
	memoryCmd.AddCommand(memorySummarizeCmd)
	memoryCmd.AddCommand(memoryStatsCmd)
	memoryCmd.AddCommand(memoryDeleteCmd)
	
	// Flags
	memoryListCmd.Flags().Int("limit", 0, "Limit number of results")
	
	memorySearchCmd.Flags().Int("limit", 10, "Limit number of results")
	memorySearchCmd.Flags().Float64("threshold", 0.0, "Similarity threshold (0.0-1.0)")
	
	memorySummarizeCmd.Flags().Bool("force", false, "Force creation of summary even if one exists")
}