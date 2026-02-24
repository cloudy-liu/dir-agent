package terminal

import (
	"strings"
	"testing"
)

func TestBuildWindowsTerminalUsesPowerShellWrapperAndTabPreferred(t *testing.T) {
	opts := LaunchOptions{
		PreferredTerminal: "windows-terminal",
		OpenMode:          "tab_preferred",
		WorkingDir:        `C:\work\repo`,
		CommandPath:       `C:\Users\cloudy\AppData\Roaming\npm\codex.ps1`,
		Args:              []string{"--model", "gpt-5"},
	}

	name, args, err := buildWindowsTerminal(opts)
	if err != nil {
		t.Fatalf("build windows terminal: %v", err)
	}
	if name != "wt.exe" {
		t.Fatalf("expected wt.exe, got %q", name)
	}
	if len(args) < 10 {
		t.Fatalf("expected richer wt args, got %#v", args)
	}
	if args[0] != "-w" || args[1] != "0" || args[2] != "new-tab" {
		t.Fatalf("expected tab preferred window args, got %#v", args[:3])
	}
	if args[3] != "-d" || args[4] != `C:\work\repo` {
		t.Fatalf("expected working dir args, got %#v", args[3:5])
	}
	if args[5] != "powershell.exe" || args[6] != "-NoExit" {
		t.Fatalf("expected powershell wrapper in wt args, got %#v", args[5:7])
	}
	script := args[len(args)-1]
	if !strings.Contains(script, "& 'C:\\Users\\cloudy\\AppData\\Roaming\\npm\\codex.ps1'") {
		t.Fatalf("expected wrapper script to invoke resolved command path, got %q", script)
	}
}

func TestBuildWindowsTerminalScriptAvoidsSemicolonCommandSeparator(t *testing.T) {
	opts := LaunchOptions{
		OpenMode:    "tab_preferred",
		WorkingDir:  `C:\work\repo`,
		CommandPath: `C:\Users\cloudy\AppData\Roaming\npm\codex.cmd`,
		Args:        []string{"--model", "gpt-5"},
	}

	_, args, err := buildWindowsTerminal(opts)
	if err != nil {
		t.Fatalf("build windows terminal: %v", err)
	}

	script := args[len(args)-1]
	if strings.Contains(script, ";") {
		t.Fatalf("script must not contain ';' because wt treats it as command separator: %q", script)
	}
}

func TestBuildWindowsTerminalUsesNewWindowMode(t *testing.T) {
	opts := LaunchOptions{
		OpenMode:    "new_window",
		WorkingDir:  `C:\work\repo`,
		CommandPath: `C:\tools\codex.cmd`,
	}

	_, args, err := buildWindowsTerminal(opts)
	if err != nil {
		t.Fatalf("build windows terminal: %v", err)
	}
	if args[0] != "-w" || args[1] != "new" || args[2] != "new-tab" {
		t.Fatalf("expected new window mode to target new window, got %#v", args[:3])
	}
}

func TestNormalizeOpenMode(t *testing.T) {
	if normalizeOpenMode("") != "tab_preferred" {
		t.Fatalf("expected empty open mode to default tab_preferred")
	}
	if normalizeOpenMode("new_window") != "new_window" {
		t.Fatalf("expected explicit new_window to be preserved")
	}
	if normalizeOpenMode("unknown-value") != "tab_preferred" {
		t.Fatalf("expected unknown open_mode to fallback tab_preferred")
	}
}
