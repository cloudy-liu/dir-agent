package terminal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"dir-agent/internal/proc"
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
	WindowsWezTermShell      string
	WindowsWezTermCmderInit  string
	WezTermWindowID          string
}

type candidate struct {
	ID      string
	Binary  string
	Builder func(LaunchOptions) (string, []string, error)
}

var ErrNoTerminalFound = errors.New("no supported terminal found")
var isWindowsTerminalRunning = detectWindowsTerminalRunning
var isWezTermRunning = detectWezTermRunning
var detectWezTermWindowID = queryWezTermWindowID

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
		cmd := buildWezTermCLICommand(name, args...)
		if err := cmd.Run(); err == nil {
			return nil
		}
		fallbackOpts := opts
		fallbackOpts.OpenMode = "new_window"
		fallbackName, fallbackArgs, fallbackErr := buildWezTermWindows(fallbackOpts)
		if fallbackErr != nil {
			return fallbackErr
		}
		fallbackCmd := hiddenCommand(fallbackName, fallbackArgs...)
		return fallbackCmd.Start()
	}

	cmd := hiddenCommand(name, args...)
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
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "windows_terminal":
		return "windows-terminal"
	default:
		return normalized
	}
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
	profile := resolveWindowsTerminalProfile(opts)

	if normalizeOpenMode(opts.OpenMode) == "new_window" {
		args := []string{"-w", "new"}
		if profile != "" {
			args = append(args, "-p", profile)
		}
		args = append(args, "-d", opts.WorkingDir)
		args = append(args, commandArgs...)
		return "wt.exe", args, nil
	}

	args := []string{"-w", "0", "new-tab"}
	if profile != "" {
		args = append(args, "-p", profile)
	}
	args = append(args, "-d", opts.WorkingDir)
	args = append(args, commandArgs...)
	return "wt.exe", args, nil
}

func resolveWindowsTerminalProfile(opts LaunchOptions) string {
	profile := strings.TrimSpace(opts.WindowsTerminalProfile)
	if profile != "" {
		return profile
	}
	if normalizeWindowsShell(opts.WindowsTerminalShell) == "cmder" {
		return "Cmder"
	}
	return ""
}

func buildWezTermWindows(opts LaunchOptions) (string, []string, error) {
	commandArgs := buildWezTermCommandArgs(opts)
	if normalizeOpenMode(opts.OpenMode) == "tab_preferred" && isWezTermRunning() {
		windowID := strings.TrimSpace(opts.WezTermWindowID)
		if windowID == "" {
			windowID = detectWezTermWindowID()
		}
		if windowID != "" {
			args := []string{"cli", "spawn", "--window-id", windowID, "--cwd", opts.WorkingDir, "--"}
			args = append(args, commandArgs...)
			return "wezterm.exe", args, nil
		}
	}

	args := []string{"start", "--cwd", opts.WorkingDir}
	args = append(args, commandArgs...)
	return "wezterm-gui.exe", args, nil
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
	commandPath := resolveWindowsCmdCommandPath(opts.CommandPath)

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
	return buildWindowsShellCommandArgs(opts, opts.WindowsTerminalShell, opts.WindowsTerminalCmderInit)
}

func buildWezTermCommandArgs(opts LaunchOptions) []string {
	return buildWindowsShellCommandArgs(opts, opts.WindowsWezTermShell, opts.WindowsWezTermCmderInit)
}

func buildWindowsShellCommandArgs(opts LaunchOptions, shell string, configuredCmderInit string) []string {
	switch normalizeWindowsShell(shell) {
	case "cmd":
		script := buildCmdScript(opts)
		return []string{"cmd.exe", "/K", script}
	case "cmder":
		return buildCmderCommandArgs(opts, configuredCmderInit)
	default:
		script := buildPowerShellScript(opts)
		return []string{"powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-Command", script}
	}
}

func buildCmderCommandArgs(opts LaunchOptions, configuredInit string) []string {
	args := []string{"cmd.exe", "/K"}
	initPath := resolveCmderInitPath(configuredInit)
	if initPath != "" {
		args = append(args, "call", initPath, "&&")
	}
	args = append(args, buildCmdCommandTokens(opts)...)
	return args
}

func buildCmdCommandTokens(opts LaunchOptions) []string {
	commandPath := resolveWindowsCmdCommandPath(opts.CommandPath)
	args := make([]string, 0, len(opts.Args)+4)
	if activate := resolveWindowsVenvActivateCmd(opts.WorkingDir); activate != "" {
		args = append(args, "call", activate, "&&")
	}
	args = append(args, commandPath)
	args = append(args, opts.Args...)
	return args
}

func resolveWindowsCmdCommandPath(commandPath string) string {
	if strings.EqualFold(filepath.Ext(commandPath), ".ps1") {
		cmdPath := strings.TrimSuffix(commandPath, filepath.Ext(commandPath)) + ".cmd"
		if _, err := os.Stat(cmdPath); err == nil {
			return cmdPath
		}
	}
	return commandPath
}

func resolveCmderInitPath(configured string) string {
	configured = strings.TrimSpace(configured)
	if configured != "" {
		return normalizeWindowsPath(configured)
	}

	cmderRoot := strings.TrimSpace(os.Getenv("CMDER_ROOT"))
	if cmderRoot == "" {
		return ""
	}
	candidate := filepath.Join(cmderRoot, "vendor", "init.bat")
	if _, err := os.Stat(candidate); err != nil {
		return ""
	}
	return normalizeWindowsPath(candidate)
}

func normalizeWindowsShell(value string) string {
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

func normalizeWindowsPath(value string) string {
	return filepath.Clean(strings.ReplaceAll(value, "/", `\`))
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

	output, err := hiddenCommand("tasklist", "/FI", "IMAGENAME eq WindowsTerminal.exe").Output()
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(string(output)), "windowsterminal.exe")
}

func detectWezTermRunning() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	output, err := hiddenCommand("tasklist", "/FI", "IMAGENAME eq wezterm-gui.exe").Output()
	if err == nil && strings.Contains(strings.ToLower(string(output)), "wezterm-gui.exe") {
		return true
	}

	output, err = hiddenCommand("tasklist", "/FI", "IMAGENAME eq wezterm.exe").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(output)), "wezterm.exe")
}

func queryWezTermWindowID() string {
	output, err := buildWezTermCLICommand("wezterm.exe", "cli", "list", "--format", "json").Output()
	if err != nil {
		return ""
	}

	type weztermListEntry struct {
		WindowID json.RawMessage `json:"window_id"`
	}

	var entries []weztermListEntry
	if err := json.Unmarshal(output, &entries); err != nil {
		return ""
	}

	for _, entry := range entries {
		windowID := parseWezTermWindowID(entry.WindowID)
		if windowID != "" {
			return windowID
		}
	}

	return ""
}

func buildWezTermCLICommand(name string, args ...string) *exec.Cmd {
	cmd := hiddenCommand(name, args...)
	if socket := resolveWezTermUnixSocketHint(); socket != "" {
		cmd.Env = append(os.Environ(), "WEZTERM_UNIX_SOCKET="+socket)
	}
	return cmd
}

func hiddenCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	proc.ApplyNoWindow(cmd)
	return cmd
}

func resolveWezTermUnixSocketHint() string {
	homeDir, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(homeDir) == "" {
		return ""
	}

	pattern := filepath.Join(homeDir, ".local", "share", "wezterm", "gui-sock-*")
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return ""
	}

	var newestPath string
	var newestModTime int64
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil || info.IsDir() {
			continue
		}
		modTime := info.ModTime().UnixNano()
		if newestPath == "" || modTime > newestModTime {
			newestPath = match
			newestModTime = modTime
		}
	}

	return newestPath
}

func parseWezTermWindowID(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	var numericID int64
	if err := json.Unmarshal(raw, &numericID); err == nil {
		return fmt.Sprintf("%d", numericID)
	}

	var stringID string
	if err := json.Unmarshal(raw, &stringID); err == nil {
		return strings.TrimSpace(stringID)
	}

	return ""
}
