package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUsageCmd() *cobra.Command {
	var daysBack int

	usageCmd := &cobra.Command{
		Use:   "usage",
		Short: "Show API usage and statistics",
		Long:  "Displays usage metrics such as token counts and request totals.",
	}

	summaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Show recent usage summary",
		Run: func(cmd *cobra.Command, args []string) {
			// Example usage metrics
			out := map[string]interface{}{
				"days":       daysBack,
				"requests":   12,
				"tokensUsed": 2345,
			}

			// Print output respecting --json
			if jsonFlag, _ := cmd.Root().PersistentFlags().GetBool("json"); jsonFlag {
				Print(out, "", cmd)
			} else {
				msg := fmt.Sprintf(
					"Showing usage summary for the past %d day(s):\n - Requests: %d\n - Tokens used: %d",
					daysBack, out["requests"], out["tokensUsed"],
				)
				Print(nil, msg, cmd)
			}
		},
	}

	summaryCmd.Flags().IntVarP(&daysBack, "days", "d", 1, "number of days to look back for usage summary")
	usageCmd.AddCommand(summaryCmd)

	return usageCmd
}
