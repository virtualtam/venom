// Package venom provides helper functions to use the cobra and viper libraries
// in cunjunction to load application configuration from command-line flags
// (cobra), environment variables (viper) and configuration file(s) (viper).
package venom

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Inject creates a new viper.Viper and binds environment variables and configuration file
// settings to command flags.
func Inject(cmd *cobra.Command, envPrefix string, configPaths []string, configName string, replaceHyphenWithCamelCase bool) error {
	v := viper.New()
	return InjectTo(v, cmd, envPrefix, configPaths, configName, replaceHyphenWithCamelCase)
}

// Inject binds environment variables and configuration file settings to
// command flags using an existing viper.Viper.
func InjectTo(v *viper.Viper, cmd *cobra.Command, envPrefix string, configPaths []string, configName string, replaceHyphenWithCamelCase bool) error {
	// Set as many paths as you like where viper should look for the
	// config file.
	for _, configPath := range configPaths {
		v.AddConfigPath(configPath)
	}

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(configName)

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable MYAPP_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(envPrefix)

	// Environment variables can't have dashes in them, so bind them to their equivalent
	// keys with underscores, e.g. --favorite-color to MYAPP_FAVORITE_COLOR
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	if err := bindFlags(cmd, v, replaceHyphenWithCamelCase); err != nil {
		return err
	}

	return nil
}

// bindFlags binds each command flag to its associated Viper configuration key
// and environment variable.
func bindFlags(cmd *cobra.Command, v *viper.Viper, replaceHyphenWithCamelCase bool) error {
	var flagErrors error

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name
		// If using camelCase in the config file, replace hyphens with a camelCased string.
		// Since viper does case-insensitive comparisons, we don't need to bother fixing the case, and only need to remove the hyphens.
		if replaceHyphenWithCamelCase {
			configName = strings.ReplaceAll(f.Name, "-", "")
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				flagErrors = errors.Join(flagErrors, err)
			}
		}
	})

	if flagErrors != nil {
		return fmt.Errorf("failed to bind command flags: %v", flagErrors)
	}

	return nil
}
