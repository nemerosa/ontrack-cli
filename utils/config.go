package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func ExpandConfig(initial string) (string, error) {
	// Parse the initial YAML string
	var data interface{}
	if err := yaml.Unmarshal([]byte(initial), &data); err != nil {
		return "", fmt.Errorf("failed to parse initial YAML: %w", err)
	}

	// Expand all @path references
	expanded, err := expandNode(data, "")
	if err != nil {
		return "", err
	}

	// Marshal back to YAML string
	result, err := yaml.Marshal(expanded)
	if err != nil {
		return "", fmt.Errorf("failed to marshal expanded YAML: %w", err)
	}

	return string(result), nil
}

// expandNode recursively processes a YAML node and expands @path references
func expandNode(node interface{}, baseDir string) (interface{}, error) {
	switch v := node.(type) {
	case string:
		// Check if the string looks like @path
		if strings.HasPrefix(v, "@") {
			path := strings.TrimPrefix(v, "@")
			// Resolve path relative to baseDir if provided
			if baseDir != "" && !filepath.IsAbs(path) {
				path = filepath.Join(baseDir, path)
			}

			// Read the file at path
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", path, err)
			}

			// Parse the YAML content
			var parsedData interface{}
			if err := yaml.Unmarshal(content, &parsedData); err != nil {
				return nil, fmt.Errorf("failed to parse YAML from %s: %w", path, err)
			}

			// Recursively expand the parsed data (using the directory of the current file as base)
			newBaseDir := filepath.Dir(path)
			return expandNode(parsedData, newBaseDir)
		}
		return v, nil

	case map[interface{}]interface{}:
		// Process each key-value pair in the map
		result := make(map[interface{}]interface{})
		for key, value := range v {
			expandedValue, err := expandNode(value, baseDir)
			if err != nil {
				return nil, err
			}
			result[key] = expandedValue
		}
		return result, nil

	case []interface{}:
		// Process each element in the array
		result := make([]interface{}, len(v))
		for i, elem := range v {
			expandedElem, err := expandNode(elem, baseDir)
			if err != nil {
				return nil, err
			}
			result[i] = expandedElem
		}
		return result, nil

	default:
		// For other types (int, bool, nil, etc.), return as-is
		return v, nil
	}
}
