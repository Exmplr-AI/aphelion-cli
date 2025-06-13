package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	CurrentProfile string             `yaml:"current_profile"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

// Profile represents a configuration profile
type Profile struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Auth     Auth   `yaml:"auth"`
}

// Auth represents authentication configuration
type Auth struct {
	Type         string `yaml:"type"`          // "auth0"
	Domain       string `yaml:"domain"`        // Auth0 domain
	ClientID     string `yaml:"client_id"`     // Auth0 client ID
	Audience     string `yaml:"audience"`      // Auth0 audience
	RedirectURI  string `yaml:"redirect_uri"`  // Auth0 redirect URI
	TokenStorage string `yaml:"token_storage"` // Token storage method
}

var (
	configInstance *Config
	configDir      string
)

// Initialize initializes the configuration system
func Initialize(configFile, profileName string) error {
	// Set up config directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	configDir = filepath.Join(home, ".aphelion")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Load configuration
	config, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	configInstance = config
	
	// Set current profile
	if profileName != "" {
		configInstance.CurrentProfile = profileName
	}
	
	// Ensure default profile exists
	if configInstance.Profiles == nil {
		configInstance.Profiles = make(map[string]Profile)
	}
	
	if _, exists := configInstance.Profiles["default"]; !exists {
		configInstance.Profiles["default"] = Profile{
			Name:     "default",
			Endpoint: viper.GetString("endpoint"),
			Auth: Auth{
				Type:         "auth0",
				Domain:       "",
				ClientID:     "",
				Audience:     "",
				RedirectURI:  "http://localhost:8765/callback",
				TokenStorage: "keyring",
			},
		}
	}
	
	// If no current profile set, use default
	if configInstance.CurrentProfile == "" {
		configInstance.CurrentProfile = "default"
	}
	
	return nil
}

// loadConfig loads configuration from file or creates default
func loadConfig(configFile string) (*Config, error) {
	var config Config
	
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create default
			config = Config{
				CurrentProfile: "default",
				Profiles:       make(map[string]Profile),
			}
			return &config, nil
		}
		return nil, err
	}
	
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &config, nil
}

// Get returns the current configuration
func Get() *Config {
	return configInstance
}

// GetCurrentProfile returns the current profile
func GetCurrentProfile() Profile {
	if configInstance == nil {
		return Profile{}
	}
	
	profile, exists := configInstance.Profiles[configInstance.CurrentProfile]
	if !exists {
		return Profile{}
	}
	
	return profile
}

// SetProfile sets a profile configuration
func SetProfile(name string, profile Profile) error {
	if configInstance == nil {
		return fmt.Errorf("configuration not initialized")
	}
	
	profile.Name = name
	configInstance.Profiles[name] = profile
	
	return Save()
}

// DeleteProfile deletes a profile
func DeleteProfile(name string) error {
	if configInstance == nil {
		return fmt.Errorf("configuration not initialized")
	}
	
	if name == "default" {
		return fmt.Errorf("cannot delete default profile")
	}
	
	if _, exists := configInstance.Profiles[name]; !exists {
		return fmt.Errorf("profile %s does not exist", name)
	}
	
	delete(configInstance.Profiles, name)
	
	// If we deleted the current profile, switch to default
	if configInstance.CurrentProfile == name {
		configInstance.CurrentProfile = "default"
	}
	
	return Save()
}

// SwitchProfile switches to a different profile
func SwitchProfile(name string) error {
	if configInstance == nil {
		return fmt.Errorf("configuration not initialized")
	}
	
	if _, exists := configInstance.Profiles[name]; !exists {
		return fmt.Errorf("profile %s does not exist", name)
	}
	
	configInstance.CurrentProfile = name
	
	return Save()
}

// ListProfiles returns all profile names
func ListProfiles() []string {
	if configInstance == nil {
		return []string{}
	}
	
	var profiles []string
	for name := range configInstance.Profiles {
		profiles = append(profiles, name)
	}
	
	return profiles
}

// Save saves the configuration to file
func Save() error {
	if configInstance == nil {
		return fmt.Errorf("configuration not initialized")
	}
	
	configFile := filepath.Join(configDir, "config.yaml")
	
	viper.Set("current_profile", configInstance.CurrentProfile)
	viper.Set("profiles", configInstance.Profiles)
	
	return viper.WriteConfigAs(configFile)
}

// GetConfigDir returns the configuration directory
func GetConfigDir() string {
	return configDir
}