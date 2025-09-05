# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Report the DeepL underlying error on translation request failure

### Fixed

- DeepL 529 response in case of too many requests does not fails and retry with backoff

## [0.1.0] - 2024-09-17

### Added

- init command initialize the repository with a .paarthurnax state file
- translate command translate all untranslated segments with deepl
- normalize command sort all destination files segments in order to limit noise after translation

[unreleased]: https://github.com/eliastre100/Paarthurnax/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/eliastre100/Paarthurnax/releases/tag/v0.1.0
