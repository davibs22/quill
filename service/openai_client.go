package service

import (
	"context"
	"fmt"
	"os"

	logger "github.com/kubescape/go-logger"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/spf13/viper"
)

type OpenAIClient struct {
	client openai.Client
	model  string
}

func NewOpenAIClient(model string) *OpenAIClient {
	apiKey := viper.GetString("preferences.openai.apiKey")
	cli := openai.NewClient(option.WithAPIKey(apiKey))

	if apiKey == "" {
		logger.InitLogger("pretty")
		logger.L().Error("OpenAI API key not configured. Set preferences.openai.apiKey in config or set --apikey flag.")
		os.Exit(1)
	}

	return &OpenAIClient{client: cli, model: model}
}

func (o *OpenAIClient) GenerateCommitMessage(diff string) (string, error) {
	prompt := fmt.Sprintf("Analyze the following code changes and generate a commit message following the Conventional Commits standard. The message should:\n\n1. Start with a conventional type (feat, fix, docs, style, refactor, perf, test, chore)\n2. Include an optional scope in parentheses when applicable\n3. Have a concise description in imperative mood (e.g., \"change\" → \"change X to do Y\")\n4. The message language should be in English.\n\nReturn ONLY the commit message itself, without any additional explanations or commentary.\n\nChanges:\n\n%s", diff)

	resp, err := o.client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Model: o.model,
		},
	)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Could not establish connection with LLM provider. Please verify if your API key is correct and if you are connected to the internet.")
		os.Exit(1)
	}

	return resp.Choices[0].Message.Content, nil
}
