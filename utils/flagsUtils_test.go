package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBuildIdFromEnv(t *testing.T) {
	// Not set
	_ = os.Unsetenv("YONTRACK_BUILD_ID")
	val, _ := GetBuildIdFromEnv()
	assert.Equal(t, 0, val)

	// Set to empty
	_ = os.Setenv("YONTRACK_BUILD_ID", "")
	val, _ = GetBuildIdFromEnv()
	assert.Equal(t, 0, val)

	// Set to non-integer
	_ = os.Setenv("YONTRACK_BUILD_ID", "abc")
	val, _ = GetBuildIdFromEnv()
	assert.Equal(t, 0, val)

	// Set to integer
	_ = os.Setenv("YONTRACK_BUILD_ID", "123")
	val, _ = GetBuildIdFromEnv()
	assert.Equal(t, 123, val)
}
