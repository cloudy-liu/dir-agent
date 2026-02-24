# 🚀 DirAgent

> 在文件管理器里一键启动 `Codex / Claude`，自动进入当前选中目录。

🌐 **Language**: [中文](README.md) | [English](README.en.md)


## 📌 目录

- [✨ 项目简介](#-项目简介)
- [🎯 核心功能](#-核心功能)
- [⚡ 快速开始（Windows 推荐）](#-快速开始windows-推荐)
- [🛠️ 安装方式（命令行）](#️-安装方式命令行)
- [⚙️ 配置说明（config.toml）](#️-配置说明configtoml)
- [🔍 参数优先级](#-参数优先级)
- [🧪 开发与验证](#-开发与验证)
- [🧯 常见问题（Troubleshooting）](#-常见问题troubleshooting)
- [📦 资源与路径](#-资源与路径)


## ✨ 项目简介

`DirAgent` 会在文件管理器中添加右键菜单：

- `Open in Codex (DirAgent)`
- `Open in Claude (DirAgent)`

当你选中：

- **目录 / 目录空白处**：显示右键菜单，并进入该目录启动 CLI
- **文件**：默认不显示 DirAgent 菜单（避免语义歧义）


## 🎯 核心功能

- 🖱️ 文件管理器右键一键启动 Codex / Claude
- 🎯 目录范围右键菜单（避免文件操作歧义）
- 🪟 Windows 菜单图标（白底 `.ico`）
- 🔁 终端策略可配置（`tab_preferred` / `new_window`）
- 🧩 可配置终端优先级、CLI 命令路径、默认参数


## ⚡ 快速开始（Windows 推荐）

直接双击以下脚本（无需手动传参）：

1. `scripts/diragent-1-build-and-verify.bat`  
   - 执行 `go test ./...`  
   - 构建 `diragent.exe`

2. `scripts/diragent-2-install-right-click.bat`  
   - 自动检测并构建 `diragent.exe`  
   - 安装 Explorer 右键菜单与图标

3. `scripts/diragent-3-uninstall-right-click.bat`  
   - 卸载右键菜单  
   - 清理已释放资源和配置


## 🛠️ 安装方式（命令行）

### Windows

前提：项目根目录存在 `diragent.exe`，或系统 `PATH` 可找到 `diragent`。

```powershell
# 安装
.\scripts\install.ps1

# 卸载
.\scripts\uninstall.ps1

# 卸载并清理 assets + config
.\scripts\uninstall.ps1 -RemoveAssets -RemoveConfig
```

### macOS / Linux

```bash
chmod +x ./scripts/install.sh ./scripts/uninstall.sh
./scripts/install.sh ./diragent
./scripts/uninstall.sh ./diragent
```

> macOS 会创建：
> - `~/Applications/DirAgent/Open in Codex (DirAgent).app`
> - `~/Applications/DirAgent/Open in Claude (DirAgent).app`


## ⚙️ 配置说明（config.toml）

配置文件路径：

- Windows：`%AppData%\dir-agent\config.toml`
- macOS/Linux：`~/.config/dir-agent/config.toml`

默认配置：

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

### 📋 全参数清单（含使用场景）

| 参数 | 类型 | 默认值 | 作用 | 什么时候改 |
|---|---|---|---|---|
| `terminals.preferred` | `string` | `""` | 指定首选终端；空值时按内置回退链自动选择 | 机器有多个终端，想固定其中一个 |
| `terminals.windows_terminal.profile` | `string` | `""` | Windows Terminal 配置文件名（例如：`Cmder`、`PowerShell`、`Command Prompt`） | 使用 `windows-terminal` 时希望固定某个 profile |
| `terminals.windows_terminal.shell` | `string` | `"powershell"` | 在 Windows Terminal 中执行 `codex`/`claude` 的壳类型（`powershell`、`cmd` 或 `cmder`） | 使用 Cmder 初始化流程时设为 `cmder` |
| `terminals.windows_terminal.cmder_init` | `string` | `""` | `shell = "cmder"` 时可选的 `init.bat` 路径 | 无法通过 `CMDER_ROOT` 自动发现时显式配置 |
| `tools.codex.command` | `string` | `"codex"` | Codex 命令名或绝对路径 | `codex` 不在 PATH / 命令名不同 |
| `tools.codex.default_args` | `string[]` | `["--dangerously-bypass-approvals-and-sandbox"]` | 每次 `Open in Codex` 自动附带参数 | 仅在不希望默认最高权限时修改 |
| `tools.claude.command` | `string` | `"claude"` | Claude 命令名或绝对路径 | `claude` 不在 PATH / 命令名不同 |
| `tools.claude.default_args` | `string[]` | `["--dangerously-skip-permissions"]` | 每次 `Open in Claude` 自动附带参数 | 仅在不希望默认最高权限时修改 |
| `behavior.resolve_file_to_parent` | `bool` | `true` | 通过 CLI 传入文件路径时是否转父目录 | 一般保持 `true` |
| `behavior.open_mode` | `string` | `"tab_preferred"` | 控制 tab/窗口策略 | 见下方详细说明 |

### 🧠 `open_mode` 详解

- `tab_preferred`（默认）  
  优先在当前终端窗口开新 Tab；若不可用则新建窗口。

- `new_window`  
  每次都新建窗口。

- 其他值  
  视为无效值，回退到 `tab_preferred`。

### 🧭 `terminals.preferred` 常见取值

- Windows：`windows-terminal` / `wezterm` / `powershell`
- macOS：`terminal.app` / `iterm2`
- Linux：`x-terminal-emulator` / `gnome-terminal` / `konsole` / `xterm`

### Windows Terminal profile/shell 示例

```toml
[terminals]
preferred = "windows-terminal"

[terminals.windows_terminal]
profile = "Cmder"
shell = "cmder"
cmder_init = "C:\\path\\to\\cmder\\vendor\\init.bat"
```


## 🔍 参数优先级

参数合并顺序（低 → 高）：

1. 程序默认值  
2. `config.toml` 的 `default_args`  
3. CLI `--` 后透传参数


## 🧪 开发与验证

### 构建

```powershell
# Windows
go build -o diragent.exe ./cmd/diragent
```

```bash
# macOS/Linux
go build -o diragent ./cmd/diragent
```

### 测试

```bash
go test ./...
```

### Windows 验收建议

1. 双击 `scripts/diragent-1-build-and-verify.bat`
2. 双击 `scripts/diragent-2-install-right-click.bat`
3. 手工验证：
   - 目录右键 `Open in Codex (DirAgent)`
   - 文件右键：不显示 DirAgent 菜单
   - 中文/空格路径正常
   - 图标显示正常
4. 双击 `scripts/diragent-3-uninstall-right-click.bat` 验证可回滚


## 🧯 常见问题（Troubleshooting）

### 1) 报错 `2147942402 (0x80070002)`，Codex 启动失败

含义：系统找不到可执行命令。  
排查顺序：

1. 在 PowerShell 执行 `Get-Command codex`
2. 若未找到，在 `config.toml` 设置 `tools.codex.command` 为正确命令或绝对路径
3. 重新执行 `scripts/diragent-2-install-right-click.bat`

### 2) 菜单已安装但看不到

- 在文件夹空白处按 `F5` 刷新
- 或重启 Explorer
- 确认安装发生在当前用户（`HKCU`）

### 3) 没按“同终端新 Tab”行为启动

- 确认 `behavior.open_mode = "tab_preferred"`
- 若首选终端不支持 Tab 复用，会回退到新窗口


## 📦 资源与路径

图标通过 `go:embed` 内嵌，安装时释放到本地：

- Windows：`.ico`
- macOS/Linux：`.png`

资源目录：

- Windows：`%AppData%\dir-agent\assets`
- macOS/Linux：`~/.local/share/dir-agent/assets`
