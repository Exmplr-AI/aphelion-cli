package cmd

import (
	"fmt"
	"os"

	"github.com/Exmplr-AI/aphelion-cli/internal/config"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Manage Aphelion CLI configuration settings and profiles.`,
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set KEY VALUE",
	Short: "Set a configuration value",
	Long: `Set a configuration value for the current profile.

Available keys:
  endpoint    - Aphelion Gateway endpoint URL
  auth.domain - Auth0 domain
  auth.client_id - Auth0 client ID
  auth.audience - Auth0 audience`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]
		
		// Get current profile
		profile := config.GetCurrentProfile()
		
		// Update the specified key
		switch key {
		case "endpoint":
			profile.Endpoint = value
		case "auth.domain":
			profile.Auth.Domain = value
		case "auth.client_id":
			profile.Auth.ClientID = value
		case "auth.audience":
			profile.Auth.Audience = value
		case "auth.redirect_uri":
			profile.Auth.RedirectURI = value
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}
		
		// Save the updated profile
		if err := config.SetProfile(config.Get().CurrentProfile, profile); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
		
		color.Green("✓ Set %s = %s", key, value)
		
		return nil
	},
}

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get [KEY]",
	Short: "Get configuration value(s)",
	Long:  `Get a specific configuration value or all values if no key is specified.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetCurrentProfile()
		
		if len(args) == 0 {
			// Show all config
			switch output {
			case "json":
				return outputJSON(profile)
			case "yaml":
				return outputYAML(profile)
			default:
				fmt.Printf("Profile: %s\n", config.Get().CurrentProfile)
				fmt.Printf("Endpoint: %s\n", profile.Endpoint)
				fmt.Printf("Auth Domain: %s\n", profile.Auth.Domain)
				fmt.Printf("Auth Client ID: %s\n", profile.Auth.ClientID)
				fmt.Printf("Auth Audience: %s\n", profile.Auth.Audience)
				fmt.Printf("Auth Redirect URI: %s\n", profile.Auth.RedirectURI)
			}
			return nil
		}
		
		// Get specific key
		key := args[0]
		var value string
		
		switch key {
		case "endpoint":
			value = profile.Endpoint
		case "auth.domain":
			value = profile.Auth.Domain
		case "auth.client_id":
			value = profile.Auth.ClientID
		case "auth.audience":
			value = profile.Auth.Audience
		case "auth.redirect_uri":
			value = profile.Auth.RedirectURI
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}
		
		fmt.Println(value)
		
		return nil
	},
}

// configListCmd represents the config list command
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration keys and values",
	Long:  `List all configuration keys and values for the current profile.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetCurrentProfile()
		
		switch output {
		case "json":
			data := map[string]interface{}{
				"profile":  config.Get().CurrentProfile,
				"endpoint": profile.Endpoint,
				"auth": map[string]string{
					"domain":       profile.Auth.Domain,
					"client_id":    profile.Auth.ClientID,
					"audience":     profile.Auth.Audience,
					"redirect_uri": profile.Auth.RedirectURI,
				},
			}
			return outputJSON(data)
		case "yaml":
			data := map[string]interface{}{
				"profile":  config.Get().CurrentProfile,
				"endpoint": profile.Endpoint,
				"auth": map[string]string{
					"domain":       profile.Auth.Domain,
					"client_id":    profile.Auth.ClientID,
					"audience":     profile.Auth.Audience,
					"redirect_uri": profile.Auth.RedirectURI,
				},
			}
			return outputYAML(data)
		default:
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Key", "Value"})
			table.SetBorder(false)
			
			table.Append([]string{"profile", config.Get().CurrentProfile})
			table.Append([]string{"endpoint", profile.Endpoint})
			table.Append([]string{"auth.domain", profile.Auth.Domain})
			table.Append([]string{"auth.client_id", profile.Auth.ClientID})
			table.Append([]string{"auth.audience", profile.Auth.Audience})
			table.Append([]string{"auth.redirect_uri", profile.Auth.RedirectURI})
			
			table.Render()
		}
		
		return nil
	},
}

// configResetCmd represents the config reset command
var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Long:  `Reset the current profile configuration to default values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentProfile := config.Get().CurrentProfile
		
		defaultProfile := config.Profile{
			Name:     currentProfile,
			Endpoint: viper.GetString("endpoint"),
			Auth: config.Auth{
				Type:         "auth0",
				Domain:       "",
				ClientID:     "",
				Audience:     "",
				RedirectURI:  "http://localhost:8765/callback",
				TokenStorage: "keyring",
			},
		}
		
		if err := config.SetProfile(currentProfile, defaultProfile); err != nil {
			return fmt.Errorf("failed to reset configuration: %w", err)
		}
		
		color.Green("✓ Reset configuration for profile '%s' to defaults", currentProfile)
		
		return nil
	},
}

// configProfilesCmd represents the config profiles command
var configProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage configuration profiles",
	Long:  `Manage multiple configuration profiles for different environments.`,
}

// configProfilesListCmd represents the config profiles list command
var configProfilesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Long:  `List all available configuration profiles.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles := config.ListProfiles()
		currentProfile := config.Get().CurrentProfile
		
		switch output {
		case "json":
			data := make(map[string]interface{})
			for _, name := range profiles {
				if name == currentProfile {
					data[name] = map[string]interface{}{"current": true}
				} else {
					data[name] = map[string]interface{}{"current": false}
				}
			}
			return outputJSON(data)
		case "yaml":
			data := make(map[string]interface{})
			for _, name := range profiles {
				if name == currentProfile {
					data[name] = map[string]interface{}{"current": true}
				} else {
					data[name] = map[string]interface{}{"current": false}
				}
			}
			return outputYAML(data)
		default:
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Profile", "Current"})
			table.SetBorder(false)
			
			for _, name := range profiles {
				current := ""
				if name == currentProfile {
					current = "✓"
				}
				table.Append([]string{name, current})
			}
			
			table.Render()
		}
		
		return nil
	},
}

// configProfilesCreateCmd represents the config profiles create command
var configProfilesCreateCmd = &cobra.Command{
	Use:   "create NAME",
	Short: "Create a new profile",
	Long:  `Create a new configuration profile with default values.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		// Check if profile already exists
		profiles := config.ListProfiles()
		for _, existing := range profiles {
			if existing == name {
				return fmt.Errorf("profile '%s' already exists", name)
			}
		}
		
		// Create new profile with defaults
		newProfile := config.Profile{
			Name:     name,
			Endpoint: viper.GetString("endpoint"),
			Auth: config.Auth{
				Type:         "auth0",
				Domain:       "",
				ClientID:     "",
				Audience:     "",
				RedirectURI:  "http://localhost:8765/callback",
				TokenStorage: "keyring",
			},
		}
		
		if err := config.SetProfile(name, newProfile); err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}
		
		color.Green("✓ Created profile '%s'", name)
		
		return nil
	},
}

// configProfilesSwitchCmd represents the config profiles switch command
var configProfilesSwitchCmd = &cobra.Command{
	Use:   "switch NAME",
	Short: "Switch to a different profile",
	Long:  `Switch to a different configuration profile.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		if err := config.SwitchProfile(name); err != nil {
			return fmt.Errorf("failed to switch profile: %w", err)
		}
		
		color.Green("✓ Switched to profile '%s'", name)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	
	// Main config commands
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configResetCmd)
	
	// Profile management
	configCmd.AddCommand(configProfilesCmd)
	configProfilesCmd.AddCommand(configProfilesListCmd)
	configProfilesCmd.AddCommand(configProfilesCreateCmd)
	configProfilesCmd.AddCommand(configProfilesSwitchCmd)
}