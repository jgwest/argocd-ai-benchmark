package types

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

func RunEvaluationsInParallel(evaluations []Evaluation, configuration EvaluationConfiguration) []ExecutionResult {
	var wg sync.WaitGroup

	inputWorkChannel := make(chan Evaluation, len(evaluations))
	outputResultsChannel := make(chan ExecutionResult, len(evaluations))

	existsAnyFocused := ExistsAnyFocused(evaluations) // Are any of the evaluations focused?

	for _, evaluation := range evaluations {
		// If there exists at least one focused, then skip all non-focused
		if existsAnyFocused && !evaluation.Focused() {
			continue
		}

		inputWorkChannel <- evaluation // Queue the work
	}

	mainContext := evalContext{
		configuration: configuration,
		resourceCache: &externalResourceCache{},
	}

	for range configuration.NumberOfWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case e := <-inputWorkChannel:

					runResult, outStr, err := runSingleEvaluation(e, mainContext)

					fmt.Println(strings.Repeat("=", 50))

					fmt.Println(outStr)
					if err != nil {
						log.Printf("Error running test on %v: %v", e.Name(), err)
					}

					res := ExecutionResult{
						Evaluation: e,
						RunResult:  runResult,
						OutStr:     outStr,
						Err:        err,
					}
					outputResultsChannel <- res

				default:
					return
				}
			}

		}()

	}

	// Wait for workers to complete
	wg.Wait()

	// Read channel results into a slice
	res := []ExecutionResult{}
out:
	for {
		select {
		case result := <-outputResultsChannel:
			res = append(res, result)
		default:
			break out
		}
	}

	return res

}

type ExecutionResult struct {
	Evaluation Evaluation
	RunResult  EvaluationRunResult
	OutStr     string
	Err        error
}
