# 🚀 DirAgent

> Launch `Codex / Claude` directly from your file manager and start in the target directory.

🌐 Language: [English](README.en.md) | [中文](README.md)

## ✨ Overview

`DirAgent` is a directory-context launcher. It turns "open terminal -> `cd` -> run command" into one right-click action.

After install, it adds:

- `Open in Codex (DirAgent)`
- `Open in Claude (DirAgent)`

![Demo Screenshot](docs/demo.png)

## 🎯 Why this project exists

This project was created to solve a few very common real-world problems:

- Many people use Agent CLIs (`codex`, `claude code`) all day, but still need to manually open a terminal and switch directories every time.
- Most users naturally navigate via the OS file manager, so the shortest path is launching Agent CLI directly during file browsing.
- DirAgent exists for exactly this: right-click any directory and start an Agent immediately (currently supports Codex and Claude Code).

## 🧠 How it works

1. Install scripts register directory-level context-menu entries.
2. On right-click launch, the selected directory is passed to `diragent`.
3. `diragent` resolves tool, args, terminal, and window mode from `config.toml`.
4. `codex` or `claude` starts directly in that directory.

Notes:

- File right-click is hidden by design to avoid ambiguity.
- `open_mode` supports tab-preferred or always-new-window behavior.

## ✅ Features

- One-click directory launch for Codex / Claude
- Directory-only menu visibility policy (hidden for file context)
- Windows Terminal profile/shell support (`powershell` / `cmd` / `cmder`)
- Configurable default launch arguments (including full-access defaults)
- Cross-platform install scripts (Windows / macOS / Linux)

## 🛠️ Installation

### 🪟 Windows

Recommended: one-click `bat` scripts:

1. `scripts/diragent-1-build-and-verify.bat`  
   Builds and runs `go test ./...`
2. `scripts/diragent-2-install-right-click.bat`  
   Installs Explorer context-menu entries
3. `scripts/diragent-3-uninstall-right-click.bat`  
   Uninstalls entries (rollback path)

You can also install from PowerShell:

```powershell
# Install
.\scripts\install.ps1

# Uninstall
.\scripts\uninstall.ps1

# Uninstall and clean assets + config
.\scripts\uninstall.ps1 -RemoveAssets -RemoveConfig
```

### 🍎🐧 macOS / Linux

```bash
chmod +x ./scripts/install.sh ./scripts/uninstall.sh
./scripts/install.sh ./diragent
./scripts/uninstall.sh ./diragent
```

## 📦 How to use release assets

Download binaries from the `Assets` section of a release, not from `Source code (zip/tar.gz)`.

Pick by OS/arch:

- Windows x64: `diragent_v0.5_windows_amd64.exe`
- Windows ARM64: `diragent_v0.5_windows_arm64.exe`
- macOS Intel: `diragent_v0.5_darwin_amd64`
- macOS Apple Silicon: `diragent_v0.5_darwin_arm64`
- Linux x64: `diragent_v0.5_linux_amd64`
- Linux ARM64: `diragent_v0.5_linux_arm64`

Optional verification (recommended):

- Download `SHA256SUMS.txt` and verify binary integrity.
- Windows (PowerShell): `Get-FileHash .\diragent_v0.5_windows_amd64.exe -Algorithm SHA256`
- macOS/Linux: `sha256sum ./diragent_v0.5_linux_amd64`

After download, choose one flow:

1. Install context menu directly from downloaded binary (recommended)
   - Windows:
   ```powershell
   .\scripts\install.ps1 -BinaryPath .\diragent_v0.5_windows_amd64.exe
   ```
   - macOS/Linux:
   ```bash
   chmod +x ./diragent_v0.5_linux_amd64
   ./scripts/install.sh ./diragent_v0.5_linux_amd64
   ```
2. Rename to `diragent` / `diragent.exe`, place it in `PATH`, then run install scripts

## ▶️ Usage

1. Right-click a folder or folder background.
2. Choose `Open in Codex (DirAgent)` or `Open in Claude (DirAgent)`.
3. The CLI starts in the selected directory.

Expected behavior:

- No DirAgent item for file right-click.
- If tab reuse is unavailable, fallback is a new window.

## ⚙️ Configuration: what, why, how

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
cmder_init = ""

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

Core concepts:

1. `tools.*`  
   Which executable is launched and which default args are always appended.
2. `terminals.*`  
   Which terminal is used to host the launch.
3. `behavior.*`  
   How path resolution and tab/window behavior should work.

Most frequently changed keys:

- `tools.codex.command` / `tools.claude.command`  
  Use absolute paths when commands are not in `PATH`.
- `tools.codex.default_args` / `tools.claude.default_args`  
  Control default launch args (currently full-access defaults).
- `terminals.preferred`  
  Pin a preferred terminal, such as `windows-terminal`.
- `terminals.windows_terminal.profile`  
  Pin a Windows Terminal profile such as `Cmder`.
- `terminals.windows_terminal.shell`  
  Supported values: `powershell`, `cmd`, `cmder`.
- `behavior.open_mode`  
  `tab_preferred` (prefer new tab) or `new_window` (always new window).

Windows Terminal + Cmder example:

```toml
[terminals]
preferred = "windows-terminal"

[terminals.windows_terminal]
profile = "Cmder"
shell = "cmder"
cmder_init = "C:\\path\\to\\cmder\\vendor\\init.bat"
```

Argument precedence (low -> high):

1. Built-in defaults
2. `default_args` from `config.toml`
3. Passthrough args after `--`

## 🧯 Troubleshooting

### `2147942402 (0x80070002)` launch error

Usually means command not found:

1. Run `Get-Command codex` in PowerShell
2. Set `tools.codex.command` in `config.toml` to the correct command/path
3. Re-run the install script

### Why is there no menu on file right-click

Expected behavior. DirAgent is intentionally directory-only.

### Why did it open a new window instead of a tab

Check `behavior.open_mode = "tab_preferred"`.  
If the terminal cannot reuse tabs, fallback is a new window.

## 🧪 Development

```bash
go test ./...
```

```powershell
go build -o diragent.exe ./cmd/diragent
```

```bash
go build -o diragent ./cmd/diragent
```

## 🤝 Contributing

Issues and PRs are welcome.

## 📄 License

MIT, see `LICENSE`.
