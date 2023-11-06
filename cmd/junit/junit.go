package junit

import (
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
)

func GetSummaryJUnitTestReports(pattern string) (int, int, int, error) {
	// Gets all the matching fields
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0, 0, 0, err
	}

	// Counts
	totalPassed := 0
	totalSkipped := 0
	totalFailed := 0

	// Getting over all matches
	for _, match := range matches {
		passed, skipped, failed, err := getSummaryJUnitTestReport(match)
		if err != nil {
			return 0, 0, 0, err
		}
		totalPassed += passed
		totalSkipped += skipped
		totalFailed += failed
	}

	// OK
	return totalPassed, totalSkipped, totalFailed, nil
}

type TestSuite struct {
	Tests    int `xml:"tests,attr"`
	Skipped  int `xml:"skipped,attr"`
	Failures int `xml:"failures,attr"`
	Errors   int `xml:"errors,attr"`
}

func getSummaryJUnitTestReport(path string) (int, int, int, error) {
	reader, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, err
	}
	buf, err := io.ReadAll(reader)
	if err != nil {
		return 0, 0, 0, err
	}

	var root TestSuite
	err = xml.Unmarshal(buf, &root)
	if err != nil {
		return 0, 0, 0, err
	}

	passed := root.Tests - root.Skipped - root.Failures - root.Errors

	return passed, root.Skipped, root.Failures + root.Errors, nil
}
