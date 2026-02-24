package scripts

import (
	"os"
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
