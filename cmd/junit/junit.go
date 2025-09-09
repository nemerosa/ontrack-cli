package junit

import (
	"encoding/xml"
	"github.com/gobwas/glob"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func GetSummaryJUnitTestReports(pattern string) (int, int, int, error) {
	g := glob.MustCompile(pattern, '/')
	var matches []string

	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if g.Match(path) {
			matches = append(matches, path)
		}
		return nil
	})

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

type TestSuites struct {
	TestSuite []TestSuite `xml:"testsuite"`
}

// getSummaryJUnitTestReport parses the JUnit XML files denoted by the path and returns the number of tests passed, skipped and failed.
//
// A JUnit XML can be formatted differently depending on the library or tool used.
// Some XML contain a <testsuites> root element, which can contain multiple <testsuite> elements.
// In that case, this method returns the sum of all tests reported in each <testsuite>.
// Other XML directly contain a single <testsuite> root element.
func getSummaryJUnitTestReport(path string) (int, int, int, error) {
	reader, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, err
	}
	buf, err := io.ReadAll(reader)
	if err != nil {
		return 0, 0, 0, err
	}

	var root TestSuites
	err = xml.Unmarshal(buf, &root)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(root.TestSuite) == 0 {
		// The file possibly contains a single <testsuite> element
		var testsuiteRoot TestSuite
		err = xml.Unmarshal(buf, &testsuiteRoot)
		if err != nil {
			return 0, 0, 0, err
		}
		root.TestSuite = append(root.TestSuite, testsuiteRoot)
	}

	var passed, skipped, failures, errors int

	for _, testSuite := range root.TestSuite {
		passed += testSuite.Tests - testSuite.Skipped - testSuite.Failures - testSuite.Errors
		skipped += testSuite.Skipped
		failures += testSuite.Failures
		errors += testSuite.Errors
	}

	return passed, skipped, failures + errors, nil
}
