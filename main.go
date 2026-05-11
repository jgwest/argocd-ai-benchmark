package main

import (
	"argocd-ai-benchmark/client"
	"argocd-ai-benchmark/client/gemini"
	"argocd-ai-benchmark/client/openrouter"
	"argocd-ai-benchmark/evaluations"
	"argocd-ai-benchmark/types"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	// OpenRouter ----------------------------------------------------------------

	// Reasoning (uses a large number of output tokens, BY DEFAULT):
	// targetModel := "google/gemini-2.5-pro" // 38/47 ($1.25/$10)
	// targetModel := "openai/gpt-oss-120b" // 43/69 | 64/69  ($0.04/$0.40)

	// Model prices sorted descending by input token price.
	// "No" reasoning (uses a small number of output tokens, BY DEFAULT):
	// targetModel := "anthropic/claude-haiku-4.5" // (36/58,51/58)($1/$5) (minimal reasoning, AFAICT)
	// targetModel := "google/gemini-2.5-flash" // 25/47 ($0.30/$2.50)
	// targetModel := "deepseek/deepseek-chat-v3.1" // 27/47 ($0.20/$0.80)
	// targetModel := "google/gemma-3-27b-it" // 27/69 | 52/69 ($0.09/$0.16)
	// targetModel := "qwen/qwen3-coder-30b-a3b-instruct" // 20/46 ($0.06/$0.25)
	// targetModel := "google/gemma-3-12b-it" // 26/69 | 55/69 ($0.03/$0.10)
	// targetModel := "ibm-granite/granite-4.0-h-micro" // 21/67 | 25/67 ($0.017/$0.11)

	// Gemini --------------------------------------------------------------------
	targetModel := "gemini-2.5-flash" // 45/79 | 74/79
	// targetModel := "gemini-2.5-pro" // 66/79 | 74/79
	// targetModel := "gemini-3-pro-preview" // 72/79 |
	// targetModel := "gemini-2.5-flash-lite" // 34/79 | 67/79
	// targetModel := "gemini-2.0-flash" // 41/79 | 69/79

	// ---------------------------------------------------------------------------

	fmt.Println("* Using model '" + targetModel + "'")

	configuration := types.EvaluationConfiguration{
		ProvideExternalResources: true,
		PrintReasoning:           true,
		NumberOfWorkers:          5,
	}

	openRouterAPIKey := os.Getenv("OPENROUTER_API_KEY")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	if openRouterAPIKey != "" {
		configuration.Client = &openrouter.OpenRouterClient{
			APIKey: openRouterAPIKey,
			Model:  targetModel,
		}
	} else if geminiAPIKey != "" {
		configuration.Client = &gemini.GeminiClient{
			APIKey: geminiAPIKey,
			Model:  targetModel,
		}
	} else {
		log.Fatal("no API keys defined: OPENROUTER_API_KEY or GEMINI_API_KEY")
		return
	}

	// Queue the work
	evaluations.AddEvaluations()
	allEvaluations := types.AllEvaluations()

	startTime := time.Now()

	results := types.RunEvaluationsInParallel(allEvaluations, configuration)

	// Calculate results
	var aggregateUsage client.ResponseUsage
	checksRun := 0
	checksPassed := 0
	checksError := 0

	for _, result := range results {

		if result.Err != nil {
			checksError++
			continue
		}

		checksRun += result.RunResult.EvaluationsRun
		checksPassed += result.RunResult.EvaluationsPassed

		if aggregateUsage == nil {
			aggregateUsage = result.RunResult.Usage
		} else {
			aggregateUsage = aggregateUsage.Aggregate(result.RunResult.Usage)
		}
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Overall Results: Passed %d, Total %d\n", checksPassed, checksRun)
	if checksError != 0 {
		fmt.Printf("- WARNING: '%d' errors occurred\n", checksError)
	}

	fmt.Println("Usage Details:")
	fmt.Print(aggregateUsage.GenerateReport())
	fmt.Println("- Elapsed time:", time.Since(startTime).Truncate(time.Second))
}
