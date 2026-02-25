# 🚀 DirAgent

在文件管理器中右键目录，一键启动 `Codex` 或 `Claude Code`，并自动进入目标目录。

🌐 语言: [中文](README.md) | [English](README.en.md)

![Demo](docs/demo.png)

## ✨ 这是什么

DirAgent 把这套重复动作收敛成一次右键：

`打开终端 -> cd 到目录 -> 启动 codex/claude`

安装后可用菜单：

- `Open in Codex (DirAgent)`
- `Open in Claude Code (DirAgent)`

## 📦 下载哪个文件

只从 **Release -> Assets** 下载，不要用 `Source code (zip/tar.gz)`。

按系统选择一个 zip：

- Windows x64: `diragent_<tag>_windows_amd64.zip`
- Windows ARM64: `diragent_<tag>_windows_arm64.zip`
- macOS Intel: `diragent_<tag>_darwin_amd64.zip`
- macOS Apple Silicon: `diragent_<tag>_darwin_arm64.zip`
- Linux x64: `diragent_<tag>_linux_amd64.zip`
- Linux ARM64: `diragent_<tag>_linux_arm64.zip`

每个 zip 已包含完整可用内容，不需要 clone 仓库。

## ⚡ 一键安装

1. 解压 zip 到你想安装 DirAgent 的目录。
2. 在解压目录运行安装入口（只用这一个）：
   - Windows: 双击 `install.bat`
   - macOS / Linux:
     ```bash
     chmod +x ./install.sh
     ./install.sh
     ```
3. 右键任意目录，选择 DirAgent 菜单项启动。

说明：
- `install` 会先清理旧安装（保留现有配置），再重新安装。
- 用户可见入口只有两个：`install` 和 `uninstall`。

## 🧹 一键卸载

- Windows: 双击 `uninstall.bat`
- macOS / Linux:
  ```bash
  chmod +x ./uninstall.sh
  ./uninstall.sh
  ```

## 🧭 配置和数据在哪

- 配置文件: `<安装目录>/config.toml`
- 资源文件: `<安装目录>/data/assets`

## 🛠️ 快速排错

- 报错 `0x80070002` / command not found：
  在 `config.toml` 把 `tools.codex.command` 或 `tools.claude.command` 改成绝对路径。
- 右键菜单没出现：
  刷新文件管理器 (`F5`) 或重启 Explorer/Finder。
- WezTerm 没有按预期开 tab：
  设置 `terminals.preferred = "wezterm"` 与 `behavior.open_mode = "tab_preferred"`。

## 👩‍💻 开发

```bash
go test ./...
```

```powershell
go build -o diragent.exe ./cmd/diragent
```

许可证：MIT（见 `LICENSE`）。
