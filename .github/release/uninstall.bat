@echo off
setlocal EnableExtensions
cd /d "%~dp0"

if not exist ".\diragent.exe" (
  echo [DirAgent][ERROR] diragent.exe not found in current folder.
  echo Please keep uninstall.bat next to diragent.exe.
  pause
  exit /b 1
)

echo [DirAgent] Uninstalling context menu and local assets/config...
powershell -NoProfile -ExecutionPolicy Bypass -File ".\scripts\uninstall.ps1" -BinaryPath ".\diragent.exe" -RemoveAssets -RemoveConfig
if errorlevel 1 (
  echo [DirAgent][ERROR] Uninstall failed.
  pause
  exit /b 1
)

echo [DirAgent] Uninstall completed.
pause
exit /b 0
