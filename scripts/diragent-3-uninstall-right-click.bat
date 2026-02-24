@echo off
setlocal EnableExtensions

set "NO_PAUSE=0"
if /I "%~1"=="--no-pause" set "NO_PAUSE=1"

cd /d "%~dp0.."

echo [DirAgent] Uninstalling Explorer context menu...
powershell -NoProfile -ExecutionPolicy Bypass -File ".\scripts\uninstall.ps1" -RemoveAssets -RemoveConfig
if errorlevel 1 (
  echo [DirAgent][ERROR] Uninstall failed.
  goto :fail
)

echo [DirAgent] Uninstall completed.
if "%NO_PAUSE%"=="1" exit /b 0
pause
exit /b 0

:fail
if "%NO_PAUSE%"=="1" exit /b 1
pause
exit /b 1
