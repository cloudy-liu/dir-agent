param(
    [string]$BinaryPath = ""
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

    throw "diragent binary not found. Provide -BinaryPath."
}

function Set-ContextMenuEntry {
    param(
        [string]$BaseKey,
        [string]$MenuKey,
        [string]$Title,
        [string]$Tool,
        [string]$IconPath,
        [string]$TargetPlaceholder,
        [string]$ExePath
    )

    $menuPath = "$BaseKey\$MenuKey"
    $commandPath = "$menuPath\command"
    $commandValue = "`"$ExePath`" launch --tool $Tool --path `"$TargetPlaceholder`""
    reg add $menuPath /ve /d $Title /f | Out-Null
    reg add $menuPath /v Icon /d $IconPath /f | Out-Null
    reg add $commandPath /ve /d $commandValue /f | Out-Null
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

$exe = Resolve-BinaryPath -ProvidedPath $BinaryPath
$installResultRaw = & $exe install-assets
$installResult = $installResultRaw | ConvertFrom-Json
$dataPath = $installResult.data_path

$codexIcon = Join-Path $dataPath "assets/icons/windows/codex.ico"
$claudeIcon = Join-Path $dataPath "assets/icons/windows/claude.ico"

$entries = @(
    @{ Base = "HKCU\Software\Classes\Directory\shell"; Placeholder = "%1" },
    @{ Base = "HKCU\Software\Classes\*\shell"; Placeholder = "%1" },
    @{ Base = "HKCU\Software\Classes\Directory\Background\shell"; Placeholder = "%V" }
)

foreach ($entry in $entries) {
    Remove-ContextMenuEntry -BaseKey $entry.Base -MenuKey "DirAgentOpenInCodex"
    Remove-ContextMenuEntry -BaseKey $entry.Base -MenuKey "DirAgentOpenInClaude"

    Set-ContextMenuEntry -BaseKey $entry.Base -MenuKey "DirAgentOpenInCodex" -Title "Open in Codex (DirAgent)" -Tool "codex" -IconPath $codexIcon -TargetPlaceholder $entry.Placeholder -ExePath $exe
    Set-ContextMenuEntry -BaseKey $entry.Base -MenuKey "DirAgentOpenInClaude" -Title "Open in Claude (DirAgent)" -Tool "claude" -IconPath $claudeIcon -TargetPlaceholder $entry.Placeholder -ExePath $exe
}

Write-Host "Installed DirAgent context menu entries with icons."
Write-Host "Binary: $exe"
Write-Host "Data path: $dataPath"
