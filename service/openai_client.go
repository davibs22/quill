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

type branchNameResponse struct {
	BranchName string `json:"branchName"`
}

func NewOpenAIClient(model string) *OpenAIClient {
	apiKey := viper.GetString("preferences.openai.apikey")
	cli := openai.NewClient(option.WithAPIKey(apiKey))

	if apiKey == "" {
		logger.InitLogger("pretty")
		logger.L().Error("OpenAI API key not configured. Set preferences.openai.apikey in config or set --apikey flag.")
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

func (o *OpenAIClient) GenerateBranchName(workItem string, workItemId string, workItemTitle string) (string, error) {
	prompt := fmt.Sprintf(`You are an assistant that generates concise Azure DevOps branch names.
Given a Work Item description and its ID, return only a branch name following this pattern:

    <Type>/<WorkItemID>-<short_description_in_snake_case>

Rules:
1. Type must be "Feat" if the Work Item introduces a new feature, or "Fix" if it is a bug fix.
2. Use underscores ("_") to separate words in <short_description_in_snake_case>.
3. The name must be short, in English, and summarize the purpose of the task.
4. Do not include any additional details or formatting: only return the branch string.
5. The branch name must always be in English. This is an important and inviolable rule that must be respected.

Example:
- WorkItem_ID: 859652
- Description: "Create a new screen for user login validation"
Expected response:
- Feat/859652-create_new_screen

Now, using the Work Item below, generate only the branch name in English (without any additional formatting):

WorkItem_ID: %s
WorkItem_Title: %s
WorkItem_Description: %s
`, workItemId, workItemTitle, workItem)

	resp, err := o.client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Model:       o.model,
			Temperature: openai.Float(0.2),
		},
	)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Could not establish connection with LLM provider. Please verify if your API key is correct and if you are connected to the internet.")
		os.Exit(1)
	}

	return resp.Choices[0].Message.Content, nil
}
