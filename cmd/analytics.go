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

// analyticsCmd represents the analytics command
var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "View usage analytics and metrics",
	Long:  `View analytics, usage metrics, and performance monitoring data.`,
}

// analyticsUserCmd represents the analytics user command
var analyticsUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get user analytics",
	Long:  `Get user-specific analytics including request metrics, session statistics, and tool usage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		timeframe, _ := cmd.Flags().GetString("timeframe")
		
		endpoint := "/analytics/user"
		if timeframe != "" {
			endpoint += fmt.Sprintf("?timeframe=%s", timeframe)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to get user analytics: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var analytics map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(analytics)
		case "yaml":
			return outputYAML(analytics)
		default:
			displayAnalytics("User Analytics", analytics)
		}
		
		return nil
	},
}

// analyticsSystemCmd represents the analytics system command
var analyticsSystemCmd = &cobra.Command{
	Use:   "system",
	Short: "Get system analytics",
	Long:  `Get system-wide analytics including total requests, unique users, and performance metrics.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		timeframe, _ := cmd.Flags().GetString("timeframe")
		
		endpoint := "/analytics/system"
		if timeframe != "" {
			endpoint += fmt.Sprintf("?timeframe=%s", timeframe)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to get system analytics: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var analytics map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(analytics)
		case "yaml":
			return outputYAML(analytics)
		default:
			displayAnalytics("System Analytics", analytics)
		}
		
		return nil
	},
}

// analyticsToolsCmd represents the analytics tools command
var analyticsToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Get tool usage analytics",
	Long:  `Get tool usage analytics including most popular tools, success rates, and error rates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		timeframe, _ := cmd.Flags().GetString("timeframe")
		userOnly, _ := cmd.Flags().GetBool("user-only")
		
		endpoint := "/analytics/tools"
		params := make([]string, 0)
		
		if timeframe != "" {
			params = append(params, fmt.Sprintf("timeframe=%s", timeframe))
		}
		if userOnly {
			params = append(params, "user_only=true")
		}
		
		if len(params) > 0 {
			endpoint += "?"
			for i, param := range params {
				if i > 0 {
					endpoint += "&"
				}
				endpoint += param
			}
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to get tool analytics: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var analytics map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(analytics)
		case "yaml":
			return outputYAML(analytics)
		default:
			title := "Tool Analytics"
			if userOnly {
				title = "Your Tool Usage Analytics"
			}
			displayAnalytics(title, analytics)
		}
		
		return nil
	},
}

// analyticsSessionsCmd represents the analytics sessions command
var analyticsSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Get session analytics",
	Long:  `Get session analytics including session counts, durations, and activity.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		timeframe, _ := cmd.Flags().GetString("timeframe")
		userOnly, _ := cmd.Flags().GetBool("user-only")
		
		endpoint := "/analytics/sessions"
		params := make([]string, 0)
		
		if timeframe != "" {
			params = append(params, fmt.Sprintf("timeframe=%s", timeframe))
		}
		if userOnly {
			params = append(params, "user_only=true")
		}
		
		if len(params) > 0 {
			endpoint += "?"
			for i, param := range params {
				if i > 0 {
					endpoint += "&"
				}
				endpoint += param
			}
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to get session analytics: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var analytics map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(analytics)
		case "yaml":
			return outputYAML(analytics)
		default:
			title := "Session Analytics"
			if userOnly {
				title = "Your Session Analytics"
			}
			displayAnalytics(title, analytics)
		}
		
		return nil
	},
}

// analyticsEndpointsCmd represents the analytics endpoints command
var analyticsEndpointsCmd = &cobra.Command{
	Use:   "endpoints",
	Short: "Get endpoint analytics",
	Long:  `Get endpoint popularity analytics including request counts, response times, and error rates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		timeframe, _ := cmd.Flags().GetString("timeframe")
		
		endpoint := "/analytics/endpoints"
		if timeframe != "" {
			endpoint += fmt.Sprintf("?timeframe=%s", timeframe)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, endpoint)
		if err != nil {
			return fmt.Errorf("failed to get endpoint analytics: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var analytics map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(analytics)
		case "yaml":
			return outputYAML(analytics)
		default:
			displayAnalytics("Endpoint Analytics", analytics)
		}
		
		return nil
	},
}

// analyticsRealtimeCmd represents the analytics realtime command
var analyticsRealtimeCmd = &cobra.Command{
	Use:   "realtime",
	Short: "Get real-time metrics",
	Long:  `Get real-time system metrics including active requests, users, and performance in the last 5 minutes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		resp, err := client.Get(ctx, "/analytics/realtime")
		if err != nil {
			return fmt.Errorf("failed to get real-time metrics: %w", err)
		}
		defer resp.Body.Close()
		
		if err := api.HandleErrorResponse(resp); err != nil {
			return err
		}
		
		var metrics map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(metrics)
		case "yaml":
			return outputYAML(metrics)
		default:
			displayAnalytics("Real-time Metrics", metrics)
		}
		
		return nil
	},
}

// displayAnalytics displays analytics data in a formatted table
func displayAnalytics(title string, data map[string]interface{}) {
	fmt.Printf("%s:\n\n", title)
	
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Value"})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	
	for key, value := range data {
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case int:
			valueStr = fmt.Sprintf("%d", v)
		case float64:
			valueStr = fmt.Sprintf("%.2f", v)
		case bool:
			valueStr = fmt.Sprintf("%t", v)
		default:
			jsonValue, _ := json.Marshal(v)
			valueStr = string(jsonValue)
		}
		
		table.Append([]string{key, valueStr})
	}
	
	table.Render()
}

func init() {
	rootCmd.AddCommand(analyticsCmd)
	
	analyticsCmd.AddCommand(analyticsUserCmd)
	analyticsCmd.AddCommand(analyticsSystemCmd)
	analyticsCmd.AddCommand(analyticsToolsCmd)
	analyticsCmd.AddCommand(analyticsSessionsCmd)
	analyticsCmd.AddCommand(analyticsEndpointsCmd)
	analyticsCmd.AddCommand(analyticsRealtimeCmd)
	
	// Flags
	analyticsUserCmd.Flags().String("timeframe", "", "Time period: hour, day, week, month")
	analyticsSystemCmd.Flags().String("timeframe", "", "Time period: hour, day, week, month")
	analyticsToolsCmd.Flags().String("timeframe", "", "Time period: hour, day, week, month")
	analyticsToolsCmd.Flags().Bool("user-only", false, "Show only current user's tool usage")
	analyticsSessionsCmd.Flags().String("timeframe", "", "Time period: hour, day, week, month")
	analyticsSessionsCmd.Flags().Bool("user-only", false, "Show only current user's sessions")
	analyticsEndpointsCmd.Flags().String("timeframe", "", "Time period: hour, day, week, month")
}