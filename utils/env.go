package utils

import (
	"fmt"
	"os"
	"strings"
)

// ReadEnvFile reads an environment file in KEY=VALUE format
// Lines starting with # are treated as comments and ignored
// Empty lines are also ignored
func ReadEnvFile(filePath string) (map[string]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		// If file doesn't exist, return empty map (not an error)
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	envMap := make(map[string]string)
	lines := strings.Split(string(content), "\n")

	for lineNum, line := range lines {
		// Trim whitespace
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first '=' only
		parts := splitOnce(line, '=')
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid env format at line %d: %s (expected KEY=VALUE)", lineNum+1, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("empty key at line %d", lineNum+1)
		}

		envMap[key] = value
	}

	return envMap, nil
}

// splitOnce splits a string on the first occurrence of sep
func splitOnce(s string, sep rune) []string {
	for i, c := range s {
		if c == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}
