package interfaces

type ALMClient interface {
	WorkItemDetails(workItemId string) (string, error)
}
