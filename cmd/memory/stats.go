package memory

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show memory statistics",
		Long:  "Display usage statistics for your memories",
		Example: `  # Show memory statistics
  aphelion memory stats

  # Show statistics in JSON format
  aphelion memory stats --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !config.IsAuthenticated() {
				return fmt.Errorf("authentication required. Please run 'aphelion auth login' first")
			}

			client := api.NewClient()
			
			var stats api.MemoryStats
			if err := client.Get("/memory/stats", &stats); err != nil {
				return fmt.Errorf("failed to get memory statistics: %w", err)
			}

			data := map[string]interface{}{
				"Total Memories":     stats.TotalMemories,
				"Total Sessions":     stats.TotalSessions,
				"Average Per Day":    fmt.Sprintf("%.2f", stats.AveragePerDay),
				"Oldest Memory":      stats.OldestMemory,
				"Most Recent Memory": stats.MostRecentMemory,
			}

			return utils.PrintOutput(data, config.GetOutputFormat())
		},
	}

	return cmd
}