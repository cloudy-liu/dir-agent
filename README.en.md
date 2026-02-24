# 馃殌 DirAgent

> One-click `Codex / Claude` launch from your file manager, with automatic directory switching.

馃寪 **Language**: [English](README.en.md) | [涓枃](README.md)


## 馃搶 Table of Contents

- [鉁?Overview](#-overview)
- [馃幆 Key Features](#-key-features)
- [鈿?Quick Start (Windows Recommended)](#-quick-start-windows-recommended)
- [馃洜锔?Installation (Command Line)](#锔?installation-command-line)
- [鈿欙笍 Configuration (`config.toml`)](#锔?configuration-configtoml)
- [馃攳 Argument Precedence](#-argument-precedence)
- [馃И Build & Verification](#-build--verification)
- [馃Н Troubleshooting](#-troubleshooting)
- [馃摝 Assets & Paths](#-assets--paths)


## 鉁?Overview

`DirAgent` adds file-manager context-menu entries:

- `Open in Codex (DirAgent)`
- `Open in Claude (DirAgent)`

Behavior:

- **Directory / directory background selected** 鈫?show context menu and launch inside that directory
- **File selected** 鈫?context menu is hidden by design


## 馃幆 Key Features

- 馃柋锔?Right-click launch for Codex / Claude
- 馃幆 Directory-only context menu scope (avoid file-action ambiguity)
- 馃獰 Windows menu icons (`.ico`, white background)
- 馃攣 Terminal strategy control (`tab_preferred` / `new_window`)
- 馃З Configurable terminal preference, CLI path, and default args


## 鈿?Quick Start (Windows Recommended)

Double-click these scripts (no manual arguments):

1. `scripts/diragent-1-build-and-verify.bat`  
   - runs `go test ./...`  
   - builds `diragent.exe`

2. `scripts/diragent-2-install-right-click.bat`  
   - auto-builds `diragent.exe` if missing  
   - installs Explorer context menu and icons

3. `scripts/diragent-3-uninstall-right-click.bat`  
   - removes context-menu entries  
   - removes extracted assets and config


## 馃洜锔?Installation (Command Line)

### Windows

Prerequisite: `diragent.exe` exists in repo root, or `diragent` is available in `PATH`.

```powershell
# Install
.\scripts\install.ps1

# Uninstall
.\scripts\uninstall.ps1

# Uninstall and clean assets + config
.\scripts\uninstall.ps1 -RemoveAssets -RemoveConfig
```

### macOS / Linux

```bash
chmod +x ./scripts/install.sh ./scripts/uninstall.sh
./scripts/install.sh ./diragent
./scripts/uninstall.sh ./diragent
```

> On macOS, it creates:
> - `~/Applications/DirAgent/Open in Codex (DirAgent).app`
> - `~/Applications/DirAgent/Open in Claude (DirAgent).app`


## 鈿欙笍 Configuration (`config.toml`)

Config file path:

- Windows: `%AppData%\dir-agent\config.toml`
- macOS/Linux: `~/.config/dir-agent/config.toml`

Default config:

```toml
[terminals]
preferred = ""

[terminals.windows_terminal]
profile = ""
shell = "powershell"

[tools.codex]
command = "codex"
default_args = ["--dangerously-bypass-approvals-and-sandbox"]

[tools.claude]
command = "claude"
default_args = ["--dangerously-skip-permissions"]

[behavior]
resolve_file_to_parent = true
open_mode = "tab_preferred"
```

### 馃搵 Full Parameter Reference

| Key | Type | Default | What it does | When to change |
|---|---|---|---|---|
| `terminals.preferred` | `string` | `""` | Preferred terminal; empty means fallback chain | Multiple terminals installed; need deterministic selection |
| `terminals.windows_terminal.profile` | `string` | `""` | Windows Terminal profile name (for example: `Cmder`, `PowerShell`, `Command Prompt`) | Use when preferred terminal is `windows-terminal` and you want a specific tab profile |
| `terminals.windows_terminal.shell` | `string` | `"powershell"` | Runner shell used to execute `codex`/`claude` inside Windows Terminal (`powershell` or `cmd`) | Set `cmd` to better align with Cmd/Cmder workflows |
| `tools.codex.command` | `string` | `"codex"` | Codex command name or absolute path | `codex` missing in PATH / custom command path |
| `tools.codex.default_args` | `string[]` | `["--dangerously-bypass-approvals-and-sandbox"]` | Default args for every Codex launch | Change only if you do not want full-access defaults |
| `tools.claude.command` | `string` | `"claude"` | Claude command name or absolute path | `claude` missing in PATH / custom command path |
| `tools.claude.default_args` | `string[]` | `["--dangerously-skip-permissions"]` | Default args for every Claude launch | Change only if you do not want full-access defaults |
| `behavior.resolve_file_to_parent` | `bool` | `true` | Convert file path to parent folder when using CLI path input | Keep `true` unless you need strict path-type behavior |
| `behavior.open_mode` | `string` | `"tab_preferred"` | Controls tab/window behavior | See mode details below |

### 馃 `open_mode` Details

- `tab_preferred` (default)  
  Reuse current terminal window with a new tab when possible; fallback to new window otherwise.

- `new_window`  
  Always open a new window.

- Any other value  
  Treated as invalid and falls back to `tab_preferred`.

### 馃Л Common `terminals.preferred` values

- Windows: `windows-terminal` / `wezterm` / `powershell`
- macOS: `terminal.app` / `iterm2`
- Linux: `x-terminal-emulator` / `gnome-terminal` / `konsole` / `xterm`

### Windows Terminal profile/shell examples

```toml
[terminals]
preferred = "windows-terminal"

[terminals.windows_terminal]
profile = "Cmder"
shell = "cmd"
```


## 馃攳 Argument Precedence

Merge order (low 鈫?high):

1. Built-in defaults  
2. `default_args` from `config.toml`  
3. Passthrough args after `--`


## 馃И Build & Verification

### Build

```powershell
# Windows
go build -o diragent.exe ./cmd/diragent
```

```bash
# macOS/Linux
go build -o diragent ./cmd/diragent
```

### Test

```bash
go test ./...
```

### Suggested Windows acceptance flow

1. Double-click `scripts/diragent-1-build-and-verify.bat`
2. Double-click `scripts/diragent-2-install-right-click.bat`
3. Verify manually:
   - folder 鈫?`Open in Codex (DirAgent)`
   - file 鈫?no DirAgent menu item
   - Chinese/space paths
   - icon visibility
4. Double-click `scripts/diragent-3-uninstall-right-click.bat` to verify rollback


## 馃Н Troubleshooting

### 1) Error `2147942402 (0x80070002)` when launching Codex

Meaning: command not found.  
Check in order:

1. Run `Get-Command codex` in PowerShell
2. If missing, set `tools.codex.command` in `config.toml`
3. Re-run `scripts/diragent-2-install-right-click.bat`

### 2) Context menu installed but not visible

- Refresh Explorer (`F5`)
- Or restart Explorer
- Ensure install happened under current-user scope (`HKCU`)

### 3) Not opening in tab as expected

- Confirm `behavior.open_mode = "tab_preferred"`
- If terminal cannot reuse tabs, fallback may open a new window


## 馃摝 Assets & Paths

Icons are embedded via `go:embed` and extracted during install:

- Windows: `.ico`
- macOS/Linux: `.png`

Asset paths:

- Windows: `%AppData%\dir-agent\assets`
- macOS/Linux: `~/.local/share/dir-agent/assets`
