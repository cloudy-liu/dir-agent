package terminal

import (
	"strings"
	"testing"
)

func TestBuildWindowsTerminalUsesPowerShellWrapperAndTabPreferred(t *testing.T) {
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return true }
	defer func() { isWindowsTerminalRunning = previous }()

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

func TestBuildWindowsTerminalTabPreferredUsesSingleTabWhenNoRunningWindow(t *testing.T) {
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return false }
	defer func() { isWindowsTerminalRunning = previous }()

	opts := LaunchOptions{
		OpenMode:    "tab_preferred",
		WorkingDir:  `C:\work\repo`,
		CommandPath: `C:\Users\cloudy\AppData\Roaming\npm\codex.cmd`,
	}

	_, args, err := buildWindowsTerminal(opts)
	if err != nil {
		t.Fatalf("build windows terminal: %v", err)
	}
	if len(args) < 4 {
		t.Fatalf("expected at least 4 args, got %#v", args)
	}
	if args[0] != "-w" || args[1] != "new" {
		t.Fatalf("expected tab_preferred without running window to use -w new, got %#v", args[:2])
	}
	if args[2] != "-d" || args[3] != `C:\work\repo` {
		t.Fatalf("expected direct single-tab launch with working dir, got %#v", args[:4])
	}
	if strings.Contains(strings.Join(args, " "), "new-tab") {
		t.Fatalf("expected no new-tab command when no running window, got %#v", args)
	}
}

func TestBuildWindowsTerminalScriptAvoidsSemicolonCommandSeparator(t *testing.T) {
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return true }
	defer func() { isWindowsTerminalRunning = previous }()

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
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return true }
	defer func() { isWindowsTerminalRunning = previous }()

	opts := LaunchOptions{
		OpenMode:    "new_window",
		WorkingDir:  `C:\work\repo`,
		CommandPath: `C:\tools\codex.cmd`,
	}

	_, args, err := buildWindowsTerminal(opts)
	if err != nil {
		t.Fatalf("build windows terminal: %v", err)
	}
	if len(args) < 4 {
		t.Fatalf("expected at least 4 args, got %#v", args)
	}
	if args[0] != "-w" || args[1] != "new" || args[2] != "-d" || args[3] != `C:\work\repo` {
		t.Fatalf("expected new window mode to launch single tab window, got %#v", args[:4])
	}
	if strings.Contains(strings.Join(args, " "), "new-tab") {
		t.Fatalf("expected no new-tab command in new_window mode, got %#v", args)
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
