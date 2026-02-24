package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Tools.Codex.Command != "codex" {
		t.Fatalf("expected codex command, got %q", cfg.Tools.Codex.Command)
	}
	if cfg.Tools.Claude.Command != "claude" {
		t.Fatalf("expected claude command, got %q", cfg.Tools.Claude.Command)
	}
	if !cfg.Behavior.ResolveFileToParent {
		t.Fatalf("expected resolve_file_to_parent default true")
	}
	if cfg.Behavior.OpenMode != "tab_preferred" {
		t.Fatalf("expected open_mode tab_preferred, got %q", cfg.Behavior.OpenMode)
	}
}

func TestLoadConfigMissingReturnsDefaults(t *testing.T) {
	tempDir := t.TempDir()
	cfgPath := filepath.Join(tempDir, "missing.toml")

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("expected no error for missing config, got %v", err)
	}
	if cfg.Tools.Codex.Command != "codex" {
		t.Fatalf("expected defaults to apply for missing config")
	}
}

func TestLoadConfigMergesDefaults(t *testing.T) {
	tempDir := t.TempDir()
	cfgPath := filepath.Join(tempDir, "config.toml")
	err := os.WriteFile(cfgPath, []byte(`
[tools.codex]
default_args = ["--model", "gpt-5"]
`), 0o644)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Tools.Codex.Command != "codex" {
		t.Fatalf("expected missing command to default to codex, got %q", cfg.Tools.Codex.Command)
	}
	if len(cfg.Tools.Codex.DefaultArgs) != 2 {
		t.Fatalf("expected codex default args to load")
	}
	if cfg.Tools.Claude.Command != "claude" {
		t.Fatalf("expected claude defaults to remain")
	}
	if cfg.Behavior.OpenMode != "tab_preferred" {
		t.Fatalf("expected default open_mode tab_preferred, got %q", cfg.Behavior.OpenMode)
	}
}

func TestPathsHelpers(t *testing.T) {
	configPath, err := ConfigPath()
	if err != nil {
		t.Fatalf("config path err: %v", err)
	}
	dataPath, err := DataPath()
	if err != nil {
		t.Fatalf("data path err: %v", err)
	}
	if configPath == "" || dataPath == "" {
		t.Fatalf("paths must not be empty")
	}
	if filepath.Ext(configPath) != ".toml" {
		t.Fatalf("expected config path to end with .toml")
	}
	if !strings.Contains(strings.ToLower(configPath), "dir-agent") {
		t.Fatalf("expected config path to include dir-agent, got %q", configPath)
	}
	if !strings.Contains(strings.ToLower(dataPath), "dir-agent") {
		t.Fatalf("expected data path to include dir-agent, got %q", dataPath)
	}
	if runtime.GOOS == "windows" && !strings.Contains(strings.ToLower(configPath), "appdata") {
		t.Fatalf("expected windows config path under AppData, got %q", configPath)
	}
}

func TestEnsureConfigFileCreatesDirAgentDefault(t *testing.T) {
	tempDir := t.TempDir()
	switch runtime.GOOS {
	case "windows":
		t.Setenv("APPDATA", tempDir)
	case "darwin":
		t.Setenv("HOME", tempDir)
	default:
		t.Setenv("XDG_CONFIG_HOME", tempDir)
		t.Setenv("HOME", tempDir)
	}

	configPath, err := EnsureConfigFile()
	if err != nil {
		t.Fatalf("ensure config file failed: %v", err)
	}
	if !strings.Contains(strings.ToLower(configPath), "dir-agent") {
		t.Fatalf("config path must contain dir-agent, got %q", configPath)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(content), "command = \"codex\"") {
		t.Fatalf("expected diragent default config content")
	}
	if !strings.Contains(string(content), "command = \"claude\"") {
		t.Fatalf("expected default claude command")
	}
}
