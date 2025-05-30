package interfaces

type LLMClient interface {
	GenerateCommitMessage(diff string) (string, error)
}
