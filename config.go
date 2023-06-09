package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Configuration Options
type Config struct {
	// The base URL of the qBittorrent instance.
	BaseURL string `mapstructure:"base_url"`

	// The HTTP port number to use for the server.
	HTTPPort int `mapstructure:"http_port"`

	// The HTTPS port number to use for the server. Optional.
	HTTPSPort int `mapstructure:"https_port"`

	// The path to the TLS certificate file. Optional.
	TLSCert string `mapstructure:"tls_cert"`

	// The path to the TLS private key file. Optional.
	TLSKey string `mapstructure:"tls_key"`

	// CORS settings
	CORSAllowOrigins     []string `mapstructure:"cors_allow_origins"`
	CORSAllowMethods     []string `mapstructure:"cors_allow_methods"`
	CORSAllowHeaders     []string `mapstructure:"cors_allow_headers"`
	CORSExposeHeaders    []string `mapstructure:"cors_expose_headers"`
	CORSAllowCredentials bool     `mapstructure:"cors_allow_credentials"`

	// Debug Mode
	Debug bool `mapstructure:"debug"`
}

// Loads the configuration from various sources and unmarshals it into the Config struct.
func (cfg *Config) Load() error {
	setupViper()

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Failed to read config file: %v", err)
			return err
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("Failed to unmarshal config values: %v", err)
		return err
	}

	if err := validateConfig(cfg); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
		return err
	}

    log.Println("Config Loaded:")
	log.Printf("%+v\n", cfg)

	return nil
}

func (cfg *Config) GetHTTPSPortStr() string {
	return fmt.Sprintf(":%d", cfg.HTTPSPort)
}

func (cfg *Config) GetHTTPPortStr() string {
	return fmt.Sprintf(":%d", cfg.HTTPPort)
}

// Validates the configuration values.
func validateConfig(cfg *Config) error {
	if cfg.BaseURL == "" {
		return errors.New("BaseURL is required in the configuration")
	}
	return nil
}

// Configures Viper to search for configuration files and sets default values.
func setupViper() {
	setSearchPath()
	setDefaultValues()
	enableEnvCfg()
}

func setSearchPath() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(getXdgCfgPath())
}

func enableEnvCfg() {
	// Prefix will be uppercased automatically
	viper.SetEnvPrefix(ENV_PREFIX)
	// Enable Viper to read environment variables
	viper.AutomaticEnv()
}

func setDefaultValues() {
	viper.SetDefault("base_url", "")
	viper.SetDefault("http_port", 9487)
	viper.SetDefault("https_port", 0) // Set to 0 to disable HTTPS by default
	viper.SetDefault("tls_cert", "")
	viper.SetDefault("tls_key", "")
	viper.SetDefault("cors_allow_origins", []string{})
	viper.SetDefault("cors_allow_methods", []string{})
	viper.SetDefault("cors_allow_headers", []string{})
	viper.SetDefault("cors_expose_headers", []string{})
	viper.SetDefault("cors_allow_credentials", false)
	viper.SetDefault("debug", false)
}

func getXdgCfgPath() string {
	cfgPath := os.Getenv("XDG_CONFIG_HOME")
	if cfgPath == "" {
		cfgPath = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath = filepath.Join(cfgPath, APP_NAME)
	return cfgPath
}
