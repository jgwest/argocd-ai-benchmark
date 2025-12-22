package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// startInteractiveChat begins an interactive chat session
func startInteractiveChat(client *openai.Client) {
	scanner := bufio.NewScanner(os.Stdin)

	// Conversation history
	var conversationHistory []openai.ChatCompletionMessage

	model := "openai/gpt-oss-20b"

	fmt.Println("🤖 OpenRouter.AI Chat Client")
	fmt.Println("Type 'exit' or 'quit' to end the conversation")
	fmt.Println("Type '/model <model-name>' to change the model")
	fmt.Println("Type '/clear' to clear conversation history")
	fmt.Printf("Using model: %s\n", model)
	fmt.Println("----------------------------------------")

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())

		// Handle special commands
		if userInput == "exit" || userInput == "quit" {
			fmt.Println("Goodbye! 👋")
			break
		}

		if userInput == "/clear" {
			conversationHistory = []openai.ChatCompletionMessage{}
			fmt.Println("Conversation history cleared.")
			continue
		}

		if userInput == "" {
			continue
		}

		// Add user message to conversation history
		conversationHistory = append(conversationHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		// Get response from AI
		response, err := chatWithHistory(client, model, conversationHistory)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		// Add AI response to conversation history
		conversationHistory = append(conversationHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: response,
		})

		// Display AI response
		fmt.Printf("AI: %s\n", response)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}

// chatWithHistory performs a chat completion with full conversation history
func chatWithHistory(client *openai.Client, model string, messages []openai.ChatCompletionMessage) (string, error) {
	ctx := context.Background()

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("unexpected number of choices: %v", resp.Choices)
	}

	return resp.Choices[0].Message.Content, nil
}
