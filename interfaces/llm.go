package interfaces

type LLMClient interface {
	GenerateCommitMessage(diff string) (string, error)
	GenerateBranchName(workItem string, workItemId string, workItemTitle string) (string, error)
}
