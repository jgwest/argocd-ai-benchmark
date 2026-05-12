package types

import (
	"argocd-ai-benchmark/client"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type EvaluationConfiguration struct {
	Client client.Client

	// Model is which model to evaluate (id from openrouter)
	// Model string

	// ProvideExternalResources: if true, the model will be provided with the corresponding resources (for example, the relevant Argo CD documentation page) which it can consult in order to solve the question.
	ProvideExternalResources bool

	// PrintReasoning: some models will use intermediate "reasoning" tokens to "think" about the answer before answering. When true, these reasoning tokens will be printed for debug purposes.
	PrintReasoning bool

	// NumberOfWorkers is the number of concurrent requests made to OpenRouter. A value of 1 will run the evaluations sequentially, which can be helpful for debugging.
	NumberOfWorkers int

	// Report implement determines how results are logged: for example, either to console out, or to file.
	Reporter ResultReporter
}

type EvalContext struct {
	Configuration EvaluationConfiguration
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
	Usage             client.ResponseUsage
}

func runSingleEvaluation(toEvaluateParam Evaluation, mainContext EvalContext) (EvaluationRunResult, string, error) {

	var res EvaluationRunResult

	var ob outBuffer

	var markdownReferenceMaterial []string

	if mainContext.Configuration.ProvideExternalResources {
		for _, resourceURL := range toEvaluateParam.initial.resourceURLS {

			markdownContentsFromURL, err := mainContext.resourceCache.getExternalContent(resourceURL, &ob)
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
		ob.println(color.New(color.Faint, color.Bold).Sprint(">"), line)
		// ob.println("> " + line)
	}

	if len(markdownReferenceMaterial) > 0 && mainContext.Configuration.ProvideExternalResources {

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

	clientInstance := mainContext.Configuration.Client.NewInstance()

	cr, err := clientInstance.SendMessage(promptText)
	if err != nil {
		return EvaluationRunResult{}, "", fmt.Errorf("error on sending message: %v", err)
	}

	responseSanitized := sanitizeString(cr.ResponseContent)

	if mainContext.Configuration.PrintReasoning {

		output := cr.ReasoningContent
		if output != "" {
			ob.println("Reasoning:")
			ob.println(output)
		}
	}

	ob.println("A:", color.New(color.FgYellow).Sprint(responseSanitized))

	matches := false

	if len(toEvaluateParam.exactAnswers) > 0 {

		// Sanitize all expected answers
		expectedSanitized := make([]string, len(toEvaluateParam.exactAnswers))

		if len(toEvaluateParam.exactAnswers) == 1 {
			expectedSanitized[0] = sanitizeString(toEvaluateParam.exactAnswers[0])
			ob.println("- Expected:", color.New(color.FgBlue).Sprint(expectedSanitized[0]))
		} else {
			ob.println("Expected one of:")
			for i, answer := range toEvaluateParam.exactAnswers {
				expectedSanitized[i] = sanitizeString(answer)
				ob.println("-", color.New(color.FgBlue).Sprint(expectedSanitized[i]))
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
		ob.println("-", color.New(color.FgGreen, color.Bold).Sprint("PASS"))
	} else {
		ob.println("-", color.New(color.FgRed, color.Bold).Sprint("FAIL"))
	}

	res.EvaluationsRun++

	res.Usage = cr.Usage

	return res, ob.out, nil
}

func (cache *externalResourceCache) getExternalContent(url string, ob *outBuffer) (string, error) {

	cache.Lock()
	defer cache.Unlock()

	val, inCache := cache.externalResourceCache[url]
	if inCache {
		if ob != nil {
			ob.println("* Retrieving content from '" + url + "' [cached]")
			return val, nil
		} else {
			fmt.Println("* Retrieving content from '" + url + "' [cached]")
			return val, nil
		}
	}

	if ob != nil {
		ob.println("* Retrieving content from '" + url)
	} else {
		fmt.Println("* Retrieving content from '" + url)
	}

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
