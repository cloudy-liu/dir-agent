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
		CommandPath:          `C:\path\to\codex.ps1`,
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
	if !strings.Contains(script, "& 'C:\\path\\to\\codex.ps1'") {
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
		CommandPath:            `C:\path\to\codex.cmd`,
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

func TestBuildWindowsTerminalUsesConfiguredProfileAndCmderShell(t *testing.T) {
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return true }
	defer func() { isWindowsTerminalRunning = previous }()

	opts := LaunchOptions{
		OpenMode:                 "tab_preferred",
		WorkingDir:               `C:\work\repo`,
		CommandPath:              `C:\path\to\codex.cmd`,
		WindowsTerminalProfile:   "Cmder",
		WindowsTerminalShell:     "cmder",
		WindowsTerminalCmderInit: `C:\path\to\cmder\vendor\init.bat`,
		Args:                     []string{"--model", "gpt-5"},
	}

	_, args, err := buildWindowsTerminal(opts)
	if err != nil {
		t.Fatalf("build windows terminal: %v", err)
	}
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "cmd.exe /K") {
		t.Fatalf("expected cmder shell wrapper to use cmd.exe /K, got %#v", args)
	}
	if !strings.Contains(joined, `C:\path\to\cmder\vendor\init.bat`) {
		t.Fatalf("expected cmder shell wrapper to include cmder init script, got %#v", args)
	}
}

func TestBuildWindowsTerminalCmderShellUsesCMDERRootFallback(t *testing.T) {
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return true }
	defer func() { isWindowsTerminalRunning = previous }()

	cmderRoot := filepath.Join(t.TempDir(), "cmder")
	initPath := filepath.Join(cmderRoot, "vendor", "init.bat")
	if err := os.MkdirAll(filepath.Dir(initPath), 0o755); err != nil {
		t.Fatalf("mkdir cmder vendor: %v", err)
	}
	writeTestFile(t, initPath)
	t.Setenv("CMDER_ROOT", cmderRoot)

	opts := LaunchOptions{
		OpenMode:               "tab_preferred",
		WorkingDir:             `C:\work\repo`,
		CommandPath:            `C:\path\to\codex.cmd`,
		WindowsTerminalShell:   "cmder",
		WindowsTerminalProfile: "Cmder",
	}

	_, args, err := buildWindowsTerminal(opts)
	if err != nil {
		t.Fatalf("build windows terminal: %v", err)
	}
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, initPath) {
		t.Fatalf("expected cmder shell to use CMDER_ROOT fallback init.bat, got %#v", args)
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
		CommandPath: `C:\path\to\codex.cmd`,
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

func TestNormalizeWindowsTerminalShell(t *testing.T) {
	if normalizeWindowsTerminalShell("") != "powershell" {
		t.Fatalf("expected empty windows terminal shell to default powershell")
	}
	if normalizeWindowsTerminalShell("cmd") != "cmd" {
		t.Fatalf("expected cmd shell to be preserved")
	}
	if normalizeWindowsTerminalShell("cmder") != "cmder" {
		t.Fatalf("expected cmder shell to be preserved")
	}
	if normalizeWindowsTerminalShell("unknown") != "powershell" {
		t.Fatalf("expected unknown windows terminal shell to fallback powershell")
	}
}

func TestBuildWindowsTerminalScriptAvoidsSemicolonCommandSeparator(t *testing.T) {
	previous := isWindowsTerminalRunning
	isWindowsTerminalRunning = func() bool { return true }
	defer func() { isWindowsTerminalRunning = previous }()

	opts := LaunchOptions{
		OpenMode:    "tab_preferred",
		WorkingDir:  `C:\work\repo`,
		CommandPath: `C:\path\to\codex.cmd`,
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

func TestBuildPowerShellScriptActivatesLocalVenv(t *testing.T) {
	tempDir := t.TempDir()
	activatePath := filepath.Join(tempDir, ".venv", "Scripts", "Activate.ps1")
	if err := os.MkdirAll(filepath.Dir(activatePath), 0o755); err != nil {
		t.Fatalf("mkdir venv scripts: %v", err)
	}
	writeTestFile(t, activatePath)

	opts := LaunchOptions{
		WorkingDir:  tempDir,
		CommandPath: `C:\path\to\codex.ps1`,
		Args:        []string{"--model", "gpt-5"},
	}

	script := buildPowerShellScript(opts)
	if !strings.Contains(script, ". "+psQuote(activatePath)) {
		t.Fatalf("expected powershell script to activate local venv, got %q", script)
	}
}

func TestBuildCmdScriptActivatesLocalVenv(t *testing.T) {
	tempDir := t.TempDir()
	activatePath := filepath.Join(tempDir, ".venv", "Scripts", "activate.bat")
	if err := os.MkdirAll(filepath.Dir(activatePath), 0o755); err != nil {
		t.Fatalf("mkdir venv scripts: %v", err)
	}
	writeTestFile(t, activatePath)

	opts := LaunchOptions{
		WorkingDir:  tempDir,
		CommandPath: `C:\path\to\codex.cmd`,
		Args:        []string{"--model", "gpt-5"},
	}

	script := buildCmdScript(opts)
	expectedPrefix := "call " + cmdQuote(activatePath) + " && "
	if !strings.Contains(script, expectedPrefix) {
		t.Fatalf("expected cmd script to activate local venv, got %q", script)
	}
}

func TestBuildXTerminalEmulatorActivatesLocalVenv(t *testing.T) {
	tempDir := t.TempDir()
	activatePath := filepath.Join(tempDir, ".venv", "bin", "activate")
	if err := os.MkdirAll(filepath.Dir(activatePath), 0o755); err != nil {
		t.Fatalf("mkdir venv bin: %v", err)
	}
	writeTestFile(t, activatePath)

	opts := LaunchOptions{
		WorkingDir:  tempDir,
		CommandPath: "/usr/local/bin/codex",
		Args:        []string{"--model", "gpt-5"},
	}

	_, args, err := buildXTerminalEmulator(opts)
	if err != nil {
		t.Fatalf("build x-terminal-emulator: %v", err)
	}
	command := args[len(args)-1]
	if !strings.Contains(command, ". "+shQuote(activatePath)) {
		t.Fatalf("expected x-terminal command to activate local venv, got %q", command)
	}
}

func writeTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("echo"), 0o644); err != nil {
		t.Fatalf("write test file %s: %v", path, err)
	}
}
