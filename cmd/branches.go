package cmd

import (
	"regexp"
)

func NormalizeBranchName(name string) string {
	re := regexp.MustCompile("[^A-Za-z0-9\\._-]")
	return re.ReplaceAllString(name, "-")
}
