package interfaces

type LLMClient interface {
	GenerateCommitMessage(diff string) (string, error)
	GenerateBranchName(diff string) (string, error)
}
