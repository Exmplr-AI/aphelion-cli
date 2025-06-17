package agent

import (
	"github.com/spf13/cobra"
)

func NewAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Agent development and execution commands",
		Long:  "Create, initialize, and run AI agents with Aphelion Gateway",
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newRunCmd())

	return cmd
}