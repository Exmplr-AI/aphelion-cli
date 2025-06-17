package main

import (
	"os"

	"github.com/Exmplr-AI/aphelion-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}