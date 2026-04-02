package types

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

type MainContext struct {
	Model                  string
	ExternalResourceCache  map[string]string
	AllowExternalResources bool
	PrintReasoning         bool
}

type TestRunResult struct {
	ChecksPassed int
	ChecksRun    int
	Usage        Usage
}

func RunTestOnFile(toEvaluateParam Evaluation, mainContext MainContext) (TestRunResult, error) {

	var res TestRunResult

	var markdownContents []string

	if mainContext.AllowExternalResources {
		for _, resourceURL := range toEvaluateParam.initial.resourceURLS {

			markdownContentsFromURL, err := getExternalContent(resourceURL, mainContext.ExternalResourceCache)
			if err != nil {
				log.Panic("unable to download from URL", err)
			}
			markdownContents = append(markdownContents, markdownContentsFromURL)
			fmt.Println()
		}
	}

	// checksPassed := 0
	fmt.Println("[", toEvaluateParam.initial.name, "| labels: ", toEvaluateParam.initial.labels, "]")

	var promptText string

	promptText += trimIndent(toEvaluateParam.initial.prompt)

	promptText = strings.TrimSpace(promptText) + "\n\n"

	switch toEvaluateParam.initial.promptType {
	case promptType_TrueOrFalse:
		promptText += "Provide ONLY the answer, expressed as a single letter, either `T` (true) or `F` (false).\n"
	case promptType_MultipleChoice:
		promptText += "Provide ONLY the answer. The answer will be a single letter from the multiple-choice list.\n"
	case promptType_Generic:
	// no-op
	default:
		return res, fmt.Errorf("unrecognized prompt type")
	}
	fmt.Println()
	lines := strings.Split(promptText, "\n")
	for _, line := range lines {
		fmt.Println("> " + line)
	}

	if len(markdownContents) > 0 && mainContext.AllowExternalResources {

		fmt.Println("> (reference text)")

		promptText += "--------------------------------------------------------------\n"
		promptText += "The following is reference material which may help to answer the above question. All text below is not part of the question.\n"
		promptText += "\n"

		for _, markdownContent := range markdownContents {
			promptText += markdownContent + "\n"
			promptText += "-----------------------------------------\n"
		}

		promptText += "--------------------------------------------------------------\n"
		promptText += "This is the end of the reference material."
		promptText += "\n"
	}

	var conversationHistory []string

	conversationHistory = append(conversationHistory, promptText)

	// Get response from AI
	response, usage, err := chatWithHistory(mainContext, conversationHistory)
	if err != nil {
		return res, err
	}

	responseSanitized := sanitizeString(response.Content)

	if mainContext.PrintReasoning {

		output := ""

		for _, toolCall := range response.ToolCalls {
			output += fmt.Sprintln("-", toolCall)
		}

		if response.Reasoning != "" {
			output += fmt.Sprintln("-", response.Reasoning)
		}

		if output != "" {
			fmt.Println("Reasoning:")
			fmt.Println(output)
		}
	}

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
		return res, fmt.Errorf("missing evaluation: %v", toEvaluateParam)
	}

	if matches {
		res.ChecksPassed++
		fmt.Println("- PASS")
	} else {
		fmt.Println("- FAIL")
	}

	// Add AI response to conversation history
	conversationHistory = append(conversationHistory, response.Content)

	res.ChecksRun++

	res.Usage = *usage

	return res, nil
}

func getExternalContent(url string, cache map[string]string) (string, error) {

	val, inCache := cache[url]
	if inCache {
		fmt.Println("* Retrieving content from '" + url + "' [cached]")
		return val, nil
	}

	fmt.Println("* Retrieving content from '" + url)
	markdownContents, err := downloadURL(url)
	if err != nil {
		return "", fmt.Errorf("unable to download URL: '%s'. error: %v", url, err)
	}

	cache[url] = markdownContents

	return markdownContents, nil
}
func sanitizeString(str string) string {
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	return str
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
