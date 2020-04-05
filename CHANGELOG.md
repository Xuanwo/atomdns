# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v0.2.0] - 2020-04-05

### Added

- upstream: Add dot support

### Changed

- match: Remove debug log
- upstream: Set return code to server failure while meet error
- match/domain_list: Use builtin map instead of immutable radix tree

### Fixed

- match/domain_list: Fix subdomain not found
- deployment/systemd: Fix exec start typo

## v0.1.0 - 2020-04-05

### Added

- Add udp/tcp upstream
- Add domain list match
- Add simple dns server

[v0.2.0]: https://github.com/Xuanwo/atomdns/compare/v0.1.0...v0.2.0
