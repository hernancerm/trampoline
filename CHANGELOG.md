# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-11-11

### Added

- Prebuilt binaries for macOS arm64, Linux arm64 and x86_64, and Windows x86_64,
  published on release and installable with mise.
- `TR_CMD` to run the core from source with babashka.
- PowerShell support.
- Config files `~/.trampoline` and `~/.trampoline.local`, merged in order. Either may be
  missing.
- `:0` suffix to list a directory itself instead of its sub-dirs.
- `$VAR` and `${VAR}` expansion in the config.
- `TR_SET_KEYBINDS=0` to bind the widget yourself.

### Changed

- Renamed the project from `zsh-trampoline` to `trampoline`, and its widget and variables
  to the `tr` prefix.
- The core is a native binary, so babashka is no longer needed.
- The default keybind is <kbd>Ctrl-t</kbd>.
- The config moved from the `ZT_CONFIG` parameter to the config files.

### Removed

- The `*` toggle in fzf to list only the configured items.

### Fixed

- `vim` is the editor when `$EDITOR` is not set.
- Paths with whitespace.
- A duplicate prompt with a long item list.

## [0.1.0] - 2024-09-28

### Added

- Initial release.

[Unreleased]: https://github.com/hernancerm/trampoline/compare/0.2.0...HEAD
[0.2.0]: https://github.com/hernancerm/trampoline/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/hernancerm/trampoline/releases/tag/0.1.0
