package main

// TestDocument represents the root structure of the YAML document
type TestDocument struct {
	Name   string   `yaml:"name"`
	Labels []string `yaml:"labels"`
	Checks []Check  `yaml:"checks"`
}

// Check represents a single check item with question and answers
type Check struct {
	Question string   `yaml:"question"`
	Answers  []string `yaml:"answers"`
}
