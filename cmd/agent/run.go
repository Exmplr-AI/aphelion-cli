package agent

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var (
	cronSchedule string
	daemon       bool
	verbose      bool
)

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [agent-file]",
		Short: "Run an agent",
		Long:  "Execute an agent script with optional cron scheduling",
		Args:  cobra.ExactArgs(1),
		RunE:  runAgent,
	}

	cmd.Flags().StringVar(&cronSchedule, "cron", "", "Cron schedule for agent execution (e.g., '*/10 * * * *')")
	cmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "Run agent as daemon")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func runAgent(cmd *cobra.Command, args []string) error {
	agentFile := args[0]

	// Validate agent file exists
	if _, err := os.Stat(agentFile); os.IsNotExist(err) {
		return fmt.Errorf("agent file not found: %s", agentFile)
	}

	// Make file path absolute
	absPath, err := filepath.Abs(agentFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if cronSchedule != "" {
		return runWithCron(absPath)
	}

	if daemon {
		return runAsDaemon(absPath)
	}

	return runOnce(absPath)
}

func runOnce(agentFile string) error {
	if verbose {
		fmt.Printf("🚀 Running agent: %s\n", agentFile)
	}

	cmd := createAgentCommand(agentFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runWithCron(agentFile string) error {
	fmt.Printf("📅 Scheduling agent with cron: %s\n", cronSchedule)
	fmt.Printf("🚀 Agent file: %s\n", agentFile)

	c := cron.New()
	
	_, err := c.AddFunc(cronSchedule, func() {
		if verbose {
			fmt.Printf("[%s] 🔄 Running scheduled agent execution\n", time.Now().Format("2006-01-02 15:04:05"))
		}
		
		cmd := createAgentCommand(agentFile)
		
		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Agent execution failed: %v\n", err)
		} else if verbose {
			fmt.Printf("✅ Agent execution completed successfully\n")
		}
	})

	if err != nil {
		return fmt.Errorf("invalid cron schedule: %w", err)
	}

	c.Start()
	defer c.Stop()

	fmt.Println("⏰ Cron scheduler started. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n🛑 Stopping cron scheduler...")
	return nil
}

func runAsDaemon(agentFile string) error {
	fmt.Printf("🔄 Running agent as daemon: %s\n", agentFile)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run agent in a goroutine
	done := make(chan error, 1)
	go func() {
		cmd := createAgentCommand(agentFile)
		
		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		
		done <- cmd.Run()
	}()

	// Wait for either completion or signal
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("agent execution failed: %w", err)
		}
		return nil
	case sig := <-sigChan:
		fmt.Printf("\n🛑 Received signal %v, stopping agent...\n", sig)
		return nil
	}
}

func createAgentCommand(agentFile string) *exec.Cmd {
	// Determine how to run the agent based on file extension
	ext := strings.ToLower(filepath.Ext(agentFile))
	
	switch ext {
	case ".py":
		return exec.Command("python3", agentFile)
	case ".js":
		return exec.Command("node", agentFile)
	case ".go":
		return exec.Command("go", "run", agentFile)
	default:
		// Try to make it executable and run directly
		os.Chmod(agentFile, 0755)
		return exec.Command(agentFile)
	}
}