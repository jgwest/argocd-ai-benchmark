package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

type MainContext struct {
	client *openai.Client
	model  string
}

func main() {
	targetModel := "openai/gpt-oss-20b"
	// targetModel := "google/gemini-2.5-flash"
	// targetModel := "google/gemini-2.5-pro"
	// targetModel := "deepseek/deepseek-chat-v3.1"
	// targetModel := "google/gemma-3-27b-it"
	// targetModel := "google/gemma-3-12b-it"

	mainContext := MainContext{
		client: getClient(),
		model:  targetModel,
	}

	checksRun := 0
	checksPassed := 0

	// List YAML files in data/ directory
	yamlFiles, err := listYAMLFiles("data/")
	if err != nil {
		log.Fatalf("Error listing YAML files: %v", err)
	}

	fmt.Printf("Found %d YAML files in data/ directory:\n", len(yamlFiles))
	for _, file := range yamlFiles {
		fmt.Printf("  - %s\n", file)
	}
	fmt.Println()

	// Run tests on each YAML file
	for _, yamlFile := range yamlFiles {
		fmt.Printf("Running tests on: %s\n", yamlFile)
		fmt.Println(strings.Repeat("=", 50))

		passed, run, err := runTestOnFile(yamlFile, mainContext)
		if err != nil {
			log.Printf("Error running test on %s: %v", yamlFile, err)
			continue
		}

		checksRun += run
		checksPassed += passed

		fmt.Printf("File %s: Passed %d/%d tests\n", yamlFile, passed, run)
		fmt.Println()
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Overall Results: Passed %d, Total %d\n", checksPassed, checksRun)
}

func runTestOnFile(path string, mainContext MainContext) (int, int, error) {
	checksRun := 0

	// Parse the YAML file
	testDoc, err := parseYAMLFile(path)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}
	checksPassed := 0
	fmt.Println(testDoc.Name)

	for x, entry := range testDoc.Checks {

		fmt.Println()
		fmt.Printf("%d)\n", x)
		lines := strings.Split(strings.TrimSpace(entry.Question), "\n")
		for _, line := range lines {
			fmt.Println("> " + line)
		}

		var conversationHistory []openai.ChatCompletionMessage

		// Add user message to conversation history
		conversationHistory = append(conversationHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: entry.Question,
		})

		// Get response from AI
		response, err := chatWithHistoryNew(mainContext, conversationHistory)
		if err != nil {
			return 0, 0, err
		}

		responseSanitized := sanitizeString(response)
		fmt.Println("A:", responseSanitized)

		// Sanitize all expected answers
		expectedSanitized := make([]string, len(entry.Answers))
		for i, answer := range entry.Answers {
			expectedSanitized[i] = sanitizeString(answer)
			fmt.Println("- Expected:", expectedSanitized[i])
		}

		// Check if response matches any of the expected answers
		matches := false
		for _, expected := range expectedSanitized {
			if responseSanitized == expected {
				matches = true
				fmt.Println("  Match:", responseSanitized)
				break
			}
		}

		if matches {
			checksPassed++
			fmt.Println("- PASS")
		} else {
			fmt.Println("- FAIL")
		}

		// Add AI response to conversation history
		conversationHistory = append(conversationHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: response,
		})

		checksRun++
	}

	return checksPassed, checksRun, nil
}

func sanitizeString(str string) string {
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	return str
}

func getClient() *openai.Client {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	// Create OpenRouter client with custom config
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	client := openai.NewClientWithConfig(config)

	return client
}

// parseYAMLFile reads and parses a YAML file into a TestDocument struct
func parseYAMLFile(filename string) (*TestDocument, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", filename, err)
	}

	var testDoc TestDocument
	err = yaml.Unmarshal(data, &testDoc)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	return &testDoc, nil
}

// chatWithHistory performs a chat completion with full conversation history
func chatWithHistoryNew(mainContext MainContext, messages []openai.ChatCompletionMessage) (string, error) {
	ctx := context.Background()

	resp, err := mainContext.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    mainContext.model,
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

// listYAMLFiles returns a list of all YAML files in the specified directory
func listYAMLFiles(dir string) ([]string, error) {
	var yamlFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a YAML file
		if !info.IsDir() && (strings.HasSuffix(strings.ToLower(path), ".yaml") || strings.HasSuffix(strings.ToLower(path), ".yml")) {
			yamlFiles = append(yamlFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %v", dir, err)
	}

	return yamlFiles, nil
}
