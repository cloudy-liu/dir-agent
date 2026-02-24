package launcher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dir-agent/internal/config"
	"dir-agent/internal/terminal"
)

var (
	ErrPathNotAccessible = errors.New("path is not accessible")
	ErrToolNotConfigured = errors.New("tool is not configured")
	ErrToolNotFound      = errors.New("tool command not found")
)

type LaunchRequest struct {
	ToolName  string
	InputPath string
	ExtraArgs []string
	Config    config.Config
}

func MergeArgs(defaultArgs []string, extraArgs []string) []string {
	merged := make([]string, 0, len(defaultArgs)+len(extraArgs))
	merged = append(merged, defaultArgs...)
	merged = append(merged, extraArgs...)
	return merged
}

func ResolveTargetDir(path string, resolveFileToParent bool) (string, error) {
	cleanPath := filepath.Clean(path)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrPathNotAccessible, cleanPath)
	}

	if info.IsDir() {
		return cleanPath, nil
	}

	if !resolveFileToParent {
		return "", fmt.Errorf("%w: %s", ErrPathNotAccessible, cleanPath)
	}

	parent := filepath.Dir(cleanPath)
	if parent == "" || parent == "." {
		return "", fmt.Errorf("%w: %s", ErrPathNotAccessible, cleanPath)
	}

	return parent, nil
}

func Launch(req LaunchRequest) error {
	targetDir, err := ResolveTargetDir(req.InputPath, req.Config.Behavior.ResolveFileToParent)
	if err != nil {
		return err
	}

	commandName, defaultArgs, err := resolveTool(req.ToolName, req.Config)
	if err != nil {
		return err
	}

	resolvedCommandPath, err := terminal.FindExecutable(commandName)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrToolNotFound, commandName)
	}

	args := MergeArgs(defaultArgs, req.ExtraArgs)
	launchOpts := terminal.LaunchOptions{
		PreferredTerminal:        req.Config.Terminals.Preferred,
		OpenMode:                 req.Config.Behavior.OpenMode,
		WorkingDir:               targetDir,
		CommandPath:              resolvedCommandPath,
		Args:                     args,
		WindowsTerminalProfile:   req.Config.Terminals.WindowsTerminal.Profile,
		WindowsTerminalShell:     req.Config.Terminals.WindowsTerminal.Shell,
		WindowsTerminalCmderInit: req.Config.Terminals.WindowsTerminal.CmderInit,
	}

	return terminal.LaunchInTerminal(launchOpts)
}

func resolveTool(toolName string, cfg config.Config) (string, []string, error) {
	switch strings.ToLower(strings.TrimSpace(toolName)) {
	case "codex":
		if cfg.Tools.Codex.Command == "" {
			return "", nil, fmt.Errorf("%w: codex", ErrToolNotConfigured)
		}
		return cfg.Tools.Codex.Command, cfg.Tools.Codex.DefaultArgs, nil
	case "claude":
		if cfg.Tools.Claude.Command == "" {
			return "", nil, fmt.Errorf("%w: claude", ErrToolNotConfigured)
		}
		return cfg.Tools.Claude.Command, cfg.Tools.Claude.DefaultArgs, nil
	default:
		return "", nil, fmt.Errorf("%w: %s", ErrToolNotConfigured, toolName)
	}
}
