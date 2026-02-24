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
