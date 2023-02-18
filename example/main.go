package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/virtualtam/venom"
)

func main() {
	// An array specifying in which directories Viper should look for configuration files.
	configPaths := []string{"."}

	// The name of our config file, without the file extension because Viper supports many different config file languages.
	configFilename := "venomous"

	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to VENOMOUS_NUMBER.
	envPrefix := "VENOMOUS"

	// Replace hyphenated flag names with camelCase in the config file
	replaceHyphenWithCamelCase := false

	cmd := newRootCommand(envPrefix, configPaths, configFilename, replaceHyphenWithCamelCase)
	cobra.CheckErr(cmd.Execute())
}

// newRootCommand initializes and returns an example Cobra command.
func newRootCommand(envPrefix string, configPaths []string, configName string, replaceHyphenWithCamelCase bool) *cobra.Command {
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
		Use:   "example",
		Short: "Cobra and Viper together at last",
		Long:  "Demonstrate how to get cobra flags to bind to viper properly",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return venom.Inject(cmd, envPrefix, configPaths, configName, replaceHyphenWithCamelCase)
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
	rootCmd.Flags().IntVarP(
		&number,
		"number",
		"n",
		7,
		"What is the magic number?",
	)
	rootCmd.Flags().StringVarP(
		&color,
		"favorite-color",
		"c",
		"red",
		fmt.Sprintf("Should come from flag first, then env var %s_FAVORITE_COLOR then the config file, then the default last", envPrefix),
	)

	return rootCmd
}
