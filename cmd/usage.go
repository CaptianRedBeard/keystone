package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	daysBack int
)

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show API usage and statistics",
	Long:  "Displays usage metrics such as token counts and request totals.",
}

var usageTodayCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show recent usage summary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Showing usage summary for the past %d day(s):\n", daysBack)
		fmt.Println(" - Requests: 12")
		fmt.Println(" - Tokens used: 2,345")
	},
}

func init() {
	rootCmd.AddCommand(usageCmd)
	usageCmd.AddCommand(usageTodayCmd)

	usageTodayCmd.Flags().IntVarP(
		&daysBack,
		"days", "d",
		1,
		"number of days to look back for usage summary",
	)
}
