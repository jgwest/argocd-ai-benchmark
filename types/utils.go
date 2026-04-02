package types

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
)

func downloadURL(url string) (string, error) {

	// Send an HTTP GET request to the URL.
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// Convert the body to a string and print it.
	content := string(body)
	return content, nil
}

// trimIndent removes common leading whitespace from each line of a string.
func trimIndent(s string) string {

	// Split the string into lines.
	lines := strings.Split(s, "\n")

	// Find the minimum indentation of non-empty lines.
	minIndent := math.MaxInt32
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue // Skip empty lines
		}

		indent := 0
		for _, r := range line {
			if r == ' ' || r == '\t' {
				indent++
			} else {
				break
			}
		}

		if indent < minIndent {
			minIndent = indent
		}
	}

	// If no indented lines were found, return the original string.
	if minIndent == math.MaxInt32 {
		return s
	}

	var trimmedLines []string
	for _, line := range lines {
		if len(line) > minIndent {
			trimmedLines = append(trimmedLines, line[minIndent:])
		} else {
			trimmedLines = append(trimmedLines, line)
		}
	}

	// Join the lines back together.
	return strings.Join(trimmedLines, "\n")
}
