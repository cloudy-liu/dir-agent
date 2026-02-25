package terminal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type LaunchOptions struct {
	PreferredTerminal        string
	OpenMode                 string
	WorkingDir               string
	CommandPath              string
	Args                     []string
	WindowsTerminalProfile   string
	WindowsTerminalShell     string
	WindowsTerminalCmderInit string
}

type candidate struct {
	ID      string
	Binary  string
	Builder func(LaunchOptions) (string, []string, error)
}

var ErrNoTerminalFound = errors.New("no supported terminal found")
var isWindowsTerminalRunning = detectWindowsTerminalRunning
var isWezTermRunning = detectWezTermRunning

func FindExecutable(binary string) (string, error) {
	return exec.LookPath(binary)
}

func LaunchInTerminal(opts LaunchOptions) error {
	candidates := terminalCandidates()
	ordered := prioritize(candidates, normalizeID(opts.PreferredTerminal))

	for _, item := range ordered {
		if _, err := FindExecutable(item.Binary); err != nil {
			continue
		}
		name, args, err := item.Builder(opts)
		if err != nil {
			return err
		}
		if err := launchTerminalCommand(item, opts, name, args); err != nil {
			continue
		}
		return nil
	}

	return ErrNoTerminalFound
}

func launchTerminalCommand(item candidate, opts LaunchOptions, name string, args []string) error {
	if shouldRunWezTermSpawnWithFallback(item, opts, args) {
		cmd := exec.Command(name, args...)
		if err := cmd.Run(); err == nil {
			return nil
		}
		fallbackOpts := opts
		fallbackOpts.OpenMode = "new_window"
		fallbackName, fallbackArgs, fallbackErr := buildWezTermWindows(fallbackOpts)
		if fallbackErr != nil {
			return fallbackErr
		}
		fallbackCmd := exec.Command(fallbackName, fallbackArgs...)
		return fallbackCmd.Start()
	}

	cmd := exec.Command(name, args...)
	return cmd.Start()
}

func shouldRunWezTermSpawnWithFallback(item candidate, opts LaunchOptions, args []string) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	if normalizeID(item.ID) != "wezterm" {
		return false
	}
	if normalizeOpenMode(opts.OpenMode) != "tab_preferred" {
		return false
	}
	return len(args) >= 2 && args[0] == "cli" && args[1] == "spawn"
}

func prioritize(candidates []candidate, preferred string) []candidate {
	if preferred == "" {
		return candidates
	}

	prioritized := make([]candidate, 0, len(candidates))
	for _, c := range candidates {
		if normalizeID(c.ID) == preferred {
			prioritized = append(prioritized, c)
			break
		}
	}
	for _, c := range candidates {
		if normalizeID(c.ID) != preferred {
			prioritized = append(prioritized, c)
		}
	}
	return prioritized
}

func normalizeID(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func terminalCandidates() []candidate {
	switch runtime.GOOS {
	case "windows":
		return []candidate{
			{ID: "windows-terminal", Binary: "wt.exe", Builder: buildWindowsTerminal},
			{ID: "wezterm", Binary: "wezterm.exe", Builder: buildWezTermWindows},
			{ID: "powershell", Binary: "powershell.exe", Builder: buildPowerShellTerminal},
		}
	case "darwin":
		return []candidate{
			{ID: "terminal.app", Binary: "osascript", Builder: buildMacTerminalApp},
			{ID: "iterm2", Binary: "osascript", Builder: buildMacITerm},
		}
	default:
		return []candidate{
			{ID: "x-terminal-emulator", Binary: "x-terminal-emulator", Builder: buildXTerminalEmulator},
			{ID: "gnome-terminal", Binary: "gnome-terminal", Builder: buildGnomeTerminal},
			{ID: "konsole", Binary: "konsole", Builder: buildKonsole},
			{ID: "xterm", Binary: "xterm", Builder: buildXTerm},
		}
	}
}

func buildWindowsTerminal(opts LaunchOptions) (string, []string, error) {
	commandArgs := buildWindowsTerminalCommandArgs(opts)
	profile := strings.TrimSpace(opts.WindowsTerminalProfile)

	if normalizeOpenMode(opts.OpenMode) == "new_window" {
		args := []string{"-w", "new"}
		if profile != "" {
			args = append(args, "-p", profile)
		}
		args = append(args, "-d", opts.WorkingDir)
		args = append(args, commandArgs...)
		return "wt.exe", args, nil
	}

	if isWindowsTerminalRunning() {
		args := []string{"-w", "0", "new-tab"}
		if profile != "" {
			args = append(args, "-p", profile)
		}
		args = append(args, "-d", opts.WorkingDir)
		args = append(args, commandArgs...)
		return "wt.exe", args, nil
	}

	args := []string{"-w", "new"}
	if profile != "" {
		args = append(args, "-p", profile)
	}
	args = append(args, "-d", opts.WorkingDir)
	args = append(args, commandArgs...)
	return "wt.exe", args, nil
}

func buildWezTermWindows(opts LaunchOptions) (string, []string, error) {
	script := buildPowerShellScript(opts)
	if normalizeOpenMode(opts.OpenMode) == "tab_preferred" && isWezTermRunning() {
		args := []string{"cli", "spawn", "--new-tab", "--cwd", opts.WorkingDir, "--", "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-Command", script}
		return "wezterm.exe", args, nil
	}

	args := []string{"start", "--cwd", opts.WorkingDir, "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-Command", script}
	return "wezterm.exe", args, nil
}

func buildPowerShellTerminal(opts LaunchOptions) (string, []string, error) {
	script := buildPowerShellScript(opts)
	return "powershell.exe", []string{"-NoExit", "-Command", script}, nil
}

func buildMacTerminalApp(opts LaunchOptions) (string, []string, error) {
	command := buildPosixShellCommand(opts)
	appleScript := fmt.Sprintf(`tell application "Terminal" to do script %q`, command)
	return "osascript", []string{"-e", appleScript}, nil
}

func buildMacITerm(opts LaunchOptions) (string, []string, error) {
	command := buildPosixShellCommand(opts)
	appleScript := fmt.Sprintf(`tell application "iTerm"
create window with default profile
tell current session of current window
write text %q
end tell
end tell`, command)
	return "osascript", []string{"-e", appleScript}, nil
}

func buildXTerminalEmulator(opts LaunchOptions) (string, []string, error) {
	command := buildPosixShellCommand(opts)
	return "x-terminal-emulator", []string{"-e", "sh", "-lc", command}, nil
}

func buildGnomeTerminal(opts LaunchOptions) (string, []string, error) {
	command := buildPosixShellCommand(opts)
	args := []string{"--working-directory", opts.WorkingDir, "--", "sh", "-lc", command}
	return "gnome-terminal", args, nil
}

func buildKonsole(opts LaunchOptions) (string, []string, error) {
	command := buildPosixShellCommand(opts)
	args := []string{"--workdir", opts.WorkingDir, "-e", "sh", "-lc", command}
	return "konsole", args, nil
}

func buildXTerm(opts LaunchOptions) (string, []string, error) {
	command := buildPosixShellCommand(opts)
	return "xterm", []string{"-e", "sh", "-lc", command}, nil
}

func buildPowerShellScript(opts LaunchOptions) string {
	// Use a newline separator to avoid wt.exe treating ';' as a command delimiter.
	script := "Set-Location -LiteralPath " + psQuote(opts.WorkingDir)
	if activate := resolveWindowsVenvActivatePs1(opts.WorkingDir); activate != "" {
		script += "\n. " + psQuote(activate)
	}
	script += "\n& " + psQuote(opts.CommandPath)
	for _, arg := range opts.Args {
		script += " " + psQuote(arg)
	}
	return script
}

func buildCmdScript(opts LaunchOptions) string {
	commandPath := opts.CommandPath
	if strings.EqualFold(filepath.Ext(commandPath), ".ps1") {
		cmdPath := strings.TrimSuffix(commandPath, filepath.Ext(commandPath)) + ".cmd"
		if _, err := os.Stat(cmdPath); err == nil {
			commandPath = cmdPath
		}
	}

	command := cmdQuote(commandPath)
	for _, arg := range opts.Args {
		command += " " + cmdQuote(arg)
	}
	if activate := resolveWindowsVenvActivateCmd(opts.WorkingDir); activate != "" {
		command = "call " + cmdQuote(activate) + " && " + command
	}
	return command
}

func buildWindowsTerminalCommandArgs(opts LaunchOptions) []string {
	switch normalizeWindowsTerminalShell(opts.WindowsTerminalShell) {
	case "cmd":
		script := buildCmdScript(opts)
		return []string{"cmd.exe", "/K", script}
	case "cmder":
		script := buildCmderScript(opts)
		return []string{"cmd.exe", "/K", script}
	default:
		script := buildPowerShellScript(opts)
		return []string{"powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-Command", script}
	}
}

func buildCmderScript(opts LaunchOptions) string {
	command := buildCmdScript(opts)
	initPath := resolveCmderInitPath(opts.WindowsTerminalCmderInit)
	if initPath == "" {
		return command
	}
	return cmdQuote(initPath) + " && " + command
}

func resolveCmderInitPath(configured string) string {
	configured = strings.TrimSpace(configured)
	if configured != "" {
		return configured
	}

	cmderRoot := strings.TrimSpace(os.Getenv("CMDER_ROOT"))
	if cmderRoot == "" {
		return ""
	}
	candidate := filepath.Join(cmderRoot, "vendor", "init.bat")
	if _, err := os.Stat(candidate); err != nil {
		return ""
	}
	return candidate
}

func normalizeWindowsTerminalShell(value string) string {
	switch normalizeID(value) {
	case "", "powershell":
		return "powershell"
	case "cmd":
		return "cmd"
	case "cmder":
		return "cmder"
	default:
		return "powershell"
	}
}

func normalizeOpenMode(openMode string) string {
	switch normalizeID(openMode) {
	case "", "tab_preferred":
		return "tab_preferred"
	case "new_window":
		return "new_window"
	default:
		return "tab_preferred"
	}
}

func psQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "''")
	return "'" + escaped + "'"
}

func shQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", `'"'"'`)
	return "'" + escaped + "'"
}

func cmdQuote(value string) string {
	escaped := strings.ReplaceAll(value, `"`, `""`)
	return `"` + escaped + `"`
}

func buildPosixShellCommand(opts LaunchOptions) string {
	command := "cd " + shQuote(opts.WorkingDir)
	if activate := resolvePosixVenvActivate(opts.WorkingDir); activate != "" {
		command += "; . " + shQuote(activate)
	}
	command += "; exec " + shQuote(opts.CommandPath)
	for _, arg := range opts.Args {
		command += " " + shQuote(arg)
	}
	return command
}

func resolveWindowsVenvActivatePs1(workingDir string) string {
	activate := filepath.Join(workingDir, ".venv", "Scripts", "Activate.ps1")
	if _, err := os.Stat(activate); err != nil {
		return ""
	}
	return activate
}

func resolveWindowsVenvActivateCmd(workingDir string) string {
	activate := filepath.Join(workingDir, ".venv", "Scripts", "activate.bat")
	if _, err := os.Stat(activate); err != nil {
		return ""
	}
	return activate
}

func resolvePosixVenvActivate(workingDir string) string {
	activate := filepath.Join(workingDir, ".venv", "bin", "activate")
	if _, err := os.Stat(activate); err != nil {
		return ""
	}
	return activate
}

func detectWindowsTerminalRunning() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	output, err := exec.Command("tasklist", "/FI", "IMAGENAME eq WindowsTerminal.exe").Output()
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(string(output)), "windowsterminal.exe")
}

func detectWezTermRunning() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	output, err := exec.Command("tasklist", "/FI", "IMAGENAME eq wezterm-gui.exe").Output()
	if err == nil && strings.Contains(strings.ToLower(string(output)), "wezterm-gui.exe") {
		return true
	}

	output, err = exec.Command("tasklist", "/FI", "IMAGENAME eq wezterm.exe").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(output)), "wezterm.exe")
}
