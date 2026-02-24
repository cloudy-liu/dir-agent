param(
    [string]$BinaryPath = "",
    [switch]$RemoveAssets,
    [switch]$RemoveConfig
)

$ErrorActionPreference = "Stop"

function Resolve-BinaryPath {
    param([string]$ProvidedPath)

    if ($ProvidedPath -and (Test-Path $ProvidedPath)) {
        return (Resolve-Path $ProvidedPath).Path
    }

    $scriptDir = Split-Path -Parent $PSCommandPath
    $localBin = Join-Path (Split-Path -Parent $scriptDir) "diragent.exe"
    if (Test-Path $localBin) {
        return (Resolve-Path $localBin).Path
    }

    $globalCmd = Get-Command diragent.exe -ErrorAction SilentlyContinue
    if ($globalCmd) {
        return $globalCmd.Source
    }

    $globalCmdNoExt = Get-Command diragent -ErrorAction SilentlyContinue
    if ($globalCmdNoExt) {
        return $globalCmdNoExt.Source
    }

    return ""
}

function Remove-ContextMenuEntry {
    param(
        [string]$BaseKey,
        [string]$MenuKey
    )
    $menuPath = "$BaseKey\$MenuKey"
    cmd /c "reg query `"$menuPath`" >nul 2>&1"
    if ($LASTEXITCODE -eq 0) {
        cmd /c "reg delete `"$menuPath`" /f >nul 2>&1"
    }
}

$entries = @(
    "HKCU\Software\Classes\Directory\shell",
    "HKCU\Software\Classes\Directory\Background\shell"
)

foreach ($base in $entries) {
    Remove-ContextMenuEntry -BaseKey $base -MenuKey "DirAgentOpenInCodex"
    Remove-ContextMenuEntry -BaseKey $base -MenuKey "DirAgentOpenInClaude"
}

if ($RemoveAssets) {
    $exe = Resolve-BinaryPath -ProvidedPath $BinaryPath
    if ($exe -eq "") {
        Write-Warning "diragent binary not found; skipped removing embedded assets."
    }
    else {
        if ($RemoveConfig) {
            & $exe uninstall-assets --remove-config | Out-Null
        }
        else {
            & $exe uninstall-assets | Out-Null
        }
    }
}

Write-Host "Removed DirAgent context menu entries."
