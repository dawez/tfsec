package config_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aquasecurity/tfsec/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinRequiredVersionFromYAML(t *testing.T) {
	content := `
min_required_version: v1.2.0
`
	c := load(t, "config.yaml", content)

	assert.Equal(t, c.MinimumRequiredVersion, "v1.2.0")
}

func TestMinRequiredVersionFromJSON(t *testing.T) {
	content := `{
	"min_required_version": "v1.2.0"
}
`
	c := load(t, "config.json", content)

	assert.Equal(t, c.MinimumRequiredVersion, "v1.2.0")
}

func TestExcludesElementsFromYAML(t *testing.T) {
	content := `
severity_overrides:
  AWS018: LOW

exclude:
  - DP001

exclude_ignores:
  - DP002
`
	c := load(t, "config.yaml", content)

	assert.Contains(t, c.SeverityOverrides, "AWS018")
	assert.Contains(t, c.ExcludedChecks, "DP001")
	assert.Contains(t, c.ExcludeIgnores, "DP002")
}

func TestExcludesElementsFromYML(t *testing.T) {
	content := `
severity_overrides:
  AWS018: LOW

exclude:
  - DP001

exclude_ignores:
  - DP002
`
	c := load(t, "config.yml", content)

	assert.Contains(t, c.SeverityOverrides, "AWS018")
	assert.Contains(t, c.ExcludedChecks, "DP001")
	assert.Contains(t, c.ExcludeIgnores, "DP002")
}

func TestExcludesElementsFromJSON(t *testing.T) {
	content := `{
  "severity_overrides": {
    "AWS018": "LOW"
  },
  "exclude": [
    "DP001"
  ],
  "exclude_ignores": [
    "DP002"
  ]
}
`
	c := load(t, "config.json", content)

	assert.Contains(t, c.SeverityOverrides, "AWS018")
	assert.Contains(t, c.ExcludedChecks, "DP001")
	assert.Contains(t, c.ExcludeIgnores, "DP002")
}

func TestWarningIsRewrittenAsMedium(t *testing.T) {
	content := `{
  "severity_overrides": {
    "AWS018": "WARNING"
  },
  "exclude": [
    "DP001"
  ]
}
`
	c := load(t, "config.json", content)

	assert.Contains(t, c.SeverityOverrides, "AWS018")
	sev := c.SeverityOverrides["AWS018"]
	assert.Equal(t, "MEDIUM", sev)
}

func load(t *testing.T, filename, content string) *config.Config {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	configFileName := fmt.Sprintf("%s/%s", dir, filename)

	err = ioutil.WriteFile(configFileName, []byte(content), os.ModePerm)
	require.NoError(t, err)

	c, err := config.LoadConfig(configFileName)
	require.NoError(t, err)

	return c
}
