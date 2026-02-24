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
	Preferred string `toml:"preferred"`
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
		Tools: ToolsConfig{
			Codex: ToolConfig{
				Command:     "codex",
				DefaultArgs: []string{},
			},
			Claude: ToolConfig{
				Command:     "claude",
				DefaultArgs: []string{},
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
	if cfg.Tools.Claude.Command == "" {
		cfg.Tools.Claude.Command = "claude"
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

[tools.codex]
command = "codex"
default_args = []

[tools.claude]
command = "claude"
default_args = []

[behavior]
resolve_file_to_parent = true
open_mode = "tab_preferred"
`)

	if err := os.WriteFile(configPath, defaultToml, 0o644); err != nil {
		return "", err
	}

	return configPath, nil
}
