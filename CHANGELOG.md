# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/) and this
project adheres to [Semantic Versioning](https://semver.org/).

## [v1.0.0](https://github.com/virtualtam/venom/releases/tag/v1.0.0) - 2023-02-18
Initial release.

### Added
- Allow setting multiple configuration paths for Viper to look for configuration files
- Run Go linters with golangci-lint

### Changed
- Refactor project as a package to be used as a library
- Refactor helper functions
- Cleanup example code and tests to avoid relying on global variables
- Bump the minimum Go version to 1.20 (to benefit from the new error helpers)
- Update documentation
- Update Github workflow to run linters, and run on tagged versions

### Fixed
- Handle all errors
