package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pelletier/go-toml/v2"
)

const (
	appFolderName = "dir-agent"
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
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(configDir, appFolderName, "config.toml")
	return configPath, nil
}

func DataPath() (string, error) {
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

func EnsureConfigFile() (string, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return "", err
	}

	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
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

	if err := os.WriteFile(configPath, defaultToml, 0o644); err != nil {
		return "", err
	}

	return configPath, nil
}
