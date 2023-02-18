package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/virtualtam/venom"
)

// I'm declaring as vars so I can test easier, I recommend declaring these as constants
var (
	// The name of our config file, without the file extension because viper supports many different config file languages.
	defaultConfigFilename = "stingoftheviper"

	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to STING_NUMBER.
	envPrefix = "STING"

	// Replace hyphenated flag names with camelCase in the config file
	replaceHyphenWithCamelCase = false
)

func main() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Build the cobra command that handles our command line tool.
func NewRootCommand() *cobra.Command {
	// Store the result of binding cobra flags and viper config. In a
	// real application these would be data structures, most likely
	// custom structs per command. This is simplified for the demo app and is
	// not recommended that you use one-off variables. The point is that we
	// aren't retrieving the values directly from viper or flags, we read the values
	// from standard Go data structures.
	color := ""
	number := 0

	// Define our command
	rootCmd := &cobra.Command{
		Use:   "venom-example",
		Short: "Cober and Viper together at last",
		Long:  "Demonstrate how to get cobra flags to bind to viper properly",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return venom.Inject(cmd, envPrefix, defaultConfigFilename, replaceHyphenWithCamelCase)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			// Print the final resolved value from binding cobra flags and viper config
			fmt.Fprintln(out, "Your favorite color is:", color)
			fmt.Fprintln(out, "The magic number is:", number)
		},
	}

	// Define cobra flags, the default value has the lowest (least significant) precedence
	rootCmd.Flags().IntVarP(&number, "number", "n", 7, "What is the magic number?")
	rootCmd.Flags().StringVarP(&color, "favorite-color", "c", "red", "Should come from flag first, then env var STING_FAVORITE_COLOR then the config file, then the default last")

	return rootCmd
}
