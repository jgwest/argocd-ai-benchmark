package client

type ResponseUsage interface {
	Aggregate(one ResponseUsage) ResponseUsage
	GenerateReport() string
}

type Client interface {
	NewInstance() ClientInstance
}

type ClientInstance interface {
	SendMessage(message string) (*ClientResponse, error)
}

type ClientResponse struct {
	ResponseContent  string
	ReasoningContent string
	Usage            ResponseUsage
}
