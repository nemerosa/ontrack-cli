package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	configFileName = ".ontrack-cli-config.yaml"
)

// Root configuration
type RootConfig struct {
	// Default configuration name
	Selected string
	// List of configurations
	Configurations []Config
}

// Configuration content
type Config struct {
	// Name of the configuration
	Name string
	// URL of the remote server
	URL string
	// Username for the remote server (when using basic authentication)
	Username string
	// Password for the remote server (when using basic authentication)
	Password string
	// Token for the remote server (when using token-based authentication)
	Token string
	// Is this configuration disabled?
	Disabled bool
}

// Gets the current configuration
func GetSelectedConfiguration() (*Config, error) {
	root := ReadRootConfiguration()
	selected := root.Selected
	if selected != "" {
		for _, item := range root.Configurations {
			if item.Name == selected {
				return &item, nil
			}
		}
		return nil, fmt.Errorf("No configuration named %s", selected)
	}
	return nil, errors.New("No current configuration")
}

// Reads the configuration
func ReadRootConfiguration() *RootConfig, error {
	var root RootConfig
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// If the config file does not exist, returns an empty root config
	if _, err := os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			return &root, nil
		}
	}

	reader, _ := os.Open(configFilePath)
	buf, _ := ioutil.ReadAll(reader)
	yaml.Unmarshal(buf, &root)
	return &root, nil
}

// Adds a new configuration and set as default
func AddConfiguration(config Config, override bool) error {
	root := ReadRootConfiguration()
	configurations := root.Configurations
	existing := false
	// Check if the configuration name already exists
	for index, item := range configurations {
		if item.Name == config.Name {
			if override {
				configurations[index] = config
				existing = true
			} else {
				return fmt.Errorf("Configuration with name %s already exists", config.Name)
			}
		}
	}
	// Default selected configuration is the added one
	// Adds the configuration to the list if not existing already

	if !existing {
		configurations = append(root.Configurations, config)
	}
	newRoot := RootConfig{
		Selected:       config.Name,
		Configurations: configurations,
	}
	// Saves the root configuration back
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	buf, _ := yaml.Marshal(newRoot)
	_, _ = os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	_ = ioutil.WriteFile(configFilePath, buf, 0600)

	// OK
	return nil
}

// Finds an existing configuration
func findConfigurationByName(root *RootConfig, name string) *Config {
	for _, item := range root.Configurations {
		if item.Name == name {
			return &item
		}
	}
	return nil
}

// Replacing an existing configuration
func replaceConfigurationByName(root *RootConfig, config *Config) {
	for index, item := range root.Configurations {
		if item.Name == config.Name {
			root.Configurations[index] = *config
		}
	}
}

// Sets the new selected configuration
func SetSelectedConfiguration(name string) error {
	root := ReadRootConfiguration()
	existing := findConfigurationByName(root, name)
	if existing == nil {
		return fmt.Errorf("Configuration with name %s does not exist", name)
	}
	newRoot := RootConfig{
		Selected:       name,
		Configurations: root.Configurations,
	}
	// Saves the root configuration back
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	buf, _ := yaml.Marshal(newRoot)
	_, _ = os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	_ = ioutil.WriteFile(configFilePath, buf, 0600)

	// OK
	return nil
}

// Disables or enabled a configuration
func SetConfigurationState(name string, disabled bool) error {
	root := ReadRootConfiguration()
	existing := findConfigurationByName(root, name)
	if existing == nil {
		return fmt.Errorf("Configuration with name %s does not exist", name)
	}
	// Adjust the existing configuration
	existing.Disabled = disabled
	replaceConfigurationByName(root, existing)
	// Saves the root configuration back
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	buf, _ := yaml.Marshal(root)
	_, _ = os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	_ = ioutil.WriteFile(configFilePath, buf, 0600)

	// OK
	return nil
}

// Gets the path to the configuration file
func getConfigFilePath() string, error {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(path, configFileName)
	return configFilePath, nil
}
