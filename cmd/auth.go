package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/exmplrai/aphelion-cli/pkg/api"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Manage authentication with the Aphelion Gateway using Auth0.`,
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Auth0",
	Long: `Login to Aphelion Gateway using Auth0 authentication.

This will open your default browser to complete the authentication flow.
Once authenticated, your session will be saved securely for future use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		color.Cyan("🔐 Starting authentication with Aphelion Gateway...")
		
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		
		color.Blue("🌐 Opening browser for authentication...")
		
		if err := client.Login(ctx); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
		
		// Verify authentication by getting user info
		userInfo, err := client.GetUserInfo()
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}
		
		color.Green("✅ Successfully authenticated as %s (%s)", userInfo.Name, userInfo.Email)
		
		return nil
	},
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display current authentication status and user information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		if !client.IsAuthenticated() {
			color.Yellow("Not authenticated. Run 'aphelion auth login' to authenticate.")
			return nil
		}
		
		userInfo, err := client.GetUserInfo()
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}
		
		color.Green("✓ Authenticated")
		color.White("User: %s (%s)", userInfo.Name, userInfo.Email)
		if userInfo.EmailVerified {
			color.Green("Email Verified: %t", userInfo.EmailVerified)
		} else {
			color.Red("Email Verified: %t", userInfo.EmailVerified)
		}
		
		return nil
	},
}

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user information",
	Long:  `Display detailed information about the currently authenticated user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated, run 'aphelion auth login'")
		}
		
		userInfo, err := client.GetUserInfo()
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}
		
		switch output {
		case "json":
			return outputJSON(userInfo)
		case "yaml":
			return outputYAML(userInfo)
		default:
			fmt.Printf("ID: %s\n", userInfo.Sub)
			fmt.Printf("Name: %s\n", userInfo.Name)
			fmt.Printf("Email: %s\n", userInfo.Email)
			if userInfo.EmailVerified {
			color.Green("Email Verified: %t", userInfo.EmailVerified)
		} else {
			color.Red("Email Verified: %t", userInfo.EmailVerified)
		}
			if userInfo.Picture != "" {
				fmt.Printf("Picture: %s\n", userInfo.Picture)
			}
		}
		
		return nil
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear authentication session",
	Long:  `Remove stored authentication tokens and clear the current session.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		
		if err := client.Logout(); err != nil {
			return fmt.Errorf("logout failed: %w", err)
		}
		
		color.Green("✓ Successfully logged out")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(whoamiCmd)
	authCmd.AddCommand(logoutCmd)
}