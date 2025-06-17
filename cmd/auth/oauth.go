package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
)

func newOAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oauth",
		Short: "Get Auth0 OAuth information",
		Long:  "Display Auth0 configuration and authentication URL for OAuth flow",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient()
			
			var authInfo map[string]interface{}
			if err := client.Get("/auth/info", &authInfo); err != nil {
				return fmt.Errorf("failed to get Auth0 info: %w", err)
			}

			fmt.Println("Auth0 OAuth Configuration:")
			fmt.Println("========================")
			
			if domain, ok := authInfo["auth0_domain"].(string); ok {
				fmt.Printf("Domain: %s\n", domain)
			}
			
			clientID := "UbXxpQBSr9AsqpS2ln2jzmsmamromaFC" // Correct client_id
			if cid, ok := authInfo["client_id"].(string); ok && cid != "" {
				clientID = cid
			}
			fmt.Printf("Client ID: %s\n", clientID)
			
			if audience, ok := authInfo["auth0_audience"].(string); ok {
				fmt.Printf("Audience: %s\n", audience)
			}
			
			fmt.Printf("Redirect URI: %s\n", "http://localhost:8765/callback")

			fmt.Println("\nTo authenticate with Auth0:")
			fmt.Println("1. Use 'aphelion auth login' for automatic browser-based authentication")
			fmt.Println("2. Or use 'aphelion auth login --no-launch-browser' to get the URL manually")

			if authURL, ok := authInfo["auth_url"].(string); ok {
				fullAuthURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=openid%%20profile%%20email&audience=%s",
					authURL,
					clientID,
					"http://localhost:8765/callback",
					authInfo["auth0_audience"],
				)
				fmt.Printf("\nSample Authorization URL:\n%s\n", fullAuthURL)
			}

			return nil
		},
	}

	return cmd
}