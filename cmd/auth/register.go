package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newRegisterCmd() *cobra.Command {
	var username, email, password string

	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new Aphelion Gateway account",
		Long:  "Create a new account with Aphelion Gateway",
		Example: `  # Register interactively
  aphelion auth register

  # Register with flags
  aphelion auth register --username myuser --email my@email.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient()
			reader := bufio.NewReader(os.Stdin)

			if username == "" {
				fmt.Print("Username: ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read username: %w", err)
				}
				username = strings.TrimSpace(input)
			}

			if email == "" {
				fmt.Print("Email: ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read email: %w", err)
				}
				email = strings.TrimSpace(input)
			}

			if password == "" {
				fmt.Print("Password: ")
				passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return fmt.Errorf("failed to read password: %w", err)
				}
				password = string(passwordBytes)
				fmt.Println()

				fmt.Print("Confirm Password: ")
				confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return fmt.Errorf("failed to read password confirmation: %w", err)
				}
				confirm := string(confirmBytes)
				fmt.Println()

				if password != confirm {
					return fmt.Errorf("passwords do not match")
				}
			}

			if username == "" || email == "" || password == "" {
				return fmt.Errorf("username, email, and password are required")
			}

			spinner := utils.NewSpinner("Creating account...")
			spinner.Start()

			registerReq := api.RegisterRequest{
				Username: username,
				Email:    email,
				Password: password,
			}

			var authResp api.AuthResponse
			err := client.Post("/auth/register", registerReq, &authResp)
			spinner.Stop()

			if err != nil {
				return fmt.Errorf("registration failed: %w", err)
			}

			if err := config.SetAuth(authResp.Token, authResp.User.ID, authResp.User.Email, authResp.User.Username); err != nil {
				return fmt.Errorf("failed to save authentication: %w", err)
			}

			utils.PrintSuccess("Successfully registered and logged in as %s", authResp.User.Username)
			return nil
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "username for the new account")
	cmd.Flags().StringVarP(&email, "email", "e", "", "email address for the new account")
	cmd.Flags().StringVarP(&password, "password", "p", "", "password for the new account (not recommended, use interactive mode)")

	return cmd
}