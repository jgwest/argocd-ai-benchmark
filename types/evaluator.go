package types

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type MainContext struct {
	Client *openai.Client
	Model  string
}

func RunTestOnFile(toEvaluateParam Evaluation, mainContext MainContext) (int, int, error) {
	checksRun := 0

	// markdownContents, err := downloadURL("https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/docs/user-guide/app_deletion.md")
	// if err != nil {
	// 	log.Panic("unable to download from URL", err)
	// }

	checksPassed := 0
	fmt.Println("[", toEvaluateParam.initial.name, "| labels: ", toEvaluateParam.initial.labels, "]")

	var promptText string

	// promptText := "The following is reference information which may be useful:\n" + markdownContents + "\nEnd reference information.\n\n\n"
	promptText += TrimIndent(toEvaluateParam.initial.prompt)

	promptText = strings.TrimSpace(promptText)

	switch toEvaluateParam.initial.promptType {
	case promptType_TrueOrFalse:
		promptText += "\n\nProvide ONLY the answer, expressed as a single letter, either `T` (true) or `F` (false).\n"
	case promptType_MultipleChoice:
		promptText += "\n\nProvide ONLY the answer. The answer will be a single letter from the multiple-choice list.\n"
	case promptType_Generic:
	// no-op
	default:
		return 0, 0, fmt.Errorf("unrecognized prompt type")
	}

	fmt.Println()
	lines := strings.Split(promptText, "\n")
	for _, line := range lines {
		fmt.Println("> " + line)
	}

	var conversationHistory []openai.ChatCompletionMessage

	// Add user message to conversation history
	conversationHistory = append(conversationHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: promptText,
	})

	// Get response from AI
	response, err := chatWithHistory(mainContext, conversationHistory)
	if err != nil {
		return 0, 0, err
	}

	responseSanitized := sanitizeString(response)
	fmt.Println("A:", responseSanitized)

	matches := false

	if len(toEvaluateParam.exactAnswers) > 0 {

		// Sanitize all expected answers
		expectedSanitized := make([]string, len(toEvaluateParam.exactAnswers))

		if len(toEvaluateParam.exactAnswers) == 1 {
			expectedSanitized[0] = sanitizeString(toEvaluateParam.exactAnswers[0])
			fmt.Println("- Expected:", expectedSanitized[0])
		} else {
			fmt.Println("Expected one of:")
			for i, answer := range toEvaluateParam.exactAnswers {
				expectedSanitized[i] = sanitizeString(answer)
				fmt.Println("-", expectedSanitized[i])
			}
			fmt.Println()

		}

		// Check if response matches any of the expected answers
		if slices.Contains(expectedSanitized, responseSanitized) {
			matches = true
			fmt.Println("  Match:", responseSanitized)
		}

	} else {
		return 0, 0, fmt.Errorf("missing evaluation: %v", toEvaluateParam)
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

	return checksPassed, checksRun, nil
}

// TrimIndent removes common leading whitespace from each line of a string.
func TrimIndent(s string) string {

	// Split the string into lines.
	lines := strings.Split(s, "\n")

	// Find the minimum indentation of non-empty lines.
	minIndent := math.MaxInt32
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue // Skip empty lines
		}

		indent := 0
		for _, r := range line {
			if r == ' ' || r == '\t' {
				indent++
			} else {
				break
			}
		}

		if indent < minIndent {
			minIndent = indent
		}
	}

	// If no indented lines were found, return the original string.
	if minIndent == math.MaxInt32 {
		return s
	}

	var trimmedLines []string
	for _, line := range lines {
		if len(line) > minIndent {
			trimmedLines = append(trimmedLines, line[minIndent:])
		} else {
			trimmedLines = append(trimmedLines, line)
		}
	}

	// Join the lines back together.
	return strings.Join(trimmedLines, "\n")
}

func sanitizeString(str string) string {
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	return str
}

// chatWithHistory performs a chat completion with full conversation history
func chatWithHistory(mainContext MainContext, messages []openai.ChatCompletionMessage) (string, error) {
	ctx := context.Background()

	resp, err := mainContext.Client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    mainContext.Model,
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

func GetClient() *openai.Client {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	client := openai.NewClientWithConfig(config)

	return client
}

func IsFocused(e Evaluation) bool {
	return e.initial.focus
}

func ExistsAnyFocused(param []Evaluation) bool {

	for _, eval := range param {
		if eval.initial.focus {
			return true
		}
	}

	return false
}

func downloadURL(url string) (string, error) {

	// Send an HTTP GET request to the URL.
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// Convert the body to a string and print it.
	content := string(body)
	return content, nil
}
