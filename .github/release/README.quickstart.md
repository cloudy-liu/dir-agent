# DirAgent Quick Start

This package is self-contained. You do not need to clone the repository.

## What is included

- `diragent` or `diragent.exe`
- `scripts/install.*`
- `scripts/uninstall.*`

## Install in 3 steps

1. Unzip this package into a folder where you want DirAgent installed.
2. Run the install script from this folder:
   - Windows PowerShell:
     `.\scripts\install.ps1 -BinaryPath .\diragent.exe`
   - macOS / Linux:
     `chmod +x ./diragent ./scripts/install.sh ./scripts/uninstall.sh && ./scripts/install.sh ./diragent`
3. Right-click any directory and choose:
   - `Open in Codex (DirAgent)`
   - `Open in Claude Code (DirAgent)`

## Config and data locations

- Config: `<install-folder>/config.toml`
- Assets: `<install-folder>/data/assets`

## Uninstall

- Windows:
  `.\scripts\uninstall.ps1 -BinaryPath .\diragent.exe -RemoveAssets -RemoveConfig`
- macOS / Linux:
  `./scripts/uninstall.sh ./diragent`
