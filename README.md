# 🚀 DirAgent

在文件管理器中右键目录，可在目标文件夹一键启动 `Codex` 或 `Claude Code`！

🌐 语言: [中文](README.md) | [English](README.en.md)

![Demo](docs/demo.png)

## ✨ 解决什么问题

对我来说，当前每次使用 `codex`或者 `cc`，都需要`打开终端 -> 复制目标目录 -> cd 到目录 -> 启动 codex/claudecode `  ，很是麻烦！

大部分人都是通过系统自带的文件管理器 GUI 方式来访问文件夹的，所以为什么不给这个操作装一个一键启动 agent呢？浏览到哪文件夹，然后一键启动 codex 或 claude code，同时可以传递配置参数！

## 📦 快速开始

从 **Release -> Assets** 下载对应平台的压缩包

按系统选择一个 zip：

- Windows x64: `diragent_<tag>_windows_amd64.zip`
- Windows ARM64: `diragent_<tag>_windows_arm64.zip`
- macOS Intel: `diragent_<tag>_darwin_amd64.zip`
- macOS Apple Silicon: `diragent_<tag>_darwin_arm64.zip`
- Linux x64: `diragent_<tag>_linux_amd64.zip`
- Linux ARM64: `diragent_<tag>_linux_arm64.zip`

### ⚡ 一键安装

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

### 🧹 一键卸载

- Windows: 双击 `uninstall.bat`
- macOS / Linux:
```bash
chmod +x ./uninstall.sh
./uninstall.sh
```

## 🧭 配置与数据位置

当安装过后，会在当前工程位置生成一个 `config/toml`文件，默认它会用系统当前可用的终端，你也可以客制化它

- 默认位置：
  - 配置文件：`<安装目录>/config.toml`
  - 资源文件：`<安装目录>/data/assets`
- 可选覆盖：
  - 设置环境变量 `DIRAGENT_HOME` 后，路径变为
    - `DIRAGENT_HOME/config.toml`
    - `DIRAGENT_HOME/data/assets`

### ⚙️ `config.toml` 配置示例

```toml
[terminals]
preferred = "" # 为空表示自动探测，常用值：windows-terminal / wezterm / powershell

[terminals.windows_terminal]
profile = ""   # 可选：指定 Windows Terminal profile 名称
shell = "powershell" # 可选：powershell / cmd / cmder
cmder_init = "" # shell=cmder 时可指定 init.bat 路径

[tools.codex]
command = "codex" # 可改为绝对路径
default_args = ["--dangerously-bypass-approvals-and-sandbox"]

[tools.claude]
command = "claude" # 可改为绝对路径
default_args = ["--dangerously-skip-permissions"]

[behavior]
resolve_file_to_parent = true # 对文件右键时，自动使用其父目录
open_mode = "tab_preferred"   # tab_preferred / new_window
```

### 🔎 关键配置项说明

- `terminals.preferred`：
  - 终端优先级。留空则自动选择可用终端。
- `terminals.windows_terminal.profile`：
  - 可选，指定 WT profile 名称（如 `Command Prompt`、`Cmder`）。
- `terminals.windows_terminal.shell`：
  - `powershell`、`cmd`、`cmder` 三选一。
- `tools.codex.command` / `tools.claude.command`：
  - 命令名或绝对路径；命令找不到时优先检查这里。
- `behavior.open_mode`：
  - `tab_preferred`：优先新 tab；
  - `new_window`：总是新窗口。

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
