# 🚀 DirAgent

> One-click `Codex / Claude` launch from your file manager, with automatic directory switching.

🌐 **Language**: [English](README.en.md) | [中文](README.md)

---

## 📌 Table of Contents

- [✨ Overview](#-overview)
- [🎯 Key Features](#-key-features)
- [⚡ Quick Start (Windows Recommended)](#-quick-start-windows-recommended)
- [🛠️ Installation (Command Line)](#️-installation-command-line)
- [⚙️ Configuration (`config.toml`)](#️-configuration-configtoml)
- [🔍 Argument Precedence](#-argument-precedence)
- [🧪 Build & Verification](#-build--verification)
- [🧯 Troubleshooting](#-troubleshooting)
- [📦 Assets & Paths](#-assets--paths)

---

## ✨ Overview

`DirAgent` adds file-manager context-menu entries:

- `Open in Codex (DirAgent)`
- `Open in Claude (DirAgent)`

Behavior:

- **Folder selected** → launch inside that folder
- **File selected** → launch inside parent folder

---

## 🎯 Key Features

- 🖱️ Right-click launch for Codex / Claude
- 🧭 Automatic directory resolution (file → parent folder)
- 🪟 Windows menu icons (`.ico`, white background)
- 🔁 Terminal strategy control (`tab_preferred` / `new_window`)
- 🧩 Configurable terminal preference, CLI path, and default args

---

## ⚡ Quick Start (Windows Recommended)

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

---

## 🛠️ Installation (Command Line)

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

---

## ⚙️ Configuration (`config.toml`)

Config file path:

- Windows: `%AppData%\dir-agent\config.toml`
- macOS/Linux: `~/.config/dir-agent/config.toml`

Default config:

```toml
[terminals]
preferred = ""

[tools.codex]
command = "codex"
default_args = []

[tools.claude]
command = "claude"
default_args = []

[behavior]
resolve_file_to_parent = true
open_mode = "tab_preferred"
```

### 📋 Full Parameter Reference

| Key | Type | Default | What it does | When to change |
|---|---|---|---|---|
| `terminals.preferred` | `string` | `""` | Preferred terminal; empty means fallback chain | Multiple terminals installed; need deterministic selection |
| `tools.codex.command` | `string` | `"codex"` | Codex command name or absolute path | `codex` missing in PATH / custom command path |
| `tools.codex.default_args` | `string[]` | `[]` | Default args for every Codex launch | Fixed model / approval / profile defaults |
| `tools.claude.command` | `string` | `"claude"` | Claude command name or absolute path | `claude` missing in PATH / custom command path |
| `tools.claude.default_args` | `string[]` | `[]` | Default args for every Claude launch | Team defaults or personal preferences |
| `behavior.resolve_file_to_parent` | `bool` | `true` | Convert selected file to parent folder | Set `false` for strict directory-only flow |
| `behavior.open_mode` | `string` | `"tab_preferred"` | Controls tab/window behavior | See mode details below |

### 🧠 `open_mode` Details

- `tab_preferred` (default)  
  Reuse current terminal window with a new tab when possible; fallback to new window otherwise.

- `new_window`  
  Always open a new window.

- Any other value  
  Treated as invalid and falls back to `tab_preferred`.

### 🧭 Common `terminals.preferred` values

- Windows: `windows-terminal` / `wezterm` / `powershell`
- macOS: `terminal.app` / `iterm2`
- Linux: `x-terminal-emulator` / `gnome-terminal` / `konsole` / `xterm`

---

## 🔍 Argument Precedence

Merge order (low → high):

1. Built-in defaults  
2. `default_args` from `config.toml`  
3. Passthrough args after `--`

---

## 🧪 Build & Verification

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
   - folder → `Open in Codex (DirAgent)`
   - file → `Open in Claude (DirAgent)` (parent directory)
   - Chinese/space paths
   - icon visibility
4. Double-click `scripts/diragent-3-uninstall-right-click.bat` to verify rollback

---

## 🧯 Troubleshooting

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

---

## 📦 Assets & Paths

Icons are embedded via `go:embed` and extracted during install:

- Windows: `.ico`
- macOS/Linux: `.png`

Asset paths:

- Windows: `%AppData%\dir-agent\assets`
- macOS/Linux: `~/.local/share/dir-agent/assets`
