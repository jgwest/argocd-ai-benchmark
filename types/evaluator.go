package types

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

type EvaluationConfiguration struct {
	// Model is which model to evaluate (id from openrouter)
	Model string

	// ProvideExternalResources: if true, the model will be provided with the corresponding resources (for example, the relevant Argo CD documentation page) which it can consult in order to solve the question.
	ProvideExternalResources bool

	// PrintReasoning: some models will use intermediate "reasoning" tokens to "think" about the answer before answering. When true, these reasoning tokens will be printed for debug purposes.
	PrintReasoning bool

	// NumberOfWorkers is the number of concurrent requests made to OpenRouter. A value of 1 will run the evaluations sequentially, which can be helpful for debugging.
	NumberOfWorkers int
}

type evalContext struct {
	configuration EvaluationConfiguration
	resourceCache *externalResourceCache
}

// externalResourceCache is a cache of the contents of external URLs (for example, Argo CD doc pages)
type externalResourceCache struct {
	sync.Mutex

	// externalResourceCache: key: URL of resource, value: contents of resource downloaded from URL
	externalResourceCache map[string]string
}

type EvaluationRunResult struct {
	EvaluationsPassed int
	EvaluationsRun    int
	Usage             Usage
}

func runSingleEvaluation(toEvaluateParam Evaluation, mainContext evalContext) (EvaluationRunResult, string, error) {

	var res EvaluationRunResult

	var ob outBuffer

	var markdownReferenceMaterial []string

	if mainContext.configuration.ProvideExternalResources {
		for _, resourceURL := range toEvaluateParam.initial.resourceURLS {

			markdownContentsFromURL, err := mainContext.resourceCache.getExternalContent(resourceURL)
			if err != nil {
				return EvaluationRunResult{}, ob.out, fmt.Errorf("unable to download from URL: %v", err)
			}
			markdownReferenceMaterial = append(markdownReferenceMaterial, markdownContentsFromURL)
			ob.println()
		}
	}

	ob.println("[", toEvaluateParam.initial.name, "| labels: ", toEvaluateParam.initial.labels, "]")

	var promptText string

	promptText += trimIndent(toEvaluateParam.initial.prompt)

	promptText = strings.TrimSpace(promptText) + "\n\n"

	switch toEvaluateParam.initial.promptType {
	case promptType_TrueOrFalse:
		promptText += "Provide ONLY the answer, expressed as a single letter, either `T` (true) or `F` (false). Don't write any other text.\n"
	case promptType_MultipleChoice:
		promptText += "Provide ONLY the answer. The answer will be a single letter from the multiple-choice list. Don't write any other text.\n"
	case promptType_Generic:
	// no-op
	default:
		return res, ob.out, fmt.Errorf("unrecognized prompt type")
	}
	ob.println()
	lines := strings.SplitSeq(promptText, "\n")
	for line := range lines {
		ob.println("> " + line)
	}

	if len(markdownReferenceMaterial) > 0 && mainContext.configuration.ProvideExternalResources {

		ob.println("> (reference text)")

		promptText += "--------------------------------------------------------------\n"
		promptText += "The following is reference material which may help to answer the above question. All text below is not part of the question.\n"
		promptText += "\n"

		for _, markdownContent := range markdownReferenceMaterial {
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
	response, usage, err := chatWithHistory(mainContext.configuration.Model, conversationHistory)
	if err != nil {
		return res, ob.out, fmt.Errorf("error on retrieving response: %v", err)
	}

	responseSanitized := sanitizeString(response.Content)

	if mainContext.configuration.PrintReasoning {

		output := ""

		for _, toolCall := range response.ToolCalls {
			output += fmt.Sprintln("-", toolCall)
		}

		if response.Reasoning != "" {
			output += fmt.Sprintln("-", response.Reasoning)
		}

		if output != "" {
			ob.println("Reasoning:")
			ob.println(output)
		}
	}

	ob.println("A:", responseSanitized)

	matches := false

	if len(toEvaluateParam.exactAnswers) > 0 {

		// Sanitize all expected answers
		expectedSanitized := make([]string, len(toEvaluateParam.exactAnswers))

		if len(toEvaluateParam.exactAnswers) == 1 {
			expectedSanitized[0] = sanitizeString(toEvaluateParam.exactAnswers[0])
			ob.println("- Expected:", expectedSanitized[0])
		} else {
			fmt.Println("Expected one of:")
			for i, answer := range toEvaluateParam.exactAnswers {
				expectedSanitized[i] = sanitizeString(answer)
				ob.println("-", expectedSanitized[i])
			}
			ob.println()

		}

		// Check if response matches any of the expected answers
		if slices.Contains(expectedSanitized, responseSanitized) {
			matches = true
			ob.println("  Match:", responseSanitized)
		}

	} else {
		return res, ob.out, fmt.Errorf("missing evaluation: %v", toEvaluateParam)
	}

	if matches {
		res.EvaluationsPassed++
		ob.println("- PASS")
	} else {
		ob.println("- FAIL")
	}

	// Add AI response to conversation history
	conversationHistory = append(conversationHistory, response.Content)

	res.EvaluationsRun++

	res.Usage = *usage

	return res, ob.out, nil
}

func (cache *externalResourceCache) getExternalContent(url string) (string, error) {

	cache.Lock()
	defer cache.Unlock()

	val, inCache := cache.externalResourceCache[url]
	if inCache {
		fmt.Println("* Retrieving content from '" + url + "' [cached]")
		return val, nil
	}

	fmt.Println("* Retrieving content from '" + url)
	markdownContents, err := downloadURL(url)
	if err != nil {
		return "", fmt.Errorf("unable to download URL: '%s'. error: %v", url, err)
	}

	if cache.externalResourceCache == nil {
		cache.externalResourceCache = map[string]string{}
	}

	cache.externalResourceCache[url] = markdownContents

	return markdownContents, nil
}

func (e Evaluation) Focused() bool {
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
