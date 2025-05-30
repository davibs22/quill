package service

import (
	"context"
	"os"

	logger "github.com/kubescape/go-logger"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
	model  string
}

func NewOpenAIClient(model string) *OpenAIClient {
	apiKey := os.Getenv("OPENAI_API_KEY")
	cli := openai.NewClient(apiKey)
	return &OpenAIClient{client: cli, model: model}
}

func (o *OpenAIClient) GenerateCommitMessage(diff string) (string, error) {
	ctx := context.Background()
	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: o.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are a helpful assistant that crafts concise git commit messages."},
			{Role: "user", Content: "Generate a git commit message for the following diff:\n\n" + diff},
		},
	})
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Could not establish connection with OpenAI API. Please verify if your API key is correct and if you are connected to the internet.")
		os.Exit(1)
	}
	return resp.Choices[0].Message.Content, nil
}
