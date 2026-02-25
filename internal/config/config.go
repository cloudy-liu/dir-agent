package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const (
	appFolderName = "dir-agent"
	configName    = "config.toml"
	dataDirName   = "data"
	homeEnvName   = "DIRAGENT_HOME"
)

type Config struct {
	Terminals TerminalsConfig `toml:"terminals"`
	Tools     ToolsConfig     `toml:"tools"`
	Behavior  BehaviorConfig  `toml:"behavior"`
}

type TerminalsConfig struct {
	Preferred       string                `toml:"preferred"`
	WindowsTerminal WindowsTerminalConfig `toml:"windows_terminal"`
}

type WindowsTerminalConfig struct {
	Profile   string `toml:"profile"`
	Shell     string `toml:"shell"`
	CmderInit string `toml:"cmder_init"`
}

type ToolsConfig struct {
	Codex  ToolConfig `toml:"codex"`
	Claude ToolConfig `toml:"claude"`
}

type ToolConfig struct {
	Command     string   `toml:"command"`
	DefaultArgs []string `toml:"default_args"`
}

type BehaviorConfig struct {
	ResolveFileToParent bool   `toml:"resolve_file_to_parent"`
	OpenMode            string `toml:"open_mode"`
}

func DefaultConfig() Config {
	return Config{
		Terminals: TerminalsConfig{
			WindowsTerminal: WindowsTerminalConfig{
				Shell: "powershell",
			},
		},
		Tools: ToolsConfig{
			Codex: ToolConfig{
				Command:     "codex",
				DefaultArgs: []string{"--dangerously-bypass-approvals-and-sandbox"},
			},
			Claude: ToolConfig{
				Command:     "claude",
				DefaultArgs: []string{"--dangerously-skip-permissions"},
			},
		},
		Behavior: BehaviorConfig{
			ResolveFileToParent: true,
			OpenMode:            "tab_preferred",
		},
	}
}

func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()

	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}

	if err := toml.Unmarshal(content, &cfg); err != nil {
		return Config{}, err
	}

	applyDefaults(&cfg)
	return cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Tools.Codex.Command == "" {
		cfg.Tools.Codex.Command = "codex"
	}
	if len(cfg.Tools.Codex.DefaultArgs) == 0 {
		cfg.Tools.Codex.DefaultArgs = []string{"--dangerously-bypass-approvals-and-sandbox"}
	}
	if cfg.Tools.Claude.Command == "" {
		cfg.Tools.Claude.Command = "claude"
	}
	if len(cfg.Tools.Claude.DefaultArgs) == 0 {
		cfg.Tools.Claude.DefaultArgs = []string{"--dangerously-skip-permissions"}
	}
	if cfg.Terminals.WindowsTerminal.Shell == "" {
		cfg.Terminals.WindowsTerminal.Shell = "powershell"
	}
	if cfg.Behavior.OpenMode == "" {
		cfg.Behavior.OpenMode = "tab_preferred"
	}
}

func ConfigPath() (string, error) {
	homeDir, err := appHomeDir()
	if err == nil {
		return filepath.Join(homeDir, configName), nil
	}

	return LegacyConfigPath()
}

func DataPath() (string, error) {
	homeDir, err := appHomeDir()
	if err == nil {
		return filepath.Join(homeDir, dataDirName), nil
	}

	return LegacyDataPath()
}

func LegacyConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, appFolderName, configName), nil
}

func LegacyDataPath() (string, error) {
	if runtime.GOOS == "windows" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(configDir, appFolderName), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".local", "share", appFolderName), nil
}

func ActiveConfigPath() (string, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return "", err
	}
	if fileExists(configPath) {
		return configPath, nil
	}

	legacyPath, err := LegacyConfigPath()
	if err != nil {
		return configPath, nil
	}
	if fileExists(legacyPath) {
		return legacyPath, nil
	}

	return configPath, nil
}

func EnsureConfigFile() (string, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return "", err
	}

	if fileExists(configPath) {
		return configPath, nil
	}

	legacyPath, legacyErr := LegacyConfigPath()
	if legacyErr == nil && !samePath(configPath, legacyPath) && fileExists(legacyPath) {
		if migrateErr := migrateLegacyConfig(legacyPath, configPath); migrateErr == nil {
			return configPath, nil
		}
		fmt.Fprintf(os.Stderr, "[diragent][WARN] migrate legacy config from %s to %s failed; continue using legacy config\n", legacyPath, configPath)
		return legacyPath, nil
	}

	if mkErr := os.MkdirAll(filepath.Dir(configPath), 0o755); mkErr == nil {
		if writeErr := writeDefaultConfig(configPath); writeErr == nil {
			return configPath, nil
		}
	}

	if legacyErr != nil {
		return "", fmt.Errorf("cannot create config at %s and no legacy path available: %w", configPath, legacyErr)
	}
	if samePath(configPath, legacyPath) {
		return "", fmt.Errorf("cannot create config at %s", configPath)
	}
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		return "", err
	}
	if err := writeDefaultConfig(legacyPath); err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stderr, "[diragent][WARN] cannot create config in install directory; fallback to legacy config path: %s\n", legacyPath)
	return legacyPath, nil
}

func appHomeDir() (string, error) {
	homeOverride := strings.TrimSpace(os.Getenv(homeEnvName))
	if homeOverride != "" {
		return homeOverride, nil
	}

	exePath, err := os.Executable()
	if err == nil {
		if resolvedPath, resolveErr := filepath.EvalSymlinks(exePath); resolveErr == nil {
			exePath = resolvedPath
		}
		return filepath.Dir(exePath), nil
	}

	return "", err
}

func migrateLegacyConfig(source string, target string) error {
	content, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	return os.WriteFile(target, content, 0o644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func samePath(first string, second string) bool {
	first = filepath.Clean(first)
	second = filepath.Clean(second)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(first, second)
	}
	return first == second
}

func writeDefaultConfig(path string) error {
	if fileExists(path) {
		return nil
	}

	defaultToml := []byte(`[terminals]
preferred = ""

[terminals.windows_terminal]
profile = ""
shell = "powershell"
cmder_init = ""

[tools.codex]
command = "codex"
default_args = ["--dangerously-bypass-approvals-and-sandbox"]

[tools.claude]
command = "claude"
default_args = ["--dangerously-skip-permissions"]

[behavior]
resolve_file_to_parent = true
open_mode = "tab_preferred"
`)

	return os.WriteFile(path, defaultToml, 0o644)
}
