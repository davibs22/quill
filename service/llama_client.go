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

func NewLlamaClient() *LlamaClient {

	endpoint := viper.GetString("preferences.llama.apiUrl")
	model := viper.GetString("preferences.llama.model")

	if endpoint == "" {
		endpoint = os.Getenv("LLAMA_API_URL")
	}
	if model == "" {
		model = os.Getenv("LLAMA_MODEL")
	}

	if endpoint == "" {
		logger.InitLogger("pretty")
		logger.L().Error("Llama API URL não configurada. Defina preferences.llama.apiUrl no config ou LLAMA_API_URL no ambiente.")
		os.Exit(1)
	}
	if model == "" {
		logger.InitLogger("pretty")
		logger.L().Error("Llama model não configurado. Defina preferences.llama.model no config ou LLAMA_MODEL no ambiente.")
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
		logger.L().Error("Falha ao montar payload para Llama API.")
		os.Exit(1)
	}

	resp, err := http.Post(l.endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Não foi possível conectar na Llama API. Verifique o endpoint.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Resposta inválida da Llama API.")
		os.Exit(1)
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("JSON inválido ou campo 'response' não encontrado.")
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
		fmt.Printf("Generated commit message: %s\n", result.Response)
		logger.InitLogger("pretty")
		logger.L().Warning("The generated commit message doesn't follow Conventional Commits standards")
		os.Exit(1)
	}

	return result.Response, nil
}
