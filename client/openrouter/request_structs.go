package openrouter

// openRouterResponse defines the structure for the entire API response.
type openRouterResponse struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
	Usage   usage    `json:"usage"`
}

// choice represents a single choice in the response.
type choice struct {
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// message contains the role and content of the message from the model.
type message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	Reasoning string     `json:"reasoning"`
	ToolCalls []toolCall `json:"tool_calls,omitempty"`
}

type toolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"` // Using interface{} for flexibility
}

type tool struct {
	Type     string       `json:"type"`
	Function toolFunction `json:"function"`
}

// For the response
type functionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // Arguments are a JSON string
}

type toolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function functionCall `json:"function"`
}
