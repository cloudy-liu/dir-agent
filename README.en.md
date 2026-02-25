# DirAgent

Launch `Codex` or `Claude Code` directly from your file manager in the target directory.

Language: [English](README.en.md) | [中文](README.md)

## Value

DirAgent removes the repeated flow:

`open terminal -> cd into folder -> run codex/claude`

After installation, just right-click a directory:

- `Open in Codex (DirAgent)`
- `Open in Claude Code (DirAgent)`

## Download the right release asset

Always download from **Release -> Assets**.  
Do not use `Source code (zip/tar.gz)`.

Pick one package by OS/arch:

- Windows x64: `diragent_<tag>_windows_amd64.zip`
- Windows ARM64: `diragent_<tag>_windows_arm64.zip`
- macOS Intel: `diragent_<tag>_darwin_amd64.zip`
- macOS Apple Silicon: `diragent_<tag>_darwin_arm64.zip`
- Linux x64: `diragent_<tag>_linux_amd64.zip`
- Linux ARM64: `diragent_<tag>_linux_arm64.zip`

Each package already includes binary + install/uninstall scripts.  
No repository clone is required.

## Install in 3 steps

1. Unzip package to your preferred install folder.
2. Run install script in that folder:

Windows (PowerShell):

```powershell
.\scripts\install.ps1 -BinaryPath .\diragent.exe
```

macOS / Linux:

```bash
chmod +x ./diragent ./scripts/install.sh ./scripts/uninstall.sh
./scripts/install.sh ./diragent
```

3. Right-click any directory and launch Codex/Claude from the menu.

## Config and data paths

- Config: `<install-folder>/config.toml`
- Assets: `<install-folder>/data/assets`

## Quick troubleshooting

- `0x80070002` or command not found:
  set `tools.codex.command` or `tools.claude.command` in `config.toml` to an absolute executable path.
- Menu not visible:
  refresh file manager (`F5`) or restart Explorer/Finder.
- WezTerm tab behavior:
  set `terminals.preferred = "wezterm"` and `behavior.open_mode = "tab_preferred"`.

## Development

```bash
go test ./...
```

```bash
go build -o diragent ./cmd/diragent
```

MIT License. See `LICENSE`.
