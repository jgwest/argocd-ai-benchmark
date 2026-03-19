package main

import (
	"argocd-ai-benchmark/checks"
	"argocd-ai-benchmark/types"
	"fmt"
	"log"
	"strings"
)

func main() {
	targetModel := "openai/gpt-oss-120b" // 35/46
	// targetModel := "google/gemini-2.5-flash" // 25/47
	// targetModel := "google/gemini-2.5-pro" // 38/47
	// targetModel := "deepseek/deepseek-chat-v3.1" // 27/47
	// targetModel := "google/gemma-3-27b-it" // 27/46
	// targetModel := "google/gemma-3-12b-it" // 22/47
	// targetModel := "qwen/qwen3-coder-30b-a3b-instruct" // 20/46

	fmt.Println("* Using model '" + targetModel + "'")

	mainContext := types.MainContext{
		Client: types.GetClient(),
		Model:  targetModel,
	}

	checksRun := 0
	checksPassed := 0

	checks.Init()

	evaluations := types.Evaluations()

	existsAnyFocused := types.ExistsAnyFocused(evaluations)

	for _, evaluation := range evaluations {

		if existsAnyFocused && !types.IsFocused(evaluation) {
			continue
		}

		fmt.Println(strings.Repeat("=", 50))

		passed, run, err := types.RunTestOnFile(evaluation, mainContext)
		if err != nil {
			log.Printf("Error running test on %v: %v", evaluation, err)
			continue
		}

		checksRun += run
		checksPassed += passed

	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Overall Results: Passed %d, Total %d\n", checksPassed, checksRun)
}
