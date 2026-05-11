package gemini

import (
	"argocd-ai-benchmark/client"
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var _ client.ClientInstance = &GeminiClientInstance{}

type GeminiClient struct {
	APIKey string
	Model  string
}

type GeminiClientInstance struct {
	client           *GeminiClient
	previousMessages []string
	// mainContext      types.EvalContext
}

var _ client.Client = &GeminiClient{}

// NewInstance implements types.Client.
func (orc *GeminiClient) NewInstance() client.ClientInstance {
	return &GeminiClientInstance{
		client:           orc,
		previousMessages: []string{},
		// mainContext:      mainContext,
	}
}

func (orc *GeminiClientInstance) SendMessage(messageParam string) (*client.ClientResponse, error) {

	// Initialize the genAIClient.
	// option.WithAPIKey(apiKey) is the explicit way to set the key.
	// If you've set the GEMINI_API_KEY environment variable,
	// you can also just use genai.NewClient(ctx)
	genAIClient, err := genai.NewClient(context.Background(), option.WithAPIKey(orc.client.APIKey))
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}
	defer genAIClient.Close() // Good practice to close the client when done.

	// Choose the model to use.
	model := genAIClient.GenerativeModel(orc.client.Model)

	// The prompt you want to send.
	prompt := genai.Text(messageParam)

	// Send the prompt to the API.
	resp, err := model.GenerateContent(context.Background(), prompt)
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("received an empty response")
	}

	if len(resp.Candidates) != 1 {
		return nil, fmt.Errorf("received an unexpected number of candidates")
	}

	var reasoningText string

	candidate := resp.Candidates[0]

	var contentPart genai.Part
	if len(candidate.Content.Parts) != 1 {

		// If there are multiple parts, the first parts will be the reasoning text, and the final part will be the answer (for gemini only)
		allPartsAreText := true
		for x, part := range candidate.Content.Parts {
			if text, ok := part.(genai.Text); ok {
				if x != len(candidate.Content.Parts)-1 {
					reasoningText += string(text)
				}
			} else {
				allPartsAreText = false
				break
			}
		}
		if !allPartsAreText {
			return nil, fmt.Errorf("unexpected number of non-text content parts: %d", len(candidate.Content.Parts))
		}

		// We know that all the parts are text, so the previous parts are likely just thinking tokens.
		// The final part will be the answer
		contentPart = candidate.Content.Parts[len(candidate.Content.Parts)-1]
	} else {
		contentPart = candidate.Content.Parts[0]
	}

	if txt, ok := contentPart.(genai.Text); !ok {
		return nil, fmt.Errorf("unexpected response type")
	} else {
		res := &client.ClientResponse{
			ResponseContent:  string(txt),
			ReasoningContent: reasoningText,
			Usage: usage{
				promptTokenCount:    int(resp.UsageMetadata.PromptTokenCount),
				candidateTokenCount: int(resp.UsageMetadata.CandidatesTokenCount),
				outputTokens:        int(candidate.TokenCount),
			},
		}

		return res, nil

	}

}

type usage struct {
	promptTokenCount    int
	candidateTokenCount int
	outputTokens        int
}

var _ client.ResponseUsage = usage{}

func (u usage) Aggregate(one client.ResponseUsage) client.ResponseUsage {

	obj, ok := (one).(usage)
	if !ok {
		log.Fatal("Unexpected casting error")
	}

	return usage{
		promptTokenCount:    u.promptTokenCount + obj.promptTokenCount,
		candidateTokenCount: u.candidateTokenCount + obj.candidateTokenCount,
		outputTokens:        u.outputTokens + obj.outputTokens,
	}

}

func (u usage) GenerateReport() string {
	var res string
	res += fmt.Sprintf("- Prompt tokens: %d\n", u.promptTokenCount)
	res += fmt.Sprintf("- Candidate tokens: %d\n", u.candidateTokenCount)
	res += fmt.Sprintf("- Output tokens: %d\n", u.outputTokens)
	return res
}
