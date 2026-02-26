# 🚀 DirAgent

在文件管理器中右键目录，一键启动 `Codex` 或 `Claude Code`，并自动进入目标目录。

🌐 语言: [中文](README.md) | [English](README.en.md)

![Demo](docs/demo.png)

## ✨ 解决什么问题

DirAgent 把这套重复动作收敛成一次右键：

`打开终端 -> cd 到目录 -> 启动 codex/claude`

安装后可用菜单：

- `Open in Codex (DirAgent)`
- `Open in Claude Code (DirAgent)`

## 📦 下载哪个文件

只从 **Release -> Assets** 下载，不要使用 `Source code (zip/tar.gz)`。

按系统选择一个 zip：

- Windows x64: `diragent_<tag>_windows_amd64.zip`
- Windows ARM64: `diragent_<tag>_windows_arm64.zip`
- macOS Intel: `diragent_<tag>_darwin_amd64.zip`
- macOS Apple Silicon: `diragent_<tag>_darwin_arm64.zip`
- Linux x64: `diragent_<tag>_linux_amd64.zip`
- Linux ARM64: `diragent_<tag>_linux_arm64.zip`

每个 zip 都是完整可用包，不需要 clone 仓库。

## ⚡ 一键安装

1. 解压 zip 到你希望安装 DirAgent 的目录。
2. 运行安装入口：
- Windows: 双击 `install.bat`
- macOS / Linux:
```bash
chmod +x ./install.sh
./install.sh
```
3. 在任意目录右键，通过 DirAgent 菜单启动。

说明：
- `install` 会先清理旧集成（保留已有配置），再重新安装。
- Release 包面对用户的入口只有两个：`install` 和 `uninstall`。

## 🧹 一键卸载

- Windows: 双击 `uninstall.bat`
- macOS / Linux:
```bash
chmod +x ./uninstall.sh
./uninstall.sh
```

## 🧭 配置与数据位置

- 配置文件：`<安装目录>/config.toml`
- 资源文件：`<安装目录>/data/assets`

## 🛠️ 快速排障

- 报错 `0x80070002` / command not found：
  在 `config.toml` 把 `tools.codex.command` 或 `tools.claude.command` 改成绝对路径。
- 右键菜单没出现：
  刷新文件管理器（`F5`）或重启 Explorer/Finder。
- WezTerm 没按预期开 tab：
  设置 `terminals.preferred = "wezterm"` 与 `behavior.open_mode = "tab_preferred"`。

## 👩‍💻 开发

运行测试：

```bash
go test ./...
```

构建二进制：

```powershell
go build -o diragent.exe ./cmd/diragent
go build -o diragentw.exe ./cmd/diragentw
```

仓库内 Windows 本地脚本：

- `scripts/install.bat`：
  先卸载旧集成（若有），再构建最新 `diragent.exe` 和 `diragentw.exe`，最后安装右键菜单。
- `scripts/uninstall.bat`：
  仅执行卸载。

许可证：MIT（见 `LICENSE`）。
