# GitHub Release 多平台产物说明

仓库已新增工作流：`.github/workflows/release.yml`。

## 触发方式

- 发布 `GitHub Release`（`published`）时自动触发。
- 也支持手动触发 `workflow_dispatch`，并输入目标 `tag`。

## 构建矩阵

当前会构建以下平台二进制：

- `linux/amd64`
- `linux/arm64`
- `windows/amd64`
- `windows/arm64`
- `darwin/amd64`
- `darwin/arm64`

## Release 资产命名

发布到 Release 的文件名统一为：

`diragent_<tag>_<goos>_<goarch>[.exe]`

示例：

- `diragent_v1.0.0_linux_amd64`
- `diragent_v1.0.0_darwin_arm64`
- `diragent_v1.0.0_windows_amd64.exe`

另外会额外上传一个校验文件：

- `SHA256SUMS.txt`
