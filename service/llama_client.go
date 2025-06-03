package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	logger "github.com/kubescape/go-logger"
	"github.com/spf13/viper"
)

type LlamaClient struct {
	endpoint string
	model    string
}

func NewLlamaClient(model string) *LlamaClient {

	endpoint := viper.GetString("preferences.ollama.apiUrl")

	if model == "" {
		model = viper.GetString("preferences.ollama.model")
	}

	if endpoint == "" {
		logger.InitLogger("pretty")
		logger.L().Error("Ollama API URL not configured. Set preferences.ollama.apiUrl in config or set --model flag.")
		os.Exit(1)
	}
	if model == "" {
		logger.InitLogger("pretty")
		logger.L().Error("Ollama model not configured. Set preferences.ollama.model in config or set --model flag.")
		os.Exit(1)
	}

	return &LlamaClient{
		endpoint: endpoint,
		model:    model,
	}
}

func (l *LlamaClient) GenerateCommitMessage(diff string) (string, error) {
	payload := map[string]interface{}{
		"model":      l.model,
		"prompt":     "Analyze the following code changes and generate a commit message following the Conventional Commits standard. The message should:\n\n1. Start with a conventional type (feat, fix, docs, style, refactor, perf, test, chore)\n2. Include an optional scope in parentheses when applicable\n3. Have a concise description in imperative mood (e.g., \"change\" → \"change X to do Y\")\n4. The message language should be in English.\n\nReturn ONLY the commit message itself, without any additional explanations or commentary.\n\nChanges:\n\n" + diff,
		"max_tokens": 100,
		"options": map[string]interface{}{
			"temperature": 0.0,
		},
		"stream": false,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Failed to build payload for ollama API.")
		os.Exit(1)
	}

	resp, err := http.Post(l.endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Could not connect to ollama API. Check the endpoint.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Invalid response from ollama API.")
		os.Exit(1)
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Invalid JSON or 'response' field not found.")
		os.Exit(1)
	}

	validTypes := []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore"}
	valid := false

	for _, t := range validTypes {
		if len(result.Response) > len(t) && result.Response[:len(t)] == t {

			rest := result.Response[len(t):]

			if len(rest) > 0 && rest[0] == '(' {
				if endScope := bytes.IndexByte([]byte(rest), ')'); endScope != -1 && endScope+2 < len(rest) {
					if rest[endScope+1] == ':' && rest[endScope+2] == ' ' {
						valid = true
					}
				}
			} else if len(rest) > 1 && rest[0] == ':' && rest[1] == ' ' {

				valid = true
			}
			break
		}
	}

	if !valid {
		logger.InitLogger("pretty")
		logger.L().Warning("The generated commit message doesn't follow Conventional Commits standards.")
		os.Exit(1)
	}

	return result.Response, nil
}

func (o *LlamaClient) GenerateBranchName(workItem string, workItemId string, workItemTitle string) (string, error) {
	payload := map[string]interface{}{
		"model": o.model,
		"prompt": fmt.Sprintf(`You are an assistant that generates concise Azure DevOps branch names.
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
`, workItemId, workItemTitle, workItem),
		"max_tokens": 100,
		"options": map[string]interface{}{
			"temperature": 0.0,
		},
		"stream": false,
		"format": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"branchName": map[string]interface{}{
					"type": "string",
				},
			},
			"required": []string{
				"branchName",
			},
		},
	}

	fmt.Printf(`You are an assistant that generates concise Azure DevOps branch names.
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
	body, err := json.Marshal(payload)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Failed to build payload for ollama API.")
		os.Exit(1)
	}
	resp, err := http.Post(o.endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Could not connect to ollama API. Check the endpoint.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Invalid response from ollama API.")
		os.Exit(1)
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Invalid JSON or 'response' field not found.")
		os.Exit(1)
	}

	return result.Response, nil
}
