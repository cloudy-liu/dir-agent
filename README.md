# DirAgent

在文件管理器中右键目录，一键启动 `Codex` 或 `Claude Code`，并自动进入目标目录。

语言: [中文](README.md) | [English](README.en.md)

## 价值主张

DirAgent 解决重复操作：

`打开终端 -> cd 到目录 -> 启动 codex/claude`

安装后，你只需要右键目录：

- `Open in Codex (DirAgent)`
- `Open in Claude Code (DirAgent)`

## Release 下载怎么选

请在 **Release -> Assets** 下载，不要使用 `Source code (zip/tar.gz)`。

按系统和架构选择：

- Windows x64: `diragent_<tag>_windows_amd64.zip`
- Windows ARM64: `diragent_<tag>_windows_arm64.zip`
- macOS Intel: `diragent_<tag>_darwin_amd64.zip`
- macOS Apple Silicon: `diragent_<tag>_darwin_arm64.zip`
- Linux x64: `diragent_<tag>_linux_amd64.zip`
- Linux ARM64: `diragent_<tag>_linux_arm64.zip`

每个压缩包已经包含二进制和安装/卸载脚本，**无需 clone 仓库**。

## 3 步安装

1. 解压到你希望安装 DirAgent 的目录。
2. 在该目录执行安装脚本：

Windows (PowerShell):

```powershell
.\scripts\install.ps1 -BinaryPath .\diragent.exe
```

macOS / Linux:

```bash
chmod +x ./diragent ./scripts/install.sh ./scripts/uninstall.sh
./scripts/install.sh ./diragent
```

3. 右键任意目录，选择 DirAgent 菜单项启动。

## 配置与数据位置

- 配置文件: `<安装目录>/config.toml`
- 资源文件: `<安装目录>/data/assets`

## 快速排错

- 报错 `0x80070002` / command not found：
  在 `config.toml` 中将 `tools.codex.command` 或 `tools.claude.command` 改为绝对路径。
- 装完菜单看不到：
  刷新文件管理器 (`F5`) 或重启 Explorer/Finder。
- WezTerm 没有按预期开 tab：
  设置 `terminals.preferred = "wezterm"` 和 `behavior.open_mode = "tab_preferred"`。

## 开发

```bash
go test ./...
```

```powershell
go build -o diragent.exe ./cmd/diragent
```

许可证：MIT（见 `LICENSE`）。
