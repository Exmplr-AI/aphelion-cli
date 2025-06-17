package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from Aphelion Gateway",
		Long:  "Clear stored authentication credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.ClearAuth(); err != nil {
				return fmt.Errorf("failed to clear authentication: %w", err)
			}

			utils.PrintSuccess("Successfully logged out")
			return nil
		},
	}

	return cmd
}