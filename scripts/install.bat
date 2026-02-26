@echo off
setlocal EnableExtensions

set "NO_PAUSE=0"
if /I "%~1"=="--no-pause" set "NO_PAUSE=1"

cd /d "%~dp0.."

echo [DirAgent] Step 1/4: uninstall previous integration if present...
powershell -NoProfile -ExecutionPolicy Bypass -File ".\scripts\uninstall.ps1" -BinaryPath ".\diragent.exe" -RemoveAssets >nul 2>&1

echo [DirAgent] Step 2/4: go build -trimpath -ldflags "-s -w -X main.version=1.0.0" -o diragent.exe ./cmd/diragent
go build -trimpath -ldflags "-s -w -X main.version=1.0.0" -o diragent.exe ./cmd/diragent
if errorlevel 1 (
  echo [DirAgent][ERROR] go build for diragent.exe failed.
  goto :fail
)

echo [DirAgent] Step 3/4: go build -trimpath -ldflags "-s -w -X main.version=1.0.0" -o diragentw.exe ./cmd/diragentw
go build -trimpath -ldflags "-s -w -X main.version=1.0.0" -o diragentw.exe ./cmd/diragentw
if errorlevel 1 (
  echo [DirAgent][ERROR] go build for diragentw.exe failed.
  goto :fail
)

if not exist ".\diragent.exe" (
  echo [DirAgent][ERROR] diragent.exe not found after build.
  goto :fail
)
if not exist ".\diragentw.exe" (
  echo [DirAgent][ERROR] diragentw.exe not found after build.
  goto :fail
)

echo [DirAgent] Step 4/4: install explorer context menu...
powershell -NoProfile -ExecutionPolicy Bypass -File ".\scripts\install.ps1" -BinaryPath ".\diragent.exe"
if errorlevel 1 (
  echo [DirAgent][ERROR] install failed.
  goto :fail
)

echo [DirAgent] Install completed with latest local build.
if "%NO_PAUSE%"=="1" exit /b 0
pause
exit /b 0

:fail
if "%NO_PAUSE%"=="1" exit /b 1
pause
exit /b 1
