# 🚀 DirAgent

Launch `Codex` or `Claude Code` from your file manager and start directly in the target directory.

🌐 Language: [English](README.en.md) | [中文](README.md)

![Demo](docs/demo.png)

## ✨ What It Solves

DirAgent removes the repeated flow:

`open terminal -> cd into folder -> run codex/claude`

After install, use the context menu:

- `Open in Codex (DirAgent)`
- `Open in Claude Code (DirAgent)`

## 📦 Which File To Download

Download only from **Release -> Assets**.
Do not use `Source code (zip/tar.gz)`.

Choose one zip by OS/arch:

- Windows x64: `diragent_<tag>_windows_amd64.zip`
- Windows ARM64: `diragent_<tag>_windows_arm64.zip`
- macOS Intel: `diragent_<tag>_darwin_amd64.zip`
- macOS Apple Silicon: `diragent_<tag>_darwin_arm64.zip`
- Linux x64: `diragent_<tag>_linux_amd64.zip`
- Linux ARM64: `diragent_<tag>_linux_arm64.zip`

Each zip is self-contained. No repository clone required.

## ⚡ One-Click Install

1. Extract the zip into your preferred install folder.
2. Run the install entrypoint:
- Windows: double-click `install.bat`
- macOS / Linux:
```bash
chmod +x ./install.sh
./install.sh
```
3. Right-click any directory and launch via DirAgent menu.

Notes:
- `install` first removes previous integration (keeps existing config), then installs again.
- User-facing entrypoints in release bundles are only two files: `install` and `uninstall`.

## 🧹 One-Click Uninstall

- Windows: double-click `uninstall.bat`
- macOS / Linux:
```bash
chmod +x ./uninstall.sh
./uninstall.sh
```

## 🧭 Config And Data Locations

- Config: `<install-folder>/config.toml`
- Assets: `<install-folder>/data/assets`

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
  uninstall previous integration if present, build latest `diragent.exe` and `diragentw.exe`, then install context menu.
- `scripts/uninstall.bat`:
  uninstall only.

License: MIT (see `LICENSE`).
