# `venom` - Configuration management helper for `cobra` and `viper` 🐍
Venom provides helper functions to use the [cobra](https://github.com/spf13/cobra)
and [viper](https://github.com/spf13/viper) libraries in cunjunction to load
application configuration from:

- command-line flags (`cobra`),
- environment variables (`viper`),
- configuration file(s) (`viper`).

## Example
An example Cobra application can be found under [example/main.go](./example/main.go).

This example demonstrates how Venom can be used to load configuration variables:

```shell
$ go run ./example
Your favorite color is: red
The magic number is: 7

$ go run ./example -n 12
Your favorite color is: red
The magic number is: 12

$ VENOMOUS_FAVORITE_COLOR=purple go run ./example
Your favorite color is: purple
The magic number is: 7
```

## Credits
These helpers have been adapted from the article
[Sting of the Viper: Getting Cobra and Viper to work together](https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/)
by Carolyn Van Slyck,
and the corresponding repository, [carolynvs/stingoftheviper](https://github.com/carolynvs/stingoftheviper)
so they can be used as a library.

The original article demonstrates how to integrate `spf13/cobra` with `spf13/viper` such that:

- command-line flags have the highest precedence,
- then environment variables,
- then config file values,
- and then defaults set on command-line flags.

## License
Venom is licensed under the MIT license.
