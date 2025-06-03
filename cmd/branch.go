package cmd

import (
	"fmt"
	"os"

	"github.com/davibs22/quill/interfaces"
	"github.com/davibs22/quill/service"
	logger "github.com/kubescape/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Generate a branch name based on the staged changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		alm, _ := cmd.Flags().GetString("alm")
		ticket, _ := cmd.Flags().GetString("ticket")

		if alm == "" {
			logger.InitLogger("pretty")
			logger.L().Error("ALM not specified.")
			os.Exit(1)
		}

		if alm != "azure" {
			logger.InitLogger("pretty")
			logger.L().Error("ALM not supported.")
			os.Exit(1)
		}

		if ticket == "" {
			logger.InitLogger("pretty")
			logger.L().Error("Ticket ID not specified.")
			os.Exit(1)
		}

		azureClient := service.NewAzureClient(ticket)

		workItemDetails, workItemTitle, err := azureClient.WorkItemDetails(ticket)
		if err != nil {
			logger.InitLogger("pretty")
			logger.L().Error("Could not get work item details.")
			os.Exit(1)
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

		var client interfaces.LLMClient
		switch provider {
		case "ollama":
			client = service.NewLlamaClient(model)
			if err != nil {
				logger.InitLogger("pretty")
				logger.L().Error("Could not create client Ollama.")
				os.Exit(1)
			}
		case "openai":
			client = service.NewOpenAIClient(model)
			if err != nil {
				logger.InitLogger("pretty")
				logger.L().Error("Could not create OpenAI client.")
				os.Exit(1)
			}
		default:
			logger.InitLogger("pretty")
			logger.L().Error("Invalid provider specified. Use 'ollama' or 'openai'.")
			os.Exit(1)
		}
		branchName, err := client.GenerateBranchName(workItemDetails, ticket, workItemTitle)
		if err != nil {
			logger.InitLogger("pretty")
			logger.L().Error("Could not generate branch name.")
			os.Exit(1)
		}

		fmt.Println(branchName)
		return nil
	},
}
