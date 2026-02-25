# 🚀 DirAgent

Launch `Codex` or `Claude Code` from your file manager and start directly in the target directory.

🌐 Language: [English](README.en.md) | [中文](README.md)

![Demo](docs/demo.png)

## ✨ What it solves

DirAgent removes the repeated flow:

`open terminal -> cd into folder -> run codex/claude`

After install, use right-click menu:

- `Open in Codex (DirAgent)`
- `Open in Claude Code (DirAgent)`

## 📦 Which file to download

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

## ⚡ One-click install

1. Extract the zip into your preferred install folder.
2. Run the single install entrypoint:
   - Windows: double-click `install.bat`
   - macOS / Linux:
     ```bash
     chmod +x ./install.sh
     ./install.sh
     ```
3. Right-click any directory and launch via DirAgent menu.

Notes:
- `install` first cleans previous install (keeps existing config), then installs again.
- User-facing entrypoints are only two files: `install` and `uninstall`.

## 🧹 One-click uninstall

- Windows: double-click `uninstall.bat`
- macOS / Linux:
  ```bash
  chmod +x ./uninstall.sh
  ./uninstall.sh
  ```

## 🧭 Config and data locations

- Config: `<install-folder>/config.toml`
- Assets: `<install-folder>/data/assets`

## 🛠️ Quick troubleshooting

- `0x80070002` or command not found:
  set `tools.codex.command` or `tools.claude.command` in `config.toml` to an absolute executable path.
- Menu not visible:
  refresh file manager (`F5`) or restart Explorer/Finder.
- WezTerm does not open tab as expected:
  set `terminals.preferred = "wezterm"` and `behavior.open_mode = "tab_preferred"`.

## 👩‍💻 Development

```bash
go test ./...
```

```bash
go build -o diragent ./cmd/diragent
```

MIT License. See `LICENSE`.
