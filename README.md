# 馃殌 DirAgent

> 鍦ㄦ枃浠剁鐞嗗櫒閲屼竴閿惎鍔?`Codex / Claude`锛岃嚜鍔ㄨ繘鍏ュ綋鍓嶉€変腑鐩綍銆?
馃寪 **Language**: [涓枃](README.md) | [English](README.en.md)


## 馃搶 鐩綍

- [鉁?椤圭洰绠€浠媇(#-椤圭洰绠€浠?
- [馃幆 鏍稿績鍔熻兘](#-鏍稿績鍔熻兘)
- [鈿?蹇€熷紑濮嬶紙Windows 鎺ㄨ崘锛塢(#-蹇€熷紑濮媤indows-鎺ㄨ崘)
- [馃洜锔?瀹夎鏂瑰紡锛堝懡浠よ锛塢(#锔?瀹夎鏂瑰紡鍛戒护琛?
- [鈿欙笍 閰嶇疆璇存槑锛坈onfig.toml锛塢(#锔?閰嶇疆璇存槑configtoml)
- [馃攳 鍙傛暟浼樺厛绾(#-鍙傛暟浼樺厛绾?
- [馃И 寮€鍙戜笌楠岃瘉](#-寮€鍙戜笌楠岃瘉)
- [馃Н 甯歌闂锛圱roubleshooting锛塢(#-甯歌闂troubleshooting)
- [馃摝 璧勬簮涓庤矾寰刔(#-璧勬簮涓庤矾寰?


## 鉁?椤圭洰绠€浠?
`DirAgent` 浼氬湪鏂囦欢绠＄悊鍣ㄤ腑娣诲姞鍙抽敭鑿滃崟锛?
- `Open in Codex (DirAgent)`
- `Open in Claude (DirAgent)`

褰撲綘閫変腑锛?
- **鐩綍 / 鐩綍绌虹櫧澶?*锛氭樉绀哄彸閿彍鍗曪紝骞惰繘鍏ヨ鐩綍鍚姩 CLI
- **鏂囦欢**锛氶粯璁や笉鏄剧ず DirAgent 鑿滃崟锛堥伩鍏嶈涔夋涔夛級


## 馃幆 鏍稿績鍔熻兘

- 馃柋锔?鏂囦欢绠＄悊鍣ㄥ彸閿竴閿惎鍔?Codex / Claude
- 馃幆 鐩綍鑼冨洿鍙抽敭鑿滃崟锛堥伩鍏嶆枃浠舵搷浣滄涔夛級
- 馃獰 Windows 鑿滃崟鍥炬爣锛堢櫧搴?`.ico`锛?- 馃攣 缁堢绛栫暐鍙厤缃紙`tab_preferred` / `new_window`锛?- 馃З 鍙厤缃粓绔紭鍏堢骇銆丆LI 鍛戒护璺緞銆侀粯璁ゅ弬鏁?

## 鈿?蹇€熷紑濮嬶紙Windows 鎺ㄨ崘锛?
鐩存帴鍙屽嚮浠ヤ笅鑴氭湰锛堟棤闇€鎵嬪姩浼犲弬锛夛細

1. `scripts/diragent-1-build-and-verify.bat`  
   - 鎵ц `go test ./...`  
   - 鏋勫缓 `diragent.exe`

2. `scripts/diragent-2-install-right-click.bat`  
   - 鑷姩妫€娴嬪苟鏋勫缓 `diragent.exe`  
   - 瀹夎 Explorer 鍙抽敭鑿滃崟涓庡浘鏍?
3. `scripts/diragent-3-uninstall-right-click.bat`  
   - 鍗歌浇鍙抽敭鑿滃崟  
   - 娓呯悊宸查噴鏀捐祫婧愬拰閰嶇疆


## 馃洜锔?瀹夎鏂瑰紡锛堝懡浠よ锛?
### Windows

鍓嶆彁锛氶」鐩牴鐩綍瀛樺湪 `diragent.exe`锛屾垨绯荤粺 `PATH` 鍙壘鍒?`diragent`銆?
```powershell
# 瀹夎
.\scripts\install.ps1

# 鍗歌浇
.\scripts\uninstall.ps1

# 鍗歌浇骞舵竻鐞?assets + config
.\scripts\uninstall.ps1 -RemoveAssets -RemoveConfig
```

### macOS / Linux

```bash
chmod +x ./scripts/install.sh ./scripts/uninstall.sh
./scripts/install.sh ./diragent
./scripts/uninstall.sh ./diragent
```

> macOS 浼氬垱寤猴細
> - `~/Applications/DirAgent/Open in Codex (DirAgent).app`
> - `~/Applications/DirAgent/Open in Claude (DirAgent).app`


## 鈿欙笍 閰嶇疆璇存槑锛坈onfig.toml锛?
閰嶇疆鏂囦欢璺緞锛?
- Windows锛歚%AppData%\dir-agent\config.toml`
- macOS/Linux锛歚~/.config/dir-agent/config.toml`

榛樿閰嶇疆锛?
```toml
[terminals]
preferred = ""

[terminals.windows_terminal]
profile = ""
shell = "powershell"

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

### 馃搵 鍏ㄥ弬鏁版竻鍗曪紙鍚娇鐢ㄥ満鏅級

| 鍙傛暟 | 绫诲瀷 | 榛樿鍊?| 浣滅敤 | 浠€涔堟椂鍊欐敼 |
|---|---|---|---|---|
| `terminals.preferred` | `string` | `""` | 鎸囧畾棣栭€夌粓绔紱绌哄€兼椂鎸夊唴缃洖閫€閾捐嚜鍔ㄩ€夋嫨 | 鏈哄櫒鏈夊涓粓绔紝鎯冲浐瀹氬叾涓竴涓?|
| `terminals.windows_terminal.profile` | `string` | `""` | Windows Terminal 閰嶇疆鏂囦欢鍚嶏紙渚嬪锛歚Cmder`銆乣PowerShell`銆乣Ubuntu`锛?| 浣跨敤 `windows-terminal` 鏃跺笇鏈涘浐瀹氭煇涓?profile |
| `terminals.windows_terminal.shell` | `string` | `"powershell"` | 鍦?Windows Terminal 涓墽琛?`codex`/`claude` 鐨勫３绫诲瀷锛坄powershell` 鎴?`cmd`锛?| 浣跨敤 Cmd/Cmder 宸ヤ綔娴佹椂鍙涓?`cmd` |
| `tools.codex.command` | `string` | `"codex"` | Codex 鍛戒护鍚嶆垨缁濆璺緞 | `codex` 涓嶅湪 PATH / 鍛戒护鍚嶄笉鍚?|
| `tools.codex.default_args` | `string[]` | `["--dangerously-bypass-approvals-and-sandbox"]` | 姣忔 `Open in Codex` 鑷姩闄勫甫鍙傛暟 | 浠呭湪涓嶅笇鏈涢粯璁ゆ渶楂樻潈闄愭椂淇敼 |
| `tools.claude.command` | `string` | `"claude"` | Claude 鍛戒护鍚嶆垨缁濆璺緞 | `claude` 涓嶅湪 PATH / 鍛戒护鍚嶄笉鍚?|
| `tools.claude.default_args` | `string[]` | `["--dangerously-skip-permissions"]` | 姣忔 `Open in Claude` 鑷姩闄勫甫鍙傛暟 | 浠呭湪涓嶅笇鏈涢粯璁ゆ渶楂樻潈闄愭椂淇敼 |
| `behavior.resolve_file_to_parent` | `bool` | `true` | 閫氳繃 CLI 浼犲叆鏂囦欢璺緞鏃舵槸鍚﹁浆鐖剁洰褰?| 涓€鑸繚鎸?`true` |
| `behavior.open_mode` | `string` | `"tab_preferred"` | 鎺у埗 tab/绐楀彛绛栫暐 | 瑙佷笅鏂硅缁嗚鏄?|

### 馃 `open_mode` 璇﹁В

- `tab_preferred`锛堥粯璁わ級  
  浼樺厛鍦ㄥ綋鍓嶇粓绔獥鍙ｅ紑鏂?Tab锛涜嫢涓嶅彲鐢ㄥ垯鏂板缓绐楀彛銆?
- `new_window`  
  姣忔閮芥柊寤虹獥鍙ｃ€?
- 鍏朵粬鍊? 
  瑙嗕负鏃犳晥鍊硷紝鍥為€€鍒?`tab_preferred`銆?
### 馃Л `terminals.preferred` 甯歌鍙栧€?
- Windows锛歚windows-terminal` / `wezterm` / `powershell`
- macOS锛歚terminal.app` / `iterm2`
- Linux锛歚x-terminal-emulator` / `gnome-terminal` / `konsole` / `xterm`

### Windows Terminal profile/shell 绀轰緥

```toml
[terminals]
preferred = "windows-terminal"

[terminals.windows_terminal]
profile = "Cmder"
shell = "cmd"
```


## 馃攳 鍙傛暟浼樺厛绾?
鍙傛暟鍚堝苟椤哄簭锛堜綆 鈫?楂橈級锛?
1. 绋嬪簭榛樿鍊? 
2. `config.toml` 鐨?`default_args`  
3. CLI `--` 鍚庨€忎紶鍙傛暟


## 馃И 寮€鍙戜笌楠岃瘉

### 鏋勫缓

```powershell
# Windows
go build -o diragent.exe ./cmd/diragent
```

```bash
# macOS/Linux
go build -o diragent ./cmd/diragent
```

### 娴嬭瘯

```bash
go test ./...
```

### Windows 楠屾敹寤鸿

1. 鍙屽嚮 `scripts/diragent-1-build-and-verify.bat`
2. 鍙屽嚮 `scripts/diragent-2-install-right-click.bat`
3. 鎵嬪伐楠岃瘉锛?   - 鐩綍鍙抽敭 `Open in Codex (DirAgent)`
   - 鏂囦欢鍙抽敭锛氫笉鏄剧ず DirAgent 鑿滃崟
   - 涓枃/绌烘牸璺緞姝ｅ父
   - 鍥炬爣鏄剧ず姝ｅ父
4. 鍙屽嚮 `scripts/diragent-3-uninstall-right-click.bat` 楠岃瘉鍙洖婊?

## 馃Н 甯歌闂锛圱roubleshooting锛?
### 1) 鎶ラ敊 `2147942402 (0x80070002)`锛孋odex 鍚姩澶辫触

鍚箟锛氱郴缁熸壘涓嶅埌鍙墽琛屽懡浠ゃ€? 
鎺掓煡椤哄簭锛?
1. 鍦?PowerShell 鎵ц `Get-Command codex`
2. 鑻ユ湭鎵惧埌锛屽湪 `config.toml` 璁剧疆 `tools.codex.command` 涓烘纭懡浠ゆ垨缁濆璺緞
3. 閲嶆柊鎵ц `scripts/diragent-2-install-right-click.bat`

### 2) 鑿滃崟宸插畨瑁呬絾鐪嬩笉鍒?
- 鍦ㄦ枃浠跺す绌虹櫧澶勬寜 `F5` 鍒锋柊
- 鎴栭噸鍚?Explorer
- 纭瀹夎鍙戠敓鍦ㄥ綋鍓嶇敤鎴凤紙`HKCU`锛?
### 3) 娌℃寜鈥滃悓缁堢鏂?Tab鈥濊涓哄惎鍔?
- 纭 `behavior.open_mode = "tab_preferred"`
- 鑻ラ閫夌粓绔笉鏀寔 Tab 澶嶇敤锛屼細鍥為€€鍒版柊绐楀彛


## 馃摝 璧勬簮涓庤矾寰?
鍥炬爣閫氳繃 `go:embed` 鍐呭祵锛屽畨瑁呮椂閲婃斁鍒版湰鍦帮細

- Windows锛歚.ico`
- macOS/Linux锛歚.png`

璧勬簮鐩綍锛?
- Windows锛歚%AppData%\dir-agent\assets`
- macOS/Linux锛歚~/.local/share/dir-agent/assets`
