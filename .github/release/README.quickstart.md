# DirAgent Quick Start

This package is self-contained. You do not need to clone the repository.

## What is included

- `diragent` or `diragent.exe`
- `install` entrypoint (one-click)
- `uninstall` entrypoint (one-click)
- `scripts/install.*`
- `scripts/uninstall.*`

## Install in 3 steps

1. Unzip this package into a folder where you want DirAgent installed.
2. Run the install entrypoint from this folder:
   - Windows: double-click `install.bat`
   - macOS / Linux: `./install.sh`
3. Right-click any directory and choose:
   - `Open in Codex (DirAgent)`
   - `Open in Claude Code (DirAgent)`

## Config and data locations

- Config: `<install-folder>/config.toml`
- Assets: `<install-folder>/data/assets`

## Uninstall

- Windows: double-click `uninstall.bat`
- macOS / Linux: `./uninstall.sh`
