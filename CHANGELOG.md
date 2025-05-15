# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [unreleased]

## [v0.3.0] - 15-05-2025

### Added

- Printout of the docker network gateway IP on `orca status`
- A timer on the postgres instance being ready to accept connections before connecting

## [v0.2.0] - 12-05-2025

- Moved orca core deployement to github packages
- Integrated orca core package into CLI

## [v0.1.0] - 12-05-2025

### Added

- CLI & converted repo to a monorepo.
- Updated build stages.

## [v0.0.0] - 06-05-2025

### Added

- Initial implementation that accepts processor registrations, and can emit windows to processors.
- Can now handle results.
- Only the dependencies that the stage needs are passed in.

### Changed

### Removed

[unreleased]: https://github.com/Predixus/Orca/compare/v0.3.0...HEAD
[v0.3.0]: https://github.com/Predixus/Orca/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/Predixus/Orca/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/Predixus/Orca/compare/v0.0.0...v0.1.0
