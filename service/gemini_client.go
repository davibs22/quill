package service

import (
	"context"
	"fmt"
	"os"

	logger "github.com/kubescape/go-logger"
	"github.com/spf13/viper"
	genai "google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(model string) *GeminiClient {
	apiKey := viper.GetString("preferences.gemini.apikey")
	if apiKey == "" {
		logger.InitLogger("pretty")
		logger.L().Error("Gemini API key not configured. Set preferences.gemini.apikey in config.")
		os.Exit(1)
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Failed to create Gemini client.")
		os.Exit(1)
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}
}

func (g *GeminiClient) GenerateCommitMessage(diff string) (string, error) {
	prompt := fmt.Sprintf(
		"Analyze the following code changes and generate a commit message following the Conventional Commits standard. "+
			"The message should:\n\n"+
			"1. Start with a conventional type (feat, fix, docs, style, refactor, perf, test, chore)\n"+
			"2. Include an optional scope in parentheses when applicable\n"+
			"3. Have a concise description in imperative mood (e.g., \"change\" → \"change X to do Y\")\n"+
			"4. The message language should be in English.\n\n"+
			"Return ONLY the commit message itself, without any additional explanations or commentary.\n\n"+
			"Changes:\n\n%s", diff,
	)

	resp, err := g.client.Models.GenerateContent(
		context.Background(),
		g.model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Could not establish connection with Gemini API. Please verify if your API key is correct.")
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned from Gemini API")
	}

	return resp.Text(), nil
}
