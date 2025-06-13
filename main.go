package main

import (
	"fmt"
	"os"

	"github.com/exmplrai/aphelion-cli/cmd"
)

var (
	// Build information, set by ldflags during build
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Set build info for version command
	cmd.SetBuildInfo(Version, GitCommit, BuildDate)
	
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}