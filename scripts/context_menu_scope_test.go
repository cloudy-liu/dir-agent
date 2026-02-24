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

	if strings.Contains(text, `HKCU\Software\Classes\*\shell`) {
		t.Fatalf("install.ps1 should not register file context menu (HKCU\\Software\\Classes\\*\\shell)")
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

	if strings.Contains(text, `HKCU\Software\Classes\*\shell`) {
		t.Fatalf("uninstall.ps1 should not remove file context menu (HKCU\\Software\\Classes\\*\\shell)")
	}
	if !strings.Contains(text, `HKCU\Software\Classes\Directory\shell`) {
		t.Fatalf("uninstall.ps1 should remove directory context menu")
	}
	if !strings.Contains(text, `HKCU\Software\Classes\Directory\Background\shell`) {
		t.Fatalf("uninstall.ps1 should remove directory background context menu")
	}
}
