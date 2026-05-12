package main

import (
	"argocd-ai-benchmark/client/gemini"
	"argocd-ai-benchmark/client/openrouter"
	"argocd-ai-benchmark/types"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	// OpenRouter ----------------------------------------------------------------

	// Reasoning (uses a large number of output tokens, BY DEFAULT):

	// targetModel := "anthropic/claude-haiku-4.5" // (36/58,51/58)($1/$5) (minimal reasoning, AFAICT)
	// ------------------

	// DONE:
	// targetModel := "z-ai/glm-5.1" //
	// targetModel := "moonshotai/kimi-k2.6" //
	// targetModel := "deepseek/deepseek-v4-pro" //
	// targetModel := "qwen/qwen3.6-plus" //
	// targetModel := "openai/gpt-oss-120b" //
	// targetModel := "qwen/qwen3.6-27b" //
	// targetModel := "openai/gpt-5.5" //
	// targetModel := "anthropic/claude-opus-4.7" //

	// DONE:
	// targetModel := "google/gemma-4-31b-it" //
	// targetModel := "google/gemma-4-26b-a4b-it" //
	// targetModel := "google/gemma-3-27b-it" //
	// targetModel := "google/gemma-3-12b-it" //
	// targetModel := "google/gemma-3-4b-it" //
	// targetModel := "google/gemma-2-27b-it" //

	// targetModel := "ibm-granite/granite-4.0-h-micro" //
	// targetModel := "ibm-granite/granite-4.1-8b" //
	// targetModel := "qwen/qwen3.5-9b"

	// targetModel := "mistralai/ministral-8b-2512"
	// targetModel := "mistralai/ministral-14b-2512"

	// Gemini (via Google Vertex AI) --------------------------------------------------------------------

	// Values for this field can be found on the individual model pages of https://docs.cloud.google.com/vertex-ai/generative-ai/docs/models

	// DONE:
	// targetModel := "gemini-3.1-pro-preview" //
	// targetModel := "gemini-3-pro-preview" //
	// targetModel := "gemini-2.5-pro" //
	// targetModel := "gemini-3-flash-preview"
	// targetModel := "gemini-2.5-flash" //
	// targetModel := "gemini-2.0-flash"

	// ---------------------------------------------------------------------------

	configuration := types.EvaluationConfiguration{
		ProvideExternalResources: false,
		PrintReasoning:           true,
		NumberOfWorkers:          1,
	}

	var resultReporter types.ResultReporter

	resultReporter = types.NewFileResultReporter(generatePathForFileReporter("/home/jgw/workspace/Projects/argocd-ai-benchmark/results/2026-Q2", targetModel, configuration))
	// resultReporter = &types.StdoutResultReporter{}

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

	configuration.Reporter = resultReporter

	if err := runBenchmark(targetModel, configuration); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

}

func generatePathForFileReporter(rootOutputPath string, targetModel string, configuration types.EvaluationConfiguration) string {
	modelPathName := targetModel
	modelPathName = strings.ReplaceAll(modelPathName, "/", "_")

	modelFileName := modelPathName
	if configuration.ProvideExternalResources {
		modelFileName += "_with_context"
	} else {
		modelFileName += "_no_context"
	}
	modelFileName += ".txt"

	rootOutputPath = strings.TrimSuffix(rootOutputPath, "/") // remove trailing '/' if present

	return rootOutputPath + "/" + modelPathName + "/" + modelFileName
}
