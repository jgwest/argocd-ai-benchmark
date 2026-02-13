package main

import (
	"argocd-ai-benchmark/checks"
	"argocd-ai-benchmark/types"
	"fmt"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
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
