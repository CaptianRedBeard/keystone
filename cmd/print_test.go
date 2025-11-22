package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
)

func TestPrint(t *testing.T) {
	var buf bytes.Buffer

	rootCmd := &cobra.Command{}
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.PersistentFlags().Bool("json", false, "Output results in JSON format")

	childCmd := &cobra.Command{}
	childCmd.SetOut(&buf)
	childCmd.SetErr(&buf)
	rootCmd.AddCommand(childCmd)

	tests := []struct {
		name     string
		jsonFlag bool
		data     interface{}
		msg      string
		want     string
	}{
		{
			name:     "plain text",
			jsonFlag: false,
			data:     nil,
			msg:      "Hello, world!",
			want:     "Hello, world!\n",
		},
		{
			name:     "json output",
			jsonFlag: true,
			data:     map[string]string{"foo": "bar"},
			msg:      "",
			want:     "{\n  \"foo\": \"bar\"\n}\n",
		},
		{
			name:     "both json and msg (json takes precedence)",
			jsonFlag: true,
			data:     map[string]string{"x": "y"},
			msg:      "ignored",
			want:     "{\n  \"x\": \"y\"\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			// Set JSON flag
			rootCmd.PersistentFlags().Lookup("json").Value.Set(boolToString(tt.jsonFlag))

			Print(tt.data, tt.msg, childCmd)

			got := buf.String()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrintError(t *testing.T) {
	var buf bytes.Buffer

	rootCmd := &cobra.Command{}
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.PersistentFlags().Bool("json", true, "Output results in JSON format")

	childCmd := &cobra.Command{}
	childCmd.SetOut(&buf)
	childCmd.SetErr(&buf)
	rootCmd.AddCommand(childCmd)

	PrintError("testcmd", "something went wrong", childCmd)

	got := buf.String()

	var out map[string]string
	if err := json.Unmarshal([]byte(got), &out); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if out["context"] != "testcmd" || out["error"] != "something went wrong" {
		t.Errorf("got %+v, want context=testcmd, error=something went wrong", out)
	}
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
