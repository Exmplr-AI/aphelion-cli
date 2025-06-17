package tools

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
)

var (
	dryRun bool
	params string
)

type ToolExecutionRequest struct {
	Parameters map[string]interface{} `json:"parameters"`
}

type ToolExecutionResult struct {
	Success   bool                   `json:"success"`
	Result    interface{}            `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Duration  string                 `json:"duration,omitempty"`
}

func newTryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "try --tool <tool-name> --params '<json>'",
		Short: "Execute a tool with given parameters",
		Long:  "Validate and execute a tool with provided parameters",
		RunE:  runTry,
	}

	cmd.Flags().StringVar(&params, "params", "{}", "Tool parameters as JSON string")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate parameters without executing")
	cmd.Flags().String("tool", "", "Tool name to execute")
	
	cmd.MarkFlagRequired("tool")

	return cmd
}

func runTry(cmd *cobra.Command, args []string) error {
	toolName, _ := cmd.Flags().GetString("tool")
	
	// Parse parameters
	var parameters map[string]interface{}
	if err := json.Unmarshal([]byte(params), &parameters); err != nil {
		return fmt.Errorf("invalid JSON parameters: %w", err)
	}

	client := api.NewClient()

	if dryRun {
		return validateToolParameters(client, toolName, parameters)
	}

	return executeToolWithParameters(client, toolName, parameters)
}

func validateToolParameters(client *api.Client, toolName string, parameters map[string]interface{}) error {
	endpoint := fmt.Sprintf("/tools/%s/validate", toolName)
	
	request := ToolExecutionRequest{
		Parameters: parameters,
	}

	var result map[string]interface{}
	if err := client.Post(endpoint, request, &result); err != nil {
		return fmt.Errorf("parameter validation failed: %w", err)
	}

	fmt.Println("✅ Parameters are valid")
	
	if viper.GetBool("verbose") {
		fmt.Println("Validation details:")
		return utils.OutputJSON(result)
	}
	
	return nil
}

func executeToolWithParameters(client *api.Client, toolName string, parameters map[string]interface{}) error {
	fmt.Printf("🚀 Executing tool: %s\n", toolName)
	
	if viper.GetBool("verbose") {
		fmt.Printf("📋 Parameters: ")
		paramsJSON, _ := json.MarshalIndent(parameters, "", "  ")
		fmt.Println(string(paramsJSON))
		fmt.Println()
	}

	spinner := utils.NewSpinner("Executing tool...")
	spinner.Start()
	defer spinner.Stop()

	endpoint := fmt.Sprintf("/tools/%s/execute", toolName)
	
	request := ToolExecutionRequest{
		Parameters: parameters,
	}

	var result ToolExecutionResult
	if err := client.Post(endpoint, request, &result); err != nil {
		spinner.Stop()
		return fmt.Errorf("tool execution failed: %w", err)
	}

	spinner.Stop()
	return outputExecutionResult(result)
}

func outputExecutionResult(result ToolExecutionResult) error {
	outputFormat := viper.GetString("output")
	
	switch outputFormat {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
		
	case "yaml":
		return utils.OutputYAML(result)
		
	default: // table format
		return outputExecutionResultTable(result)
	}
}

func outputExecutionResultTable(result ToolExecutionResult) error {
	if result.Success {
		fmt.Println("✅ Tool executed successfully")
	} else {
		fmt.Println("❌ Tool execution failed")
		if result.Error != "" {
			fmt.Printf("Error: %s\n", result.Error)
		}
		return nil
	}
	
	if result.Duration != "" {
		fmt.Printf("⏱️  Duration: %s\n", result.Duration)
	}
	
	if result.Result != nil {
		fmt.Println("\n📋 Result:")
		resultJSON, err := json.MarshalIndent(result.Result, "", "  ")
		if err != nil {
			fmt.Printf("Raw result: %v\n", result.Result)
		} else {
			fmt.Println(string(resultJSON))
		}
	}
	
	if len(result.Metadata) > 0 && viper.GetBool("verbose") {
		fmt.Println("\n🔍 Metadata:")
		metadataJSON, _ := json.MarshalIndent(result.Metadata, "", "  ")
		fmt.Println(string(metadataJSON))
	}
	
	return nil
}