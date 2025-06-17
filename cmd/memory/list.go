package memory

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Exmplr-AI/aphelion-cli/pkg/api"
	"github.com/Exmplr-AI/aphelion-cli/pkg/config"
	"github.com/Exmplr-AI/aphelion-cli/internal/utils"
)

func newListCmd() *cobra.Command {
	var limit int
	var sort string
	var dateFrom, dateTo string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List your memories",
		Long:  "List all memories with pagination and sorting options",
		Example: `  # List recent memories
  aphelion memory list

  # List with custom limit and sorting
  aphelion memory list --limit 10 --sort oldest

  # List memories from specific date range
  aphelion memory list --date-from 2023-01-01 --date-to 2023-12-31`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !config.IsAuthenticated() {
				return fmt.Errorf("authentication required. Please run 'aphelion auth login' first")
			}

			client := api.NewClient()
			
			// Try paginated endpoint first, fallback to basic endpoint
			var memories []api.Memory
			
			if limit != 10 || sort != "newest" || dateFrom != "" || dateTo != "" {
				// Use paginated endpoint when filters are specified
				params := map[string]string{
					"limit": strconv.Itoa(limit),
					"sort":  sort,
				}

				if dateFrom != "" {
					params["date_from"] = dateFrom
				}
				if dateTo != "" {
					params["date_to"] = dateTo
				}

				var response api.MemoriesResponse
				if err := client.GetWithQuery("/memory/paginated", params, &response); err != nil {
					// Fallback to basic endpoint
					if err := client.Get("/memory", &memories); err != nil {
						return fmt.Errorf("failed to list memories: %w", err)
					}
				} else {
					memories = response.Memories
				}
			} else {
				// Use basic endpoint for simple list
				var response api.MemoriesResponse
				if err := client.Get("/memory", &response); err != nil {
					return fmt.Errorf("failed to list memories: %w", err)
				}
				memories = response.Memories
			}

			if len(memories) == 0 {
				utils.PrintInfo("No memories found")
				return nil
			}

			utils.PrintInfo("Found %d memories", len(memories))
			
			var data []map[string]interface{}
			for _, memory := range memories {
				data = append(data, map[string]interface{}{
					"ID":         memory.ID,
					"Session ID": memory.SessionID,
					"Summary":    memory.Summary,
					"Created":    memory.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}

			return utils.PrintOutput(data, config.GetOutputFormat())
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 10, "number of memories to return (1-100)")
	cmd.Flags().StringVarP(&sort, "sort", "s", "newest", "sort order (newest, oldest)")
	cmd.Flags().StringVar(&dateFrom, "date-from", "", "filter from date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&dateTo, "date-to", "", "filter to date (YYYY-MM-DD)")

	return cmd
}