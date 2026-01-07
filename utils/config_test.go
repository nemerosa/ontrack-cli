package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandConfig_SimpleYAML(t *testing.T) {
	// Test that simple YAML without @path references is unchanged
	input := `
name: test-project
version: 1.0.0
enabled: true
`
	result, err := ExpandConfig(input)
	require.NoError(t, err)
	assert.Contains(t, result, "name: test-project")
	assert.Contains(t, result, "version: 1.0.0")
	assert.Contains(t, result, "enabled: true")
}

func TestExpandConfig_FileReference(t *testing.T) {
	// Create a temporary directory and file
	tmpDir := t.TempDir()
	externalFile := filepath.Join(tmpDir, "external.yaml")

	externalContent := `
external: true
data: from-file
`
	err := os.WriteFile(externalFile, []byte(externalContent), 0644)
	require.NoError(t, err)

	// Create config that references the external file
	input := `
name: main-config
include: '@` + externalFile + `'`

	result, err := ExpandConfig(input)
	require.NoError(t, err)
	assert.Contains(t, result, "name: main-config")
	assert.Contains(t, result, "external: true")
	assert.Contains(t, result, "data: from-file")
}

func TestExpandConfig_FileNotFound(t *testing.T) {
	// Test that referencing a non-existent file returns an error
	input := `
name: main-config
include: '@/nonexistent/path/file.yaml'
`

	_, err := ExpandConfig(input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}

func TestExpandConfig_NestedFileReferences(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create the innermost file
	innerFile := filepath.Join(tmpDir, "inner.yaml")
	innerContent := `
level: 3
value: innermost
`
	err := os.WriteFile(innerFile, []byte(innerContent), 0644)
	require.NoError(t, err)

	// Create the middle file that references the inner file
	middleFile := filepath.Join(tmpDir, "middle.yaml")
	middleContent := `
level: 2
nested: '@inner.yaml'
`
	err = os.WriteFile(middleFile, []byte(middleContent), 0644)
	require.NoError(t, err)

	// Create the main config that references the middle file
	input := `
level: 1
config: '@` + middleFile + `'`

	result, err := ExpandConfig(input)
	require.NoError(t, err)
	assert.Contains(t, result, "level: 1")
	assert.Contains(t, result, "level: 2")
	assert.Contains(t, result, "level: 3")
	assert.Contains(t, result, "value: innermost")
}

func TestExpandConfig_ArrayWithReferences(t *testing.T) {
	// Create temporary files
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "item1.yaml")
	err := os.WriteFile(file1, []byte("name: item1\nvalue: 1"), 0644)
	require.NoError(t, err)

	file2 := filepath.Join(tmpDir, "item2.yaml")
	err = os.WriteFile(file2, []byte("name: item2\nvalue: 2"), 0644)
	require.NoError(t, err)

	// Create config with array containing file references
	input := `
items:
  - '@` + file1 + `'
  - '@` + file2 + `'
  - name: item3
    value: 3
`

	result, err := ExpandConfig(input)
	require.NoError(t, err)
	assert.Contains(t, result, "name: item1")
	assert.Contains(t, result, "value: 1")
	assert.Contains(t, result, "name: item2")
	assert.Contains(t, result, "value: 2")
	assert.Contains(t, result, "name: item3")
	assert.Contains(t, result, "value: 3")
}

func TestExpandConfig_MapWithReferences(t *testing.T) {
	// Create temporary files
	tmpDir := t.TempDir()

	dbFile := filepath.Join(tmpDir, "db.yaml")
	err := os.WriteFile(dbFile, []byte("host: localhost\nport: 5432"), 0644)
	require.NoError(t, err)

	cacheFile := filepath.Join(tmpDir, "cache.yaml")
	err = os.WriteFile(cacheFile, []byte("host: redis\nttl: 3600"), 0644)
	require.NoError(t, err)

	// Create config with nested maps containing file references
	input := `
services:
  database: '@` + dbFile + `'
  cache: '@` + cacheFile + `'
  api:
    host: api.example.com
    port: 8080
`

	result, err := ExpandConfig(input)
	require.NoError(t, err)
	assert.Contains(t, result, "services:")
	assert.Contains(t, result, "host: localhost")
	assert.Contains(t, result, "port: 5432")
	assert.Contains(t, result, "host: redis")
	assert.Contains(t, result, "ttl: 3600")
	assert.Contains(t, result, "host: api.example.com")
	assert.Contains(t, result, "port: 8080")
}

func TestExpandConfig_StringStartingWithAtButNotReference(t *testing.T) {
	// Test that strings starting with @ but not being valid file paths are handled
	tmpDir := t.TempDir()

	// Create a file that will be referenced
	validFile := filepath.Join(tmpDir, "valid.yaml")
	err := os.WriteFile(validFile, []byte("valid: true"), 0644)
	require.NoError(t, err)

	// Reference an invalid path - this should error
	input := `
name: test
twitter: '@username'
`

	// Since @username will be treated as a file path, it should fail to read
	_, err = ExpandConfig(input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}

func TestExpandConfig_RelativePathResolution(t *testing.T) {
	// Create a directory structure:
	// tmpDir/
	//   main.yaml
	//   subdir/
	//     config.yaml
	//     data.yaml

	tmpDir := t.TempDir()
	subdir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subdir, 0755)
	require.NoError(t, err)

	// Create data.yaml in subdir
	dataFile := filepath.Join(subdir, "data.yaml")
	err = os.WriteFile(dataFile, []byte("key: value"), 0644)
	require.NoError(t, err)

	// Create config.yaml in subdir that references data.yaml relatively
	configFile := filepath.Join(subdir, "config.yaml")
	configContent := `
section:
  data: '@data.yaml'
`
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// Create main config that references subdir/config.yaml
	input := `
main: true
include: '@` + configFile + `'
`

	result, err := ExpandConfig(input)
	require.NoError(t, err)
	assert.Contains(t, result, "main: true")
	assert.Contains(t, result, "key: value")
}

func TestExpandConfig_DifferentTypes(t *testing.T) {
	// Test that different YAML types are preserved
	input := `
string: hello
integer: 42
float: 3.14
boolean: true
null_value: null
list:
  - item1
  - item2
nested:
  key: value
`

	result, err := ExpandConfig(input)
	require.NoError(t, err)
	assert.Contains(t, result, "string: hello")
	assert.Contains(t, result, "integer: 42")
	assert.Contains(t, result, "float: 3.14")
	assert.Contains(t, result, "boolean: true")
	// null values might be omitted or shown as "null_value: null"
	assert.Contains(t, result, "- item1")
	assert.Contains(t, result, "- item2")
	assert.Contains(t, result, "key: value")
}

func TestExpandConfig_ComplexNestedStructure(t *testing.T) {
	// Create a complex structure with multiple levels and types
	tmpDir := t.TempDir()

	// Create validation config
	validationFile := filepath.Join(tmpDir, "validation.yaml")
	validationContent := `
type: junit
path: test-results.xml
`
	err := os.WriteFile(validationFile, []byte(validationContent), 0644)
	require.NoError(t, err)

	// Create promotion config
	promotionFile := filepath.Join(tmpDir, "promotion.yaml")
	promotionContent := `
name: RELEASE
validations:
  - '@validation.yaml'
promotions: []
`
	err = os.WriteFile(promotionFile, []byte(promotionContent), 0644)
	require.NoError(t, err)

	// Main config
	input := `
project:
  name: test-project
  branches:
    - name: main
      promotions:
        - '@` + promotionFile + `'
`

	result, err := ExpandConfig(input)
	require.NoError(t, err)
	assert.Contains(t, result, "project:")
	assert.Contains(t, result, "name: test-project")
	assert.Contains(t, result, "name: main")
	assert.Contains(t, result, "name: RELEASE")
	assert.Contains(t, result, "type: junit")
	assert.Contains(t, result, "path: test-results.xml")
}

func TestRenderConfig(t *testing.T) {
	vars := map[string]string{
		"version": "1.0.0",
		"name":    "test",
	}
	envMap := map[string]string{
		"HOME": "/home/user",
	}

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple variables",
			content:  "name: {{ vars \"name\" }}\nversion: {{ vars \"version\" }}",
			expected: "name: test\nversion: 1.0.0",
		},
		{
			name:     "default variable",
			content:  "other: {{ getvar \"other\" \"default\" }}",
			expected: "other: default",
		},
		{
			name:     "environment variable",
			content:  "home: {{ env \"HOME\" }}",
			expected: "home: /home/user",
		},
		{
			name:     "default environment variable",
			content:  "path: {{ getenv \"PATH\" \"/usr/bin\" }}",
			expected: "path: /usr/bin",
		},
		{
			name:     "sprig functions",
			content:  "upper: {{ \"hello\" | upper }}",
			expected: "upper: HELLO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderConfig(tt.content, vars, envMap)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
