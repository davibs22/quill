package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/davibs22/quill/interfaces"
	"github.com/davibs22/quill/service"
	logger "github.com/kubescape/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "0.1.0"

var (
	provider string
	model    string
	cfgPath  string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&provider, "provider", "openai", "LLM provider: openai|ollama")
	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", "gpt-4o", "Model name (used in OpenAI)")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show version information")
	configCmd.Flags().String("set-openai-api-key", "", "Sets the API key for OpenAI")
	configCmd.Flags().String("set-openai-model", "", "Sets the model for OpenAI (e.g.: gpt-4o-mini)")
	configCmd.Flags().String("set-ollama-model", "", "Sets the model for Ollama (e.g.: llama3.2:latest)")
	configCmd.Flags().String("set-ollama-api-url", "", "Sets the API url for Ollama")
	configCmd.Flags().String("set-provider-default", "", "Sets the default provider (e.g.: openai|ollama)")
	rootCmd.AddCommand(configCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		showVersion, err := cmd.Flags().GetBool("version")
		if err != nil {
			return err
		}
		if showVersion {
			fmt.Printf("Quill version %s\n", version)
			os.Exit(0)
		}
		return nil
	}
}

var rootCmd = &cobra.Command{
	Use:   "quill",
	Short: "Generate AI-powered Git commit messages",
	RunE:  runRoot,
}

func initConfig() {
	path, err := getConfigPath()
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Failed to find config path.")
		os.Exit(1)
	}
	cfgPath = fmt.Sprintf("%s/config.yaml", path)
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		initializeConfig(path)
		logger.InitLogger("pretty")
		logger.L().Error("Enter the necessary flags for Quill to work.")
		logger.L().Info("Execute \"quill --help\" for more information.")
		os.Exit(1)
	}
	viper.SetConfigFile(cfgPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return
	}
	if !rootCmd.PersistentFlags().Changed("provider") {
		p := viper.GetString("preferences.providerDefault")
		if p != "" {
			switch p {
			case "openai", "ollama":
				provider = p
			default:
				logger.InitLogger("pretty")
				logger.L().Error(fmt.Sprintf("Invalid provider in config: %s", p))
				os.Exit(1)
			}
		}
	}
	if !rootCmd.PersistentFlags().Changed("model") {
		if provider == "openai" {
			m := viper.GetString("preferences.openai.model")
			if m != "" {
				model = m
			}
		}
		if provider == "ollama" {
			m := viper.GetString("preferences.ollama.model")
			if m != "" {
				model = m
			}
		}
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	diffOut, err := exec.Command("git", "diff", "--staged").CombinedOutput()
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Could not establish connection with Git. Please verify if you are in a Git repository.")
		os.Exit(1)
	}
	if len(diffOut) == 0 {
		logger.InitLogger("pretty")
		logger.L().Error("No changes found in Git.")
		os.Exit(1)
	}
	if len(diffOut) > 10000 {
		logger.InitLogger("pretty")
		logger.L().Error("Diff is too large. Please reduce the number of changes.")
		os.Exit(1)
	}

	var client interfaces.LLMClient
	switch provider {
	case "openai":
		client = service.NewOpenAIClient(model)
	case "ollama":
		client = service.NewLlamaClient(model)
	default:
		return fmt.Errorf("invalid provider: %s", provider)
	}

	msg, err := client.GenerateCommitMessage(string(diffOut))
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Could not establish connection with LLM provider. Please verify if your API key is correct and if you are connected to the internet.")
		os.Exit(1)
	}

	fmt.Printf("%s", msg)
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
