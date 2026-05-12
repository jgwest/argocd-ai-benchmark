package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	path := "/home/jgw/workspace/Projects/argocd-ai-benchmark/results/2026-Q2"

	w := csv.NewWriter(os.Stdout)
	// if err := w.Write([]string{"file", "checksPassed", "checksError", "checksTotal", "modelName"}); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	err := filepath.WalkDir(path, func(fullPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(fullPath) != ".txt" {
			return nil
		}
		return processFile(fullPath, w)
	})

	w.Flush()
	if ferr := w.Error(); ferr != nil && err == nil {
		err = ferr
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func processFile(path string, w *csv.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%s: open: %w", path, err)
	}
	defer f.Close()

	raw, err := scanForJSONSummaryLine(f)
	if err != nil {
		if errors.Is(err, errNoJSONSummaryLine) {
			return fmt.Errorf("%s: %w", path, err)
		}
		return fmt.Errorf("%s: read: %w", path, err)
	}

	var summary jsonSummary
	if err := json.Unmarshal(raw, &summary); err != nil {
		return fmt.Errorf("%s: parse JSON: %w", path, err)
	}

	withContext := true
	if strings.Contains(path, "_no_context") {
		withContext = false
	}

	contextSuffix := " (With Context)"
	if !withContext {
		contextSuffix = " (Without Context)"
	}

	suffixAddMap := map[string]string{
		"anthropic/claude-opus-4.7":       "April 2026",
		"deepseek/deepseek-v4-pro":        "April 2026",
		"gemini-2.0-flash":                "February 2025",
		"gemini-2.5-flash":                "June 2025",
		"gemini-2.5-pro":                  "June 2025",
		"gemini-3-flash-preview":          "December 2025",
		"gemini-3-pro-preview":            "November 2025",
		"gemini-3.1-pro-preview":          "Feb 2026",
		"google/gemma-2-27b-it":           "June 2024",
		"google/gemma-3-4b-it":            "March 2025",
		"google/gemma-3-12b-it":           "March 2025",
		"google/gemma-3-27b-it":           "March 2025",
		"google/gemma-4-26b-a4b-it":       "April 2026",
		"google/gemma-4-31b-it":           "April 2026",
		"ibm-granite/granite-4.0-h-micro": "October 2025",
		"ibm-granite/granite-4.1-8b":      "April 2026",
		"mistralai/ministral-8b-2512":     "December 2025",
		"mistralai/ministral-14b-2512":    "December 2025",
		"moonshotai/kimi-k2.6":            "April 2026",
		"openai/gpt-5.5":                  "April 2026",
		"openai/gpt-oss-120b":             "August 2025",
		"qwen/qwen3.5-9b":                 "Feb 2026",
		"qwen/qwen3.6-27b":                "April 2026",
		"qwen/qwen3.6-plus":               "April 2026",
		"z-ai/glm-5.1":                    "April 2026",
	}

	keys := make([]string, 0, len(suffixAddMap))
	for k := range suffixAddMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) > len(keys[j])
		}
		return keys[i] < keys[j]
	})

	var dateLabel string
	for _, modelKey := range keys {
		if strings.Contains(summary.ModelName, modelKey) {
			dateLabel = suffixAddMap[modelKey]
			break
		}
	}
	if dateLabel == "" {
		return fmt.Errorf("%s: no suffixAddMap key matches modelName %q", path, summary.ModelName)
	}
	contextSuffix += " (" + dateLabel + ")"

	return w.Write([]string{
		// path,
		summary.ModelName + contextSuffix,
		strconv.Itoa(summary.ChecksPassed),
		strconv.Itoa(summary.ChecksError),
		strconv.Itoa(summary.ChecksTotal),
	})
}

func scanForJSONSummaryLine(f *os.File) ([]byte, error) {
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		jsonPart, found := strings.CutPrefix(line, jsonSummaryPrefix)
		if found {
			return []byte(jsonPart), nil
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return nil, errNoJSONSummaryLine
}

const jsonSummaryPrefix = "JSON Summary: "

var errNoJSONSummaryLine = errors.New("no JSON Summary line found")

type jsonSummary struct {
	ChecksPassed int    `json:"checksPassed"`
	ChecksError  int    `json:"checksError"`
	ChecksTotal  int    `json:"checksTotal"`
	ModelName    string `json:"modelName"`
}
