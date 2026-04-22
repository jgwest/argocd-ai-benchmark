package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func chatWithHistory(modelName string, messages []string) (*Message, *Usage, error) {

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	url := "https://openrouter.ai/api/v1/chat/completions"

	messagesMap := []map[string]string{}

	for _, message := range messages {
		messagesMap = append(messagesMap, map[string]string{
			"role":    "user",
			"content": message,
		})
	}

	// Define the request body
	requestBody, err := json.Marshal(map[string]any{
		"model":    modelName,
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
		return nil, nil, fmt.Errorf("Error marshalling request body: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Error reading response body: %v", err)
	}

	// Parse the JSON response into our structs
	var apiResponse OpenRouterResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		// Also print the raw body to help debug if parsing fails
		fmt.Println("Raw Response Body:", string(body))
		return nil, nil, fmt.Errorf("Error unmarshalling response JSON: %v", err)
	}

	if len(apiResponse.Choices) != 1 {
		return nil, nil, fmt.Errorf("Unexpected choices parameter: %v", string(body))
	}

	return &apiResponse.Choices[0].Message, &apiResponse.Usage, nil
}

// OpenRouterResponse defines the structure for the entire API response.
type OpenRouterResponse struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a single choice in the response.
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Message contains the role and content of the message from the model.
type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	Reasoning string     `json:"reasoning"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"` // Using interface{} for flexibility
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// For the response
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // Arguments are a JSON string
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// Usage provides the token and cost details when 'usage accounting' is enabled.
type Usage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	Cost             float64 `json:"cost"`
}

func MergeUsage(u1 Usage, u2 Usage) Usage {
	return Usage{
		PromptTokens:     u1.PromptTokens + u2.PromptTokens,
		CompletionTokens: u1.CompletionTokens + u2.CompletionTokens,
		TotalTokens:      u1.TotalTokens + u2.TotalTokens,
		Cost:             u1.Cost + u2.Cost,
	}
}
