@echo off
setlocal EnableExtensions

set "NO_PAUSE=0"
if /I "%~1"=="--no-pause" set "NO_PAUSE=1"

cd /d "%~dp0.."

if not exist ".\diragent.exe" (
  echo [DirAgent] diragent.exe not found. Building first...
  call ".\scripts\diragent-1-build-and-verify.bat" --no-pause
  if errorlevel 1 (
    echo [DirAgent][ERROR] Build failed. Install aborted.
    goto :fail
  )
)

echo [DirAgent] Installing Explorer context menu...
powershell -NoProfile -ExecutionPolicy Bypass -File ".\scripts\install.ps1"
if errorlevel 1 (
  echo [DirAgent][ERROR] Installation failed.
  goto :fail
)

echo [DirAgent] Installation completed.
echo [DirAgent] You can now right-click files/folders and use "Open in Codex/Claude (DirAgent)".
if "%NO_PAUSE%"=="1" exit /b 0
pause
exit /b 0

:fail
if "%NO_PAUSE%"=="1" exit /b 1
pause
exit /b 1
