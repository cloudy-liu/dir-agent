package scripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWindowsInstallScriptRegistersDirectoryOnlyContextMenu(t *testing.T) {
	content, err := os.ReadFile("install.ps1")
	if err != nil {
		t.Fatalf("read install.ps1: %v", err)
	}
	text := string(content)

	if strings.Contains(text, `@{ Base = "HKCU\Software\Classes\*\shell"; Placeholder = "%1" }`) {
		t.Fatalf("install.ps1 must not register file context menu")
	}
	if !strings.Contains(text, `Remove-ContextMenuEntry -BaseKey "HKCU\Software\Classes\*\shell" -MenuKey "DirAgentOpenInCodex"`) {
		t.Fatalf("install.ps1 should clean legacy file context menu entry for codex")
	}
	if !strings.Contains(text, `Remove-ContextMenuEntry -BaseKey "HKCU\Software\Classes\*\shell" -MenuKey "DirAgentOpenInClaude"`) {
		t.Fatalf("install.ps1 should clean legacy file context menu entry for claude")
	}
	if !strings.Contains(text, `HKCU\Software\Classes\Directory\shell`) {
		t.Fatalf("install.ps1 should register directory context menu")
	}
	if !strings.Contains(text, `HKCU\Software\Classes\Directory\Background\shell`) {
		t.Fatalf("install.ps1 should register directory background context menu")
	}
	if !strings.Contains(text, `-Title "Open in Claude Code (DirAgent)"`) {
		t.Fatalf("install.ps1 should expose Claude Code menu title")
	}
	if strings.Contains(text, `-Title "Open in Claude (DirAgent)"`) {
		t.Fatalf("install.ps1 should not use legacy Claude menu title")
	}
}

func TestWindowsUninstallScriptRemovesDirectoryOnlyContextMenu(t *testing.T) {
	content, err := os.ReadFile("uninstall.ps1")
	if err != nil {
		t.Fatalf("read uninstall.ps1: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, `HKCU\Software\Classes\*\shell`) {
		t.Fatalf("uninstall.ps1 should remove legacy file context menu if present")
	}
	if !strings.Contains(text, `HKCU\Software\Classes\Directory\shell`) {
		t.Fatalf("uninstall.ps1 should remove directory context menu")
	}
	if !strings.Contains(text, `HKCU\Software\Classes\Directory\Background\shell`) {
		t.Fatalf("uninstall.ps1 should remove directory background context menu")
	}
}

func TestUnixInstallScriptUsesClaudeCodeLabel(t *testing.T) {
	content, err := os.ReadFile("install.sh")
	if err != nil {
		t.Fatalf("read install.sh: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "Open in Claude Code (DirAgent)") {
		t.Fatalf("install.sh should use Claude Code label")
	}
	if strings.Contains(text, "Open in Claude (DirAgent)") {
		t.Fatalf("install.sh should not use legacy Claude label")
	}
}

func TestUnixUninstallScriptCleansLegacyAndCurrentClaudeApps(t *testing.T) {
	content, err := os.ReadFile("uninstall.sh")
	if err != nil {
		t.Fatalf("read uninstall.sh: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "Open in Claude Code (DirAgent).app") {
		t.Fatalf("uninstall.sh should remove current Claude Code app name")
	}
	if !strings.Contains(text, "Open in Claude (DirAgent).app") {
		t.Fatalf("uninstall.sh should also remove legacy Claude app name")
	}
}

func TestWindowsInstallScriptPrefersGuiLauncherBinary(t *testing.T) {
	content, err := os.ReadFile("install.ps1")
	if err != nil {
		t.Fatalf("read install.ps1: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "function Resolve-LauncherPath") {
		t.Fatalf("install.ps1 should define Resolve-LauncherPath to prefer diragentw.exe")
	}
	if !strings.Contains(text, "$launcher = Resolve-LauncherPath -ExePath $exe") {
		t.Fatalf("install.ps1 should resolve launcher path from the primary executable path")
	}
	if !strings.Contains(text, `$commandValue = "`+"`"+`"$LauncherPath`+"`"+`" launch --tool $Tool --path `+"`"+`"$TargetPlaceholder`+"`"+`""`) {
		t.Fatalf("install.ps1 should build menu command with launcher path")
	}
}

func TestWindowsLocalBatEntrypointsOnlyInstallAndUninstall(t *testing.T) {
	batFiles, err := filepath.Glob("*.bat")
	if err != nil {
		t.Fatalf("glob bat files: %v", err)
	}
	if len(batFiles) != 2 {
		t.Fatalf("scripts directory should contain exactly 2 bat files, got %d: %v", len(batFiles), batFiles)
	}
	expected := map[string]bool{
		"install.bat":   false,
		"uninstall.bat": false,
	}
	for _, file := range batFiles {
		name := filepath.Base(file)
		if _, ok := expected[name]; !ok {
			t.Fatalf("unexpected bat entrypoint in scripts: %s", name)
		}
		expected[name] = true
	}
	for name, found := range expected {
		if !found {
			t.Fatalf("missing required bat entrypoint: %s", name)
		}
	}
}

func TestWindowsInstallBatBuildsLatestAndInstalls(t *testing.T) {
	content, err := os.ReadFile("install.bat")
	if err != nil {
		t.Fatalf("read install.bat: %v", err)
	}
	text := string(content)

	uninstallIdx := strings.Index(text, `uninstall.ps1`)
	buildDiragentIdx := strings.Index(text, `go build -trimpath -ldflags "-s -w -X main.version=1.0.0" -o diragent.exe ./cmd/diragent`)
	buildDiragentwIdx := strings.Index(text, `go build -trimpath -ldflags "-s -w -X main.version=1.0.0" -o diragentw.exe ./cmd/diragentw`)
	installIdx := strings.Index(text, `.\scripts\install.ps1`)
	if uninstallIdx < 0 {
		t.Fatalf("install.bat should uninstall previous integration first")
	}
	if buildDiragentIdx < 0 || buildDiragentwIdx < 0 {
		t.Fatalf("install.bat should build both diragent.exe and diragentw.exe")
	}
	if installIdx < 0 {
		t.Fatalf("install.bat should invoke install.ps1 after building")
	}
	if !(uninstallIdx < buildDiragentIdx && buildDiragentIdx < buildDiragentwIdx && buildDiragentwIdx < installIdx) {
		t.Fatalf("install.bat should run uninstall -> build diragent -> build diragentw -> install in order")
	}
}

func TestWindowsUninstallBatOnlyUninstalls(t *testing.T) {
	content, err := os.ReadFile("uninstall.bat")
	if err != nil {
		t.Fatalf("read uninstall.bat: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, `uninstall.ps1`) {
		t.Fatalf("uninstall.bat should invoke uninstall.ps1")
	}
	if strings.Contains(text, `go build`) {
		t.Fatalf("uninstall.bat should not build binaries")
	}
	if strings.Contains(text, `.\scripts\install.ps1`) {
		t.Fatalf("uninstall.bat should not invoke install.ps1")
	}
}

func TestReleaseWorkflowBuildsGuiLauncherForWindows(t *testing.T) {
	content, err := os.ReadFile("../.github/workflows/release.yml")
	if err != nil {
		t.Fatalf("read release workflow: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "diragentw_") {
		t.Fatalf("release workflow should build and name diragentw windows artifact")
	}
}

func TestReleasePackageScriptIncludesGuiLauncherInWindowsBundle(t *testing.T) {
	content, err := os.ReadFile("../.github/scripts/package-release.sh")
	if err != nil {
		t.Fatalf("read package-release.sh: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "diragentw_") {
		t.Fatalf("package-release.sh should look up diragentw artifacts")
	}
	if !strings.Contains(text, "diragentw.exe") {
		t.Fatalf("package-release.sh should copy diragentw.exe into windows package")
	}
}
