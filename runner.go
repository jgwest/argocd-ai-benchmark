package main

import (
	"argocd-ai-benchmark/client"
	"argocd-ai-benchmark/evaluations"
	"argocd-ai-benchmark/types"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func runBenchmark(targetModel string, configuration types.EvaluationConfiguration) error {

	configuration.Reporter.Init(targetModel)

	configuration.Reporter.ReportCumulativeResult("* Using model '" + targetModel + "'")

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

	resultReporter := configuration.Reporter

	resultReporter.ReportCumulativeResult(strings.Repeat("=", 80))

	resultReporter.ReportCumulativeResult(fmt.Sprintf("Overall Results: Passed %d, Total %d\n", checksPassed, checksRun))

	resultReporter.ReportCumulativeResult("Usage Details:")
	resultReporter.ReportCumulativeResult(aggregateUsage.GenerateReport())
	resultReporter.ReportCumulativeResult(fmt.Sprintln("- Elapsed time:", time.Since(startTime).Truncate(time.Second)))

	if resultReporter.IncludeJsonSummary() {

		var jsonSummary JsonSummary = JsonSummary{
			ChecksPassed: checksPassed,
			ChecksTotal:  checksRun,
			ChecksError:  checksError,
			ModelName:    targetModel,
		}

		jsonData, err := json.Marshal(jsonSummary)
		if err != nil {
			return err
		}

		resultReporter.ReportCumulativeResult("JSON Summary: " + string(jsonData))
	}

	if err := resultReporter.Close(targetModel); err != nil {
		return fmt.Errorf("error on writing results: %v ", err)
	}
	if checksError != 0 {
		return fmt.Errorf("ERROR: Some API errors errors occurred during evaluation: '%d'", checksError)
	}

	return nil
}

// func generatePathForFileReporter(targetModel string, configuration types.EvaluationConfiguration) string {
// 	modelPathName := targetModel
// 	modelPathName = strings.ReplaceAll(modelPathName, "/", "_")

// 	modelFileName := modelPathName
// 	if configuration.ProvideExternalResources {
// 		modelFileName += "_with_context"
// 	} else {
// 		modelFileName += "_no_context"
// 	}
// 	modelFileName += ".txt"

// 	return "/tmp/out/" + modelPathName + "/" + modelFileName
// }

type JsonSummary struct {
	ChecksPassed int    `json:"checksPassed"`
	ChecksError  int    `json:"checksError"`
	ChecksTotal  int    `json:"checksTotal"`
	ModelName    string `json:"modelName"`
}
