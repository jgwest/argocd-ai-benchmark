package types

import (
	"fmt"
	"os"
	"path/filepath"
)

type ResultReporter interface {
	// Init once before the reporter is first used
	Init(modelName string) error

	// ReportIndividualResult will be called to output text related to a single result
	ReportIndividualResult(str string)

	// ReportCumulativeResult will be called to output text related to a multiple results, often for cumulative/aggregate summary
	ReportCumulativeResult(str string)

	// IncludeJsonSummary should return true if json summary of the aggregate result should be output, false otherwise
	IncludeJsonSummary() bool

	// Close is called once, when the model tests have completed
	Close(modelName string) error
}

type FileResultReporter struct {
	output         outBuffer
	fileOutputPath string
}

func NewFileResultReporter(fileOutputPath string) *FileResultReporter {
	return &FileResultReporter{
		fileOutputPath: fileOutputPath,
	}
}

var _ ResultReporter = &FileResultReporter{}

func (frr *FileResultReporter) Init(modelName string) error {
	// Only log simple status text to console for the file reporter
	fmt.Println("* Running model '"+modelName+"', to output path:", frr.fileOutputPath)
	return nil
}

func (frr *FileResultReporter) ReportIndividualResult(str string) {
	frr.output.println(str)

}

func (frr *FileResultReporter) ReportCumulativeResult(str string) {
	frr.output.println(str)
}

func (frr *FileResultReporter) Close(modelName string) error {

	parentPath := filepath.Dir(frr.fileOutputPath)
	if parentPath == "" {
		return fmt.Errorf("unexpected empty parent path")
	}

	if err := os.MkdirAll(parentPath, 0755); err != nil {
		return fmt.Errorf("unable to mkdirall: %v", err)
	}

	if err := os.WriteFile(frr.fileOutputPath, []byte(frr.output.out), 0644); err != nil {
		return fmt.Errorf("unable to write to file '%s': %v", frr.fileOutputPath, err)
	}

	fmt.Println("* Completed run of model '"+modelName+"', to output path:", frr.fileOutputPath)

	return nil
}

func (frr *FileResultReporter) IncludeJsonSummary() bool {
	return true
}

type StdoutResultReporter struct {
}

var _ ResultReporter = &StdoutResultReporter{}

func (frr *StdoutResultReporter) ReportIndividualResult(str string) {
	fmt.Println(str)

}

func (frr *StdoutResultReporter) ReportCumulativeResult(str string) {
	fmt.Println(str)
}

func (frr *StdoutResultReporter) Close(modelName string) error {
	return nil
}

func (frr *StdoutResultReporter) Init(modelName string) error {
	return nil
}
func (frr *StdoutResultReporter) IncludeJsonSummary() bool {
	return false
}
