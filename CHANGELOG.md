# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.3.0] - 2022-04-22

### Changed
- Updated to latest wavefront sdk patch release 0.9.11 for build dependancy
- Updated to latest plugin sdk version v0.16.0-alpha4  for build dependancy
- Updated to require go 1.18 for build dependancy

### Fixed
- Coerce metric timestamp sent to wavefront to 1 second precision using heuristic to determine Sensu metric timestamp precision. Needed to correctly handle prometheus metrics with millisecond precision ingested by agents.

## [0.2.2] - 2022-03-14

### Changed
- Updated go version to 1.17

## [0.2.1] - 2022-03-14

### Changed
- Updated release platform targets

## [0.2.0] - 2022-03-14

### Changed
- Update Sensu Go and SDK dependencies with the correct modules

## [0.1.0] - 2020-01-07

### Added
- Initial release
