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

	cmd := &cobra.Command{
		RunE: configCmd.RunE,
	}
	cmd.Flags().String("set-openai-api-key", "", "Set OpenAI API Key")
	cmd.Flags().String("set-openai-model", "", "Set OpenAI model")
	cmd.Flags().String("set-ollama-model", "", "Set Ollama model")
	cmd.Flags().String("set-ollama-api-url", "", "Set Ollama API URL")
	cmd.Flags().String("set-provider-default", "", "Set provider default")

	cmd.Flags().Set("set-openai-api-key", "sk-1234567890")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "sk-1234567890", viper.GetString("preferences.openai.apiKey"))

	cmd.Flags().Set("set-openai-model", "gpt-4o-mini")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "gpt-4o-mini", viper.GetString("preferences.openai.model"))

	cmd.Flags().Set("set-ollama-model", "llama3.2")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "llama3.2", viper.GetString("preferences.ollama.model"))

	err = cmd.Execute()
	assert.NoError(t, err)

	cmd.Flags().Set("set-provider-default", "ollama")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "ollama", viper.GetString("preferences.providerDefault"))
}
func TestConfigEdgeCases(t *testing.T) {
	tempDir := t.TempDir()
	err := initializeConfig(tempDir)
	assert.NoError(t, err)

	configFilePath := filepath.Join(tempDir, configFile)
	viper.SetConfigFile(configFilePath)
	err = viper.ReadInConfig()
	assert.NoError(t, err)

	cmd := &cobra.Command{
		RunE: configCmd.RunE,
	}
	cmd.Flags().String("set-openai-api-key", "", "Set OpenAI API Key")
	cmd.Flags().String("set-openai-model", "", "Set OpenAI model")
	cmd.Flags().String("set-ollama-model", "", "Set Ollama model")
	cmd.Flags().String("set-ollama-api-url", "", "Set Ollama API URL")
	cmd.Flags().String("set-provider-default", "", "Set provider default")
	cmd.Flags().String("set-gemini-model", "", "Set Gemini model")

	cmd.Flags().Set("set-gemini-model", "gemini-2.0-flash")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "gemini-2.0-flash", viper.GetString("preferences.gemini.model"))

	cmd.Flags().Set("set-provider-default", "gemini")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "gemini", viper.GetString("preferences.providerDefault"))

	cmd.Flags().Set("set-openai-model", "gpt-4")
	cmd.Flags().Set("set-ollama-model", "llama3.2")
	cmd.Flags().Set("set-gemini-model", "gemini-2.0-flash")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "gpt-4", viper.GetString("preferences.openai.model"))
	assert.Equal(t, "llama3.2", viper.GetString("preferences.ollama.model"))
	assert.Equal(t, "gemini-2.0-flash", viper.GetString("preferences.gemini.model"))

	cmd.Flags().Set("set-openai-model", "")
	cmd.Flags().Set("set-ollama-model", "")
	cmd.Flags().Set("set-gemini-model", "")
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "gpt-4", viper.GetString("preferences.openai.model"))
	assert.Equal(t, "llama3.2", viper.GetString("preferences.ollama.model"))
	assert.Equal(t, "gemini-2.0-flash", viper.GetString("preferences.gemini.model"))
}
