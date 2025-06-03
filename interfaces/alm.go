package interfaces

type ALMClient interface {
	WorkItemDetails(workItemId string) (string, string, error)
}
