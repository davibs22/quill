package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	logger "github.com/kubescape/go-logger"
	"github.com/spf13/viper"
)

type AzureClient struct {
	workItemId string
	apiKey     string
	company    string
	project    string
}

type workItemResponse struct {
	Fields struct {
		Title              string `json:"System.Title"`
		Description        string `json:"System.Description"`
		AcceptanceCriteria string `json:"Microsoft.VSTS.Common.AcceptanceCriteria"`
	} `json:"fields"`
}

func NewAzureClient(workItemId string) *AzureClient {
	apiKey := viper.GetString("preferences.azure.apikey")
	company := viper.GetString("preferences.azure.company")
	project := viper.GetString("preferences.azure.project")

	if apiKey == "" {
		logger.InitLogger("pretty")
		logger.L().Error("Azure API key not configured. Set preferences.azure.apikey in config or set --azure-apikey flag.")
		os.Exit(1)
	}

	if company == "" {
		logger.InitLogger("pretty")
		logger.L().Error("Azure company not configured. Set preferences.azure.company in config or set --azure-company flag.")
		os.Exit(1)
	}

	if project == "" {
		logger.InitLogger("pretty")
		logger.L().Error("Azure project not configured. Set preferences.azure.project in config or set --azure-project flag.")
		os.Exit(1)
	}

	return &AzureClient{workItemId: workItemId, apiKey: apiKey, company: company, project: project}
}

func (a *AzureClient) WorkItemDetails(workItemId string) (string, string, error) {
	url := fmt.Sprintf(
		"https://dev.azure.com/%s/%s/_apis/wit/workitems/%s?api-version=7.1-preview.3",
		a.company, a.project, workItemId,
	)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Error creating request to Azure DevOps API.")
		os.Exit(1)
	}

	pat := a.apiKey
	if pat == "" {
		logger.InitLogger("pretty")
		logger.L().Error("Azure API key not configured. Set preferences.azure.apikey in config or set --azure-apikey flag.")
		os.Exit(1)
	}

	toEncode := fmt.Sprintf(":%s", pat)

	encoded := base64.StdEncoding.EncodeToString([]byte(toEncode))

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+encoded)

	res, err := client.Do(req)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Error making request to Azure DevOps API.")
		os.Exit(1)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.InitLogger("pretty")
		logger.L().Error("Error getting work item details from Azure DevOps API.")
		os.Exit(1)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Error reading response body from Azure DevOps API.")
		os.Exit(1)
	}

	var wi workItemResponse
	if err := json.Unmarshal(bodyBytes, &wi); err != nil {
		logger.InitLogger("pretty")
		logger.L().Error("Error unmarshalling response body from Azure DevOps API.")
		os.Exit(1)
	}

	title := wi.Fields.Title
	description := wi.Fields.Description + "\n\n" + wi.Fields.AcceptanceCriteria

	return description, title, nil
}
