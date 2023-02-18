package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPrecedence(t *testing.T) {
	envPrefix := "TEST"
	configName := "test"

	// Run the tests in a temporary directory
	tmpDir, err := os.MkdirTemp("", "snakebite")
	if err != nil {
		t.Fatalf("failed to create a temporary test directory: %v", err)
	}

	testDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(testDir); err != nil {
			t.Fatalf("failed to change directory: %v", err)
		}
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	t.Run("Set favorite-color from the config file and the number from the default flag value", func(t *testing.T) {
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
				if err != nil {
					t.Fatalf("failed to read test config file: %v", err)
				}

				if err := os.WriteFile(filepath.Join(tmpDir, configFileTOML), configB, 0644); err != nil {
					t.Fatalf("failed to write test config file: %v", err)
				}
				defer os.Remove(filepath.Join(tmpDir, configFileTOML))

				// Run ./example
				cmd := newRootCommand(envPrefix, configName, tc.replaceHyphenWithCamelCase)
				output := &bytes.Buffer{}
				cmd.SetOut(output)
				if err := cmd.Execute(); err != nil {
					t.Fatalf("failed to execute command: %q", err)
				}

				got := output.String()
				want := `Your favorite color is: blue
The magic number is: 7
`
				if got != want {
					t.Errorf("want output %q, got %q", want, got)
				}
			})
		}
	})

	t.Run("Set favorite-color from the environment variable value and the number from the default flag value", func(t *testing.T) {
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

		got := output.String()
		want := `Your favorite color is: purple
The magic number is: 7
`
		if got != want {
			t.Errorf("want output %q, got %q", want, got)
		}
	})

	t.Run("Set number from the flag value and the color from the default flag value", func(t *testing.T) {
		// Run ./example --number 2
		cmd := newRootCommand(envPrefix, configName, false)
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"--number", "2"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("failed to execute command: %q", err)
		}

		got := output.String()
		want := `Your favorite color is: red
The magic number is: 2
`
		if got != want {
			t.Errorf("want output %q, got %q", want, got)
		}
	})
}
