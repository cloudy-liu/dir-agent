@echo off
setlocal EnableExtensions
cd /d "%~dp0"

if not exist ".\diragent.exe" (
  echo [DirAgent][ERROR] diragent.exe not found in current folder.
  echo Please keep install.bat next to diragent.exe.
  pause
  exit /b 1
)

echo [DirAgent] Cleaning previous install (keep existing config)...
powershell -NoProfile -ExecutionPolicy Bypass -File ".\scripts\uninstall.ps1" -BinaryPath ".\diragent.exe" -RemoveAssets >nul 2>&1

echo [DirAgent] Installing context menu...
powershell -NoProfile -ExecutionPolicy Bypass -File ".\scripts\install.ps1" -BinaryPath ".\diragent.exe"
if errorlevel 1 (
  echo [DirAgent][ERROR] Install failed.
  pause
  exit /b 1
)

echo [DirAgent] Install completed.
echo [DirAgent] You can now right-click a folder and choose DirAgent.
pause
exit /b 0
