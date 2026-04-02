package main

import (
	"argocd-ai-benchmark/checks"
	"argocd-ai-benchmark/types"
	"fmt"
	"log"
	"strings"
)

func main() {

	// Reasoning (uses a large number of output tokens, BY DEFAULT):
	// targetModel := "google/gemini-2.5-pro" // 38/47 ($1.25/$10)
	// targetModel := "openai/gpt-oss-120b" // 35/46 ($0.04/$0.40)

	// Model prices sorted descending by input token price.
	// "No" reasoning (uses a small number of output tokens, BY DEFAULT):
	// targetModel := "anthropic/claude-haiku-4.5" // ($1/$5) (minimal reasoning, AFAICT)
	targetModel := "google/gemini-2.5-flash" // 25/47 ($0.30/$2.50)
	// targetModel := "deepseek/deepseek-chat-v3.1" // 27/47 ($0.20/$0.80)
	// targetModel := "google/gemma-3-27b-it" // 27/46 ($0.09/$0.16)
	// targetModel := "qwen/qwen3-coder-30b-a3b-instruct" // 20/46 ($0.06/$0.25)
	// targetModel := "google/gemma-3-12b-it" // 22/47 ($0.03/$0.10)

	fmt.Println("* Using model '" + targetModel + "'")

	mainContext := types.MainContext{
		// Client:                 types.GetClient(),
		Model:                  targetModel,
		ExternalResourceCache:  map[string]string{},
		AllowExternalResources: false,
		PrintReasoning:         true,
	}

	checksRun := 0
	checksPassed := 0

	checks.Init()

	evaluations := types.Evaluations()

	existsAnyFocused := types.ExistsAnyFocused(evaluations)

	var aggregateUsage types.Usage

	for _, evaluation := range evaluations {

		if existsAnyFocused && !types.IsFocused(evaluation) {
			continue
		}

		fmt.Println(strings.Repeat("=", 50))

		runResult, err := types.RunTestOnFile(evaluation, mainContext) // JGW-TODO: Rename Test
		if err != nil {
			log.Printf("Error running test on %v: %v", evaluation, err)
			continue
		}

		checksRun += runResult.ChecksRun
		checksPassed += runResult.ChecksPassed

		aggregateUsage = types.MergeUsage(aggregateUsage, runResult.Usage)
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Overall Results: Passed %d, Total %d\n", checksPassed, checksRun)

	fmt.Println("Usage Details:")

	fmt.Printf("- Prompt Tokens: %d\n", aggregateUsage.PromptTokens)
	fmt.Printf("- Completion Tokens: %d\n", aggregateUsage.CompletionTokens)
	fmt.Printf("- Total Tokens: %d\n", aggregateUsage.TotalTokens)
	fmt.Printf("- Cost: $%.8f\n", aggregateUsage.Cost)
}
