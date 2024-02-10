package junit

import (
	"testing"
)

func TestJUnitParsingSimple(t *testing.T) {

	pattern := "junit_reports/junit_report_simple.xml"

	passed, skipped, failed, err := GetSummaryJUnitTestReports(pattern)
	if err != nil {
		t.Errorf("Error reading the JUnit XML reports: %v", err)
	}

	if passed != 1 {
		t.Errorf("Passed - Expected: 1, Actual: %v", passed)
	}

	if skipped != 0 {
		t.Errorf("Skipped - Expected: 0, Actual: %v", skipped)
	}

	if failed != 0 {
		t.Errorf("Failed - Expected: 0, Actual: %v", failed)
	}
}

func TestJUnitParsingComposite(t *testing.T) {

	pattern := "junit_reports/junit_report_composite.xml"

	passed, skipped, failed, err := GetSummaryJUnitTestReports(pattern)
	if err != nil {
		t.Errorf("Error reading the JUnit XML reports: %v", err)
	}

	if passed != 2 {
		t.Errorf("Passed - Expected: 2, Actual: %v", passed)
	}

	if skipped != 1 {
		t.Errorf("Skipped - Expected: 1, Actual: %v", skipped)
	}

	if failed != 1 {
		t.Errorf("Failed - Expected: 1, Actual: %v", failed)
	}

}

func TestJunitParsingTestSuitesSingleSuite(t *testing.T) {
	pattern := "junit_reports/junit_report_simple_testsuites_single_suite.xml"

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

func TestJunitParsingTestSuitesMultipleSuites(t *testing.T) {
	pattern := "junit_reports/junit_report_simple_testsuites_multiple_suites.xml"

	passed, skipped, failed, err := GetSummaryJUnitTestReports(pattern)
	if err != nil {
		t.Errorf("Error reading the JUnit XML reports: %v", err)
	}

	if passed != 5 {
		t.Errorf("Passed - Expected: 5, Actual: %v", passed)
	}

	if skipped != 2 {
		t.Errorf("Skipped - Expected: 2, Actual: %v", skipped)
	}

	if failed != 2 {
		t.Errorf("Failed - Expected: 2, Actual: %v", failed)
	}
}

func TestJUnitParsingGlob(t *testing.T) {

	pattern := "junit_reports/glob/*.xml"

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
