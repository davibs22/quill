package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davibs22/quill/service"
	logger "github.com/kubescape/go-logger"
	"github.com/spf13/cobra"
)

type branchNameResponse struct {
	BranchName string `json:"branchName"`
}

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

		client := service.NewLlamaClient(model)
		if err != nil {
			logger.InitLogger("pretty")
			logger.L().Error("Could not create client Ollama.")
			os.Exit(1)
		}

		branchName, err := client.GenerateBranchName(workItemDetails, ticket, workItemTitle)
		if err != nil {
			logger.InitLogger("pretty")
			logger.L().Error("Could not generate branch name.")
			os.Exit(1)
		}

		var bn branchNameResponse
		if err := json.Unmarshal([]byte(branchName), &bn); err != nil {
			logger.InitLogger("pretty")
			logger.L().Error("Error unmarshalling response body from Azure DevOps API.")
			os.Exit(1)
		}

		fmt.Printf(bn.BranchName)
		return nil
	},
}
