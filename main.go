package main

import (
	"fmt"
	"keystone/cmd"
	"keystone/internal/config"
	"os"
)

func main() {
	cfgFile := cmd.GetConfigPath()

	if err := config.LoadConfig(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Failed to load config: %v\n", err)
		os.Exit(1)
	}

	cmd.Execute()
}
