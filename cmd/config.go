package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	logger "github.com/kubescape/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName       = "quill"
	configFile    = "config.yaml"
	linuxConfig   = ".config/quill"
	windowsConfig = "AppData/Roaming/quill"
	macConfig     = "Library/Application Support/quill"
)

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "linux":
		return filepath.Join(homeDir, linuxConfig), nil
	case "windows":
		return filepath.Join(homeDir, windowsConfig), nil
	case "darwin":
		return filepath.Join(homeDir, macConfig), nil
	default:
		return filepath.Join(homeDir, "."+appName), nil
	}
}

func initializeConfig(configPath string) error {
	if err := os.MkdirAll(configPath, 0755); err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Failed to create config directory.")
		os.Exit(1)
	}

	viper.SetConfigType("yaml")
	viper.Set("preferences", map[string]interface{}{
		"providerDefault": "",
		"openai": map[string]interface{}{
			"apiKey": "",
			"model":  "",
		},
		"llama": map[string]interface{}{
			"apiUrl": "",
			"model":  "",
		},
	})

	configFile := filepath.Join(configPath, "config.yaml")
	if err := viper.WriteConfigAs(configFile); err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Failed to write config file.")
		os.Exit(1)
	}

	logger.InitLogger("pretty")
	logger.L().Success(fmt.Sprintf("Created default config at: %s", configFile))
	logger.L().Info("Please edit the config file to set your API keys and preferences.")
	return nil
}

func editConfig(configPath string) error {
	fmt.Println("Configuration editor will be implemented here.")
	return nil
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Quill",
	Long:  "Configure Quill settings and preferences",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := getConfigPath()
		if err != nil {
			logger.InitLogger("pretty")
			logger.L().Error("Failed to get config path.")
			os.Exit(1)
		}
		if _, err := os.Stat(filepath.Join(configPath, configFile)); os.IsNotExist(err) {
			logger.InitLogger("pretty")
			logger.L().Warning(fmt.Sprintf("Config file not found at: %s", configPath))
			return initializeConfig(configPath)
		}

		logger.InitLogger("pretty")
		logger.L().Info(fmt.Sprintf("Using config file at: %s", configPath))
		return editConfig(configPath)
	},
}
