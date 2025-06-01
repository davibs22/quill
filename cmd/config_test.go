package cmd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfigFlags(t *testing.T) {
	tempDir := t.TempDir()
	err := initializeConfig(tempDir)
	assert.NoError(t, err)

	configFilePath := filepath.Join(tempDir, configFile)
	viper.SetConfigFile(configFilePath)
	err = viper.ReadInConfig()
	assert.NoError(t, err)

	// Create a new command for testing
	cmd := &cobra.Command{
		RunE: configCmd.RunE,
	}
	cmd.Flags().String("set-openai-api-key", "", "Set OpenAI API Key")
	cmd.Flags().String("set-openai-model", "", "Set OpenAI model")
	cmd.Flags().String("set-ollama-model", "", "Set Ollama model")
	cmd.Flags().String("set-ollama-api-url", "", "Set Ollama API URL")
	cmd.Flags().String("set-provider-default", "", "Set provider default")

	// Test setting OpenAI API Key
	cmd.Flags().Set("set-openai-api-key", "sk-1234567890")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "sk-1234567890", viper.GetString("preferences.openai.apiKey"))

	// Test setting OpenAI model
	cmd.Flags().Set("set-openai-model", "gpt-4o-mini")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "gpt-4o-mini", viper.GetString("preferences.openai.model"))

	// Test setting Ollama model
	cmd.Flags().Set("set-ollama-model", "llama2")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "llama2", viper.GetString("preferences.ollama.model"))

	// Test setting Ollama API URL
	cmd.Flags().Set("set-ollama-api-url", "http://localhost:11434/api/generate")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:11434/api/generate", viper.GetString("preferences.ollama.apiUrl"))

	// Test setting provider default
	cmd.Flags().Set("set-provider-default", "ollama")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "ollama", viper.GetString("preferences.providerDefault"))
}
