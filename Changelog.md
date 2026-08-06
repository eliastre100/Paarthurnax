# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-08-06

### Added

- The `init` command lets you choose the source locale when a project contains more than one locale.
- Add destination languages interactively with `paarthurnax locales add`. 
- Translate missing segments on language addition
- Show translation progress in a more transparent way
- Report the underlying DeepL error when a translation request fails.

### Changed

- Locale configuration is now saved in `.paarthurnax`. 
- Projects can use any supported source locale instead of relying on a fixed source language and destination-language list.
- The project is now under an AGPL-3 license 

### Fixed

- Remove every destination plural form when its source plural key is deleted.
- Retry DeepL 429 (too many requests) responses with backoff instead of failing immediately.
- No longer send an empty `Content-Type` header in DeepL requests.
- The plural translations no longer have a chance to fail the parameter sanity check by translating the count hint

## [0.1.0] - 2024-09-17

### Added

- init command initialize the repository with a .paarthurnax state file
- translate command translate all untranslated segments with deepl
- normalize command sort all destination files segments in order to limit noise after translation

[unreleased]: https://github.com/eliastre100/Paarthurnax/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/eliastre100/Paarthurnax/releases/tag/v1.0.0
[0.1.0]: https://github.com/eliastre100/Paarthurnax/releases/tag/v0.1.0
