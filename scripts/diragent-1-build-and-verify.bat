@echo off
setlocal EnableExtensions

set "NO_PAUSE=0"
if /I "%~1"=="--no-pause" set "NO_PAUSE=1"

cd /d "%~dp0.."

echo [DirAgent] Step 1/2: go test ./...
go test ./...
if errorlevel 1 (
  echo [DirAgent][ERROR] go test failed.
  goto :fail
)

echo [DirAgent] Step 2/2: go build -trimpath -ldflags "-s -w -X main.version=1.0.0" -o diragent.exe ./cmd/diragent
go build -trimpath -ldflags "-s -w -X main.version=1.0.0" -o diragent.exe ./cmd/diragent
if errorlevel 1 (
  echo [DirAgent][ERROR] go build failed.
  goto :fail
)

if not exist ".\diragent.exe" (
  echo [DirAgent][ERROR] diragent.exe not found after build.
  goto :fail
)

echo [DirAgent] Build and verification completed successfully.
if "%NO_PAUSE%"=="1" exit /b 0
pause
exit /b 0

:fail
if "%NO_PAUSE%"=="1" exit /b 1
pause
exit /b 1
