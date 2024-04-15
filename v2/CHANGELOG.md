# Change Log

All notable changes in pubsub-go will be reflected in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v2.0.1 - 2024-04-15

### Changed v2.0.1

- Updated documentation for more clarity

## v2.0.0 - 2024-03-28

### Added v2.0.0

- Created `pkg` directory for organizing the codebase
- Added `CreateSubscription` method to allow creation of a single subscription
- Added `Close` method to allow consumers of the library to close the PubSub client

### Changed v2.0.0

- Changed `CreateSubscriptions` method signature to be variadic
- Changed `CreateTopic` method signature to be variadic
- Changed `NewPubsub` method name to `NewPubSub`
- Changed `Publish` method signature to accept the object to publish (marshals provided interface as JSON when not a `[]byte`)

### Removed v2.0.0

- Removed `IsLocal` property from configuration
- Removed `ServiceAccountFilePath` property from configuration

## v1.1.0 - 2023-08-02

### Changed v1.1.0

- Enhanced the `Receive` method to accept options for controlling subscription behavior

## v1.0.2 - 2023-04-14

### Changed v1.0.2

- Updated dependencies

## v1.0.1 - 2022-08-26

### Added v1.0.1

- Added CreateSubscriptions method to create a new subscription for a topic
- Added CreateTopic method to create a new topic

## v1.0.0 - 2022-06-22

### Added v1.0.0

- Initial release of the library
