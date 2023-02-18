package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrecedence(t *testing.T) {
	envPrefix := "TEST"
	configName := "test"

	// Run the tests in a temporary directory
	tmpDir, err := os.MkdirTemp("", "snakebite")
	require.NoError(t, err, "error creating a temporary test directory")

	testDir, err := os.Getwd()
	require.NoError(t, err, "error getting the current working directory")

	defer func() {
		if err := os.Chdir(testDir); err != nil {
			t.Fatalf("failed to change directory: %q", err)
		}
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err, "error changing to the temporary test directory")

	t.Run("Set favorite-color with the config file", func(t *testing.T) {
		testcases := []struct {
			name                       string
			configFile                 string
			replaceHyphenWithCamelCase bool
		}{
			{name: "hyphen", configFile: "testdata/config-hyphen.toml"},
			{name: "camelCase", configFile: "testdata/config-camel.toml", replaceHyphenWithCamelCase: true},
		}

		configFileTOML := fmt.Sprintf("%s.toml", configName)

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				// Copy the config file into our temporary test directory
				configB, err := os.ReadFile(filepath.Join(testDir, tc.configFile))
				require.NoError(t, err, "error reading test config file")

				err = os.WriteFile(filepath.Join(tmpDir, configFileTOML), configB, 0644)
				require.NoError(t, err, "error writing test config file")
				defer os.Remove(filepath.Join(tmpDir, configFileTOML))

				// Run ./example
				cmd := newRootCommand(envPrefix, configName, tc.replaceHyphenWithCamelCase)
				output := &bytes.Buffer{}
				cmd.SetOut(output)
				if err := cmd.Execute(); err != nil {
					t.Fatalf("failed to execute command: %q", err)
				}

				gotOutput := output.String()
				wantOutput := `Your favorite color is: blue
The magic number is: 7
`
				assert.Equal(t, wantOutput, gotOutput, "expected the color from the config file and the number from the flag default")
			})
		}
	})

	t.Run("Set favorite-color with an environment variable", func(t *testing.T) {
		// Run TEST_FAVORITE_COLOR=purple ./example
		colorEnvVar := fmt.Sprintf("%s_FAVORITE_COLOR", envPrefix)

		os.Setenv(colorEnvVar, "purple")
		defer os.Unsetenv(colorEnvVar)

		cmd := newRootCommand(envPrefix, configName, false)
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("failed to execute command: %q", err)
		}

		gotOutput := output.String()
		wantOutput := `Your favorite color is: purple
The magic number is: 7
`
		assert.Equal(t, wantOutput, gotOutput, "expected the color to use the environment variable value and the number to use the flag default")
	})

	t.Run("Set number with a flag", func(t *testing.T) {
		// Run ./example --number 2
		cmd := newRootCommand(envPrefix, configName, false)
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"--number", "2"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("failed to execute command: %q", err)
		}

		gotOutput := output.String()
		wantOutput := `Your favorite color is: red
The magic number is: 2
`
		assert.Equal(t, wantOutput, gotOutput, "expected the number to use the flag value and the color to use the flag default")
	})
}
