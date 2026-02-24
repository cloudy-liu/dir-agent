package terminal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildWindowsTerminalUsesPowerShellWrapperAndTabPreferred(t *testing.T) {
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return true }
	defer func() { isWindowsTerminalRunning = previous }()

	opts := LaunchOptions{
		PreferredTerminal:    "windows-terminal",
		OpenMode:             "tab_preferred",
		WorkingDir:           `C:\work\repo`,
		CommandPath:          `C:\Users\cloudy\AppData\Roaming\npm\codex.ps1`,
		WindowsTerminalShell: "powershell",
		Args:                 []string{"--model", "gpt-5"},
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

func TestBuildWindowsTerminalUsesConfiguredProfileAndCmdShell(t *testing.T) {
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return true }
	defer func() { isWindowsTerminalRunning = previous }()

	opts := LaunchOptions{
		OpenMode:               "tab_preferred",
		WorkingDir:             `C:\work\repo`,
		CommandPath:            `C:\Users\cloudy\AppData\Roaming\npm\codex.cmd`,
		WindowsTerminalProfile: "Cmder",
		WindowsTerminalShell:   "cmd",
		Args:                   []string{"--model", "gpt-5"},
	}

	_, args, err := buildWindowsTerminal(opts)
	if err != nil {
		t.Fatalf("build windows terminal: %v", err)
	}
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "new-tab") {
		t.Fatalf("expected new-tab when WT is running, got %#v", args)
	}
	if !strings.Contains(joined, "-p Cmder") {
		t.Fatalf("expected configured windows terminal profile to be applied, got %#v", args)
	}
	if !strings.Contains(joined, "cmd.exe /K") {
		t.Fatalf("expected cmd shell wrapper to be used, got %#v", args)
	}
}

func TestBuildWindowsTerminalCmdShellPrefersCmdWrapperForPs1Command(t *testing.T) {
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return true }
	defer func() { isWindowsTerminalRunning = previous }()

	tempDir := t.TempDir()
	ps1Path := filepath.Join(tempDir, "codex.ps1")
	cmdPath := filepath.Join(tempDir, "codex.cmd")
	writeTestFile(t, ps1Path)
	writeTestFile(t, cmdPath)

	opts := LaunchOptions{
		OpenMode:               "tab_preferred",
		WorkingDir:             `C:\work\repo`,
		CommandPath:            ps1Path,
		WindowsTerminalProfile: "Cmder",
		WindowsTerminalShell:   "cmd",
	}

	_, args, err := buildWindowsTerminal(opts)
	if err != nil {
		t.Fatalf("build windows terminal: %v", err)
	}
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, cmdPath) {
		t.Fatalf("expected cmd shell wrapper to invoke .cmd sibling when command path is .ps1, got %#v", args)
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

func writeTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("echo"), 0o644); err != nil {
		t.Fatalf("write test file %s: %v", path, err)
	}
}
