package service

type LLMClient interface {
	GenerateCommitMessage(diff string) (string, error)
}
