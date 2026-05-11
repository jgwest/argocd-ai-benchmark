package openrouter

import (
	"argocd-ai-benchmark/client"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var _ client.ClientInstance = &openRouterClientInstance{}

type OpenRouterClient struct {
	APIKey string
	Model  string
}

type openRouterClientInstance struct {
	client           *OpenRouterClient
	previousMessages []string
	// mainContext      types.EvalContext
}

var _ client.Client = &OpenRouterClient{}

// NewInstance implements types.Client.
func (orc *OpenRouterClient) NewInstance() client.ClientInstance {
	return &openRouterClientInstance{
		client:           orc,
		previousMessages: []string{},
		// mainContext:      mainContext,
	}
}

func (orc *openRouterClientInstance) SendMessage(messageParam string) (*client.ClientResponse, error) {

	url := "https://openrouter.ai/api/v1/chat/completions"

	messagesMap := []map[string]string{}

	for _, message := range orc.previousMessages {
		messagesMap = append(messagesMap, map[string]string{
			"role":    "user",
			"content": message,
		})
	}

	messagesMap = append(messagesMap, map[string]string{
		"role":    "user",
		"content": messageParam,
	})

	// Define the request body
	requestBody, err := json.Marshal(map[string]any{
		"model":    orc.client.Model,
		"messages": messagesMap,
		"usage": map[string]bool{
			"include": true,
		},
		// https://openrouter.ai/docs/use-cases/reasoning-tokens
		// "reasoning": map[string]any{
		// 	"effort":  "high",
		// 	"exclude": false,
		// 	"enabled": true,
		// },
	})

	if err != nil {
		return nil, fmt.Errorf("Error marshalling request body: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+orc.client.APIKey)

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	// Parse the JSON response into our structs
	var apiResponse openRouterResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		// Also print the raw body to help debug if parsing fails
		fmt.Println("Raw Response Body:", string(body))
		return nil, fmt.Errorf("Error unmarshalling response JSON: %v", err)
	}

	if len(apiResponse.Choices) != 1 {
		return nil, fmt.Errorf("Unexpected choices parameter: %v", string(body))
	}

	response := apiResponse.Choices[0].Message

	var reasoning string
	// if orc.mainContext.Configuration.PrintReasoning {

	output := ""

	for _, toolCall := range response.ToolCalls {
		output += fmt.Sprintln("-", toolCall)
	}

	if response.Reasoning != "" {
		output += fmt.Sprintln("-", response.Reasoning)
	}

	reasoning += output
	// }

	orc.previousMessages = append(orc.previousMessages, response.Content)

	res := &client.ClientResponse{
		ResponseContent:  response.Content,
		ReasoningContent: reasoning,
		Usage:            apiResponse.Usage,
	}

	return res, nil
}

// usage provides the token and cost details when 'usage accounting' is enabled.
type usage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	Cost             float64 `json:"cost"`
}

var _ client.ResponseUsage = usage{}

func (u usage) Aggregate(one client.ResponseUsage) client.ResponseUsage {

	obj, ok := (one).(usage)
	if !ok {
		log.Fatal("Unexpected casting error")
	}

	// return MergeUsage(obj, u)

	return usage{
		PromptTokens:     obj.PromptTokens + u.PromptTokens,
		CompletionTokens: obj.CompletionTokens + u.CompletionTokens,
		TotalTokens:      obj.TotalTokens + u.TotalTokens,
		Cost:             obj.Cost + u.Cost,
	}
}

func (u usage) GenerateReport() string {
	var res string
	res += fmt.Sprintf("- Prompt Tokens: %d\n", u.PromptTokens)
	res += fmt.Sprintf("- Completion Tokens: %d\n", u.CompletionTokens)
	res += fmt.Sprintf("- Total Tokens: %d\n", u.TotalTokens)
	res += fmt.Sprintf("- Cost: $%.8f\n", u.Cost)
	return res
}

// func MergeUsage(u1 Usage, u2 Usage) Usage {
// 	return Usage{
// 		PromptTokens:     u1.PromptTokens + u2.PromptTokens,
// 		CompletionTokens: u1.CompletionTokens + u2.CompletionTokens,
// 		TotalTokens:      u1.TotalTokens + u2.TotalTokens,
// 		Cost:             u1.Cost + u2.Cost,
// 	}
// }
