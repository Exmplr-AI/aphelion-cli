package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/exmplrai/aphelion-cli/internal/config"
	"github.com/exmplrai/aphelion-cli/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	profile     string
	verbose     bool
	output      string
	endpoint    string
	
	// Build information
	version   string
	gitCommit string
	buildDate string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aphelion",
	Short: "Aphelion CLI - Interact with the Aphelion Gateway API",
	Long: `Aphelion CLI is a command-line tool for interacting with the Aphelion Gateway,
a unified AI agent platform with tool discovery and memory capabilities.

This CLI provides a gcloud-like experience for managing:
- Authentication (automatic Auth0 configuration)
- Agent sessions and tool execution
- Service registry and API management
- Memory operations and semantic search
- Analytics and usage metrics

Get started in seconds:
  aphelion auth login

No configuration needed - the CLI automatically discovers all Auth0 settings!

For more information, visit: https://github.com/exmplrai/aphelion-cli`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize logger
		if verbose {
			logger.SetLevel("debug")
		}
		
		// Initialize configuration
		return config.Initialize(cfgFile, profile)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

// SetBuildInfo sets build information for version command
func SetBuildInfo(v, commit, date string) {
	version = v
	gitCommit = commit
	buildDate = date
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aphelion/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "default", "configuration profile to use")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "output format (table|json|yaml)")
	rootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "Aphelion Gateway endpoint URL")

	// Bind flags to viper
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".aphelion" (without extension).
		configDir := fmt.Sprintf("%s/.aphelion", home)
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Environment variables
	viper.SetEnvPrefix("APHELION")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("endpoint", "https://simplified-aphelion-api-172201620564.us-central1.run.app")
	viper.SetDefault("output", "table")
	viper.SetDefault("verbose", false)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Debug("Using config file:", viper.ConfigFileUsed())
	}
}