# 🚀 DirAgent

Right-click a folder in your file manager and launch `Codex` or `Claude Code` in that target folder with one click.

🌐 Language: [English](README.en.md) | [中文](README.md)

![Demo](docs/demo.png)

## ✨ What Problem It Solves

For me, every time I use `codex` or `cc`, I need to do:
`open terminal -> copy target path -> cd to target folder -> run codex/claudecode`.
It is repetitive and annoying.

Most people browse folders from the system file-manager GUI anyway, so it makes sense to add a one-click agent launcher there:
browse to any folder, then launch `codex` or `claude code` directly, while still being able to pass config options.

## 📦 Quick Start

Download the zip package for your platform from **Release -> Assets**.

Choose one zip by OS/arch:

- Windows x64: `diragent_<tag>_windows_amd64.zip`
- Windows ARM64: `diragent_<tag>_windows_arm64.zip`
- macOS Intel: `diragent_<tag>_darwin_amd64.zip`
- macOS Apple Silicon: `diragent_<tag>_darwin_arm64.zip`
- Linux x64: `diragent_<tag>_linux_amd64.zip`
- Linux ARM64: `diragent_<tag>_linux_arm64.zip`

### ⚡ One-Click Install

1. Extract the zip to your preferred DirAgent install folder.
2. Run the install entrypoint:
- Windows: double-click `install.bat`
- macOS / Linux:
```bash
chmod +x ./install.sh
./install.sh
```
3. Right-click any directory and launch through the DirAgent menu.

Notes:
- `install` first cleans previous integration (keeps existing config), then installs again.
- Release bundles expose only two user entrypoints: `install` and `uninstall`.

### 🧹 One-Click Uninstall

- Windows: double-click `uninstall.bat`
- macOS / Linux:
```bash
chmod +x ./uninstall.sh
./uninstall.sh
```

## 🧭 Config And Data Locations

After installation, a `config.toml` file is created in your install location. By default, DirAgent uses an available terminal on your system, and you can customize it.

- Default locations:
  - Config: `<install-folder>/config.toml`
  - Assets: `<install-folder>/data/assets`
- Optional override:
  - After setting `DIRAGENT_HOME`, paths become:
    - `DIRAGENT_HOME/config.toml`
    - `DIRAGENT_HOME/data/assets`

### ⚙️ `config.toml` Example

```toml
[terminals]
preferred = "" # empty means auto-detect; common values: windows-terminal / wezterm / powershell

[terminals.windows_terminal]
profile = ""   # optional: Windows Terminal profile name
shell = "powershell" # optional: powershell / cmd / cmder
cmder_init = "" # optional init.bat path when shell=cmder

[tools.codex]
command = "codex" # can be absolute path
default_args = ["--dangerously-bypass-approvals-and-sandbox"]

[tools.claude]
command = "claude" # can be absolute path
default_args = ["--dangerously-skip-permissions"]

[behavior]
resolve_file_to_parent = true # when right-clicking a file, use its parent folder
open_mode = "tab_preferred"   # tab_preferred / new_window
```

### 🔎 Key Configuration Fields

- `terminals.preferred`:
  - Terminal priority. Empty means auto-select from available terminals.
- `terminals.windows_terminal.profile`:
  - Optional WT profile name (for example `Command Prompt`, `Cmder`).
- `terminals.windows_terminal.shell`:
  - One of `powershell`, `cmd`, `cmder`.
- `tools.codex.command` / `tools.claude.command`:
  - Command name or absolute path. Check these first if command is not found.
- `behavior.open_mode`:
  - `tab_preferred`: prefer opening a new tab.
  - `new_window`: always open a new window.

## 🛠️ Quick Troubleshooting

- `0x80070002` or command not found:
  set `tools.codex.command` or `tools.claude.command` in `config.toml` to an absolute executable path.
- Menu not visible:
  refresh file manager (`F5`) or restart Explorer/Finder.
- WezTerm does not open tab as expected:
  set `terminals.preferred = "wezterm"` and `behavior.open_mode = "tab_preferred"`.

## 👩‍💻 Development

Run tests:

```bash
go test ./...
```

Build binaries:

```powershell
go build -o diragent.exe ./cmd/diragent
go build -o diragentw.exe ./cmd/diragentw
```

Local Windows scripts from repository:

- `scripts/install.bat`:
  uninstall previous integration if present, build latest `diragent.exe` and `diragentw.exe`, then install the context menu.
- `scripts/uninstall.bat`:
  uninstall only.

License: MIT (see `LICENSE`).
