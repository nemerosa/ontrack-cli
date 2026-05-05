package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSearchFilter_ByName(t *testing.T) {
	filter := buildSearchFilter("mybuild", "")

	assert.Equal(t, 1, filter["maximumCount"])
	assert.Equal(t, "mybuild", filter["buildName"])
	assert.Equal(t, true, filter["buildExactMatch"])
	assert.NotContains(t, filter, "property")
	assert.NotContains(t, filter, "propertyValue")
}

func TestBuildSearchFilter_ByVersion(t *testing.T) {
	filter := buildSearchFilter("", "1.2.3")

	assert.Equal(t, 1, filter["maximumCount"])
	assert.Equal(t, releasePropertyType, filter["property"])
	assert.Equal(t, "1.2.3", filter["propertyValue"])
	// Regression: must use "property", not "propertyName" (issue #37)
	assert.NotContains(t, filter, "propertyName")
	assert.NotContains(t, filter, "buildName")
	assert.NotContains(t, filter, "buildExactMatch")
}
