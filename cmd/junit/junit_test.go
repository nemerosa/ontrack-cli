package junit

import (
	"testing"
)

func TestJUnitParsing(t *testing.T) {

	pattern := "junit_reports/*.xml"

	passed, skipped, failed, err := GetSummaryJUnitTestReports(pattern)
	if err != nil {
		t.Errorf("Error reading the JUnit XML reports: %v", err)
	}

	if passed != 3 {
		t.Errorf("Passed - Expected: 3, Actual: %v", passed)
	}

	if skipped != 1 {
		t.Errorf("Skipped - Expected: 1, Actual: %v", skipped)
	}

	if failed != 1 {
		t.Errorf("Failed - Expected: 1, Actual: %v", failed)
	}

}
