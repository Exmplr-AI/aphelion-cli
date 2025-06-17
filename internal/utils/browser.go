package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenBrowser opens the specified URL in the default browser
func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// OpenBrowserWithFallback attempts to open browser and provides fallback instructions
func OpenBrowserWithFallback(url string) {
	PrintInfo("Opening browser for authentication...")
	
	if err := OpenBrowser(url); err != nil {
		PrintWarning("Could not open browser automatically")
		fmt.Printf("\nPlease open the following URL in your browser:\n\n%s\n\n", url)
	} else {
		fmt.Printf("If the browser doesn't open automatically, please visit:\n%s\n\n", url)
	}
}