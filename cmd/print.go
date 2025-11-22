package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// getJSONFlag returns true if the --json flag is set on the command.
func getJSONFlag(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	if f := cmd.Root().PersistentFlags().Lookup("json"); f != nil {
		val, _ := cmd.Root().PersistentFlags().GetBool("json")
		return val
	}
	return false
}

// Print outputs obj as JSON if --json is set, otherwise prints msg.
func Print(obj interface{}, msg string, cmd *cobra.Command) {
	out := io.Writer(os.Stdout)
	if cmd != nil {
		out = cmd.OutOrStdout()
	}

	if getJSONFlag(cmd) && obj != nil {
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			fmt.Fprintf(out, "Error marshaling JSON output: %v\n", err)
			return
		}
		fmt.Fprintln(out, string(data))
	} else if msg != "" {
		fmt.Fprintln(out, msg)
	}
}

// PrintError prints an error in JSON if --json is set, otherwise human-readable.
func PrintError(context, errMsg string, cmd *cobra.Command) {
	out := io.Writer(os.Stderr)
	if cmd != nil {
		out = cmd.ErrOrStderr()
	}

	if getJSONFlag(cmd) {
		data, err := json.MarshalIndent(map[string]string{
			"context": context,
			"error":   errMsg,
		}, "", "  ")
		if err != nil {
			fmt.Fprintf(out, "Error marshaling JSON error: %v\n", err)
			return
		}
		fmt.Fprintln(out, string(data))
	} else {
		fmt.Fprintf(out, "Error [%s]: %s\n", context, errMsg)
	}
}
