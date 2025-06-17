package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	APIUrl      string    `yaml:"api_url" mapstructure:"api_url"`
	AccessToken string    `yaml:"access_token" mapstructure:"access_token"`
	UserID      string    `yaml:"user_id" mapstructure:"user_id"`
	Email       string    `yaml:"email" mapstructure:"email"`
	Username    string    `yaml:"username" mapstructure:"username"`
	LastLogin   time.Time `yaml:"last_login" mapstructure:"last_login"`
	Output      string    `yaml:"output" mapstructure:"output"`
}

var globalConfig *Config

func InitConfig() {
	viper.SetDefault("api_url", "https://api.aphelion.exmplr.ai")
	viper.SetDefault("output", "table")

	globalConfig = &Config{
		APIUrl: viper.GetString("api_url"),
		Output: viper.GetString("output"),
	}

	if viper.IsSet("access_token") {
		globalConfig.AccessToken = viper.GetString("access_token")
	}
	if viper.IsSet("user_id") {
		globalConfig.UserID = viper.GetString("user_id")
	}
	if viper.IsSet("email") {
		globalConfig.Email = viper.GetString("email")
	}
	if viper.IsSet("username") {
		globalConfig.Username = viper.GetString("username")
	}
	if viper.IsSet("last_login") {
		globalConfig.LastLogin = viper.GetTime("last_login")
	}
}

func GetConfig() *Config {
	if globalConfig == nil {
		InitConfig()
	}
	return globalConfig
}

func SaveConfig() error {
	if globalConfig == nil {
		return fmt.Errorf("config not initialized")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".aphelion")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")
	
	data, err := yaml.Marshal(globalConfig)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	viper.Set("access_token", globalConfig.AccessToken)
	viper.Set("user_id", globalConfig.UserID)
	viper.Set("email", globalConfig.Email)
	viper.Set("username", globalConfig.Username)
	viper.Set("last_login", globalConfig.LastLogin)

	return nil
}

func SetAuth(token, userID, email, username string) error {
	config := GetConfig()
	config.AccessToken = token
	config.UserID = userID
	config.Email = email
	config.Username = username
	config.LastLogin = time.Now()
	
	return SaveConfig()
}

func ClearAuth() error {
	config := GetConfig()
	config.AccessToken = ""
	config.UserID = ""
	config.Email = ""
	config.Username = ""
	config.LastLogin = time.Time{}
	
	return SaveConfig()
}

func IsAuthenticated() bool {
	config := GetConfig()
	return config.AccessToken != ""
}

func GetAPIUrl() string {
	return GetConfig().APIUrl
}

func GetAccessToken() string {
	return GetConfig().AccessToken
}

func GetUserID() string {
	return GetConfig().UserID
}

func GetOutputFormat() string {
	return GetConfig().Output
}