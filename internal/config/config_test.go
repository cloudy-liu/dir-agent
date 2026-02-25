package config

import (
	"os"
	"path/filepath"
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
	if len(cfg.Tools.Codex.DefaultArgs) != 1 || cfg.Tools.Codex.DefaultArgs[0] != "--dangerously-bypass-approvals-and-sandbox" {
		t.Fatalf("expected codex default args to include highest-permission flag, got %#v", cfg.Tools.Codex.DefaultArgs)
	}
	if len(cfg.Tools.Claude.DefaultArgs) != 1 || cfg.Tools.Claude.DefaultArgs[0] != "--dangerously-skip-permissions" {
		t.Fatalf("expected claude default args to include highest-permission flag, got %#v", cfg.Tools.Claude.DefaultArgs)
	}
	if !cfg.Behavior.ResolveFileToParent {
		t.Fatalf("expected resolve_file_to_parent default true")
	}
	if cfg.Behavior.OpenMode != "tab_preferred" {
		t.Fatalf("expected open_mode tab_preferred, got %q", cfg.Behavior.OpenMode)
	}
	if cfg.Terminals.WindowsTerminal.Shell != "powershell" {
		t.Fatalf("expected windows terminal shell default powershell, got %q", cfg.Terminals.WindowsTerminal.Shell)
	}
	if cfg.Terminals.WindowsTerminal.CmderInit != "" {
		t.Fatalf("expected windows terminal cmder_init default empty, got %q", cfg.Terminals.WindowsTerminal.CmderInit)
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
[terminals.windows_terminal]
profile = "Cmder"
shell = "cmder"
cmder_init = "C:\\path\\to\\cmder\\vendor\\init.bat"

[tools.codex]
default_args = []
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
	if len(cfg.Tools.Codex.DefaultArgs) != 1 || cfg.Tools.Codex.DefaultArgs[0] != "--dangerously-bypass-approvals-and-sandbox" {
		t.Fatalf("expected empty codex default args to be backfilled with highest-permission default, got %#v", cfg.Tools.Codex.DefaultArgs)
	}
	if cfg.Tools.Claude.Command != "claude" {
		t.Fatalf("expected claude defaults to remain")
	}
	if len(cfg.Tools.Claude.DefaultArgs) != 1 || cfg.Tools.Claude.DefaultArgs[0] != "--dangerously-skip-permissions" {
		t.Fatalf("expected claude default args to remain highest-permission defaults, got %#v", cfg.Tools.Claude.DefaultArgs)
	}
	if cfg.Behavior.OpenMode != "tab_preferred" {
		t.Fatalf("expected default open_mode tab_preferred, got %q", cfg.Behavior.OpenMode)
	}
	if cfg.Terminals.WindowsTerminal.Profile != "Cmder" {
		t.Fatalf("expected windows terminal profile to load, got %q", cfg.Terminals.WindowsTerminal.Profile)
	}
	if cfg.Terminals.WindowsTerminal.Shell != "cmder" {
		t.Fatalf("expected windows terminal shell to load cmder, got %q", cfg.Terminals.WindowsTerminal.Shell)
	}
	if cfg.Terminals.WindowsTerminal.CmderInit != "C:\\path\\to\\cmder\\vendor\\init.bat" {
		t.Fatalf("expected windows terminal cmder_init to load, got %q", cfg.Terminals.WindowsTerminal.CmderInit)
	}
}

func TestPathsHelpers(t *testing.T) {
	t.Setenv("DIRAGENT_HOME", t.TempDir())

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
	if filepath.Base(configPath) != "config.toml" {
		t.Fatalf("expected config file name config.toml, got %q", configPath)
	}
	if filepath.Base(dataPath) != "data" {
		t.Fatalf("expected data dir name data, got %q", dataPath)
	}
	if filepath.Dir(configPath) != os.Getenv("DIRAGENT_HOME") {
		t.Fatalf("expected config path under DIRAGENT_HOME, got %q", configPath)
	}
	if filepath.Dir(dataPath) != os.Getenv("DIRAGENT_HOME") {
		t.Fatalf("expected data path under DIRAGENT_HOME, got %q", dataPath)
	}
}

func TestEnsureConfigFileCreatesDirAgentDefault(t *testing.T) {
	t.Setenv("DIRAGENT_HOME", t.TempDir())
	legacyRoot := t.TempDir()
	t.Setenv("APPDATA", legacyRoot)
	t.Setenv("XDG_CONFIG_HOME", legacyRoot)
	t.Setenv("HOME", legacyRoot)

	configPath, err := EnsureConfigFile()
	if err != nil {
		t.Fatalf("ensure config file failed: %v", err)
	}
	if filepath.Dir(configPath) != os.Getenv("DIRAGENT_HOME") {
		t.Fatalf("config path must be in DIRAGENT_HOME, got %q", configPath)
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
	if !strings.Contains(string(content), "default_args = [\"--dangerously-bypass-approvals-and-sandbox\"]") {
		t.Fatalf("expected codex highest-permission default args in config file")
	}
	if !strings.Contains(string(content), "default_args = [\"--dangerously-skip-permissions\"]") {
		t.Fatalf("expected claude highest-permission default args in config file")
	}
	if !strings.Contains(string(content), "[terminals.windows_terminal]") {
		t.Fatalf("expected windows terminal section in default config file")
	}
	if !strings.Contains(string(content), "shell = \"powershell\"") {
		t.Fatalf("expected powershell as windows terminal shell default")
	}
	if !strings.Contains(string(content), "cmder_init = \"\"") {
		t.Fatalf("expected empty windows terminal cmder_init in default config file")
	}
}

func TestEnsureConfigFileMigratesLegacyConfigToDiragentHome(t *testing.T) {
	homeDir := t.TempDir()
	legacyRoot := t.TempDir()
	t.Setenv("DIRAGENT_HOME", homeDir)
	t.Setenv("APPDATA", legacyRoot)
	t.Setenv("XDG_CONFIG_HOME", legacyRoot)
	t.Setenv("HOME", legacyRoot)

	legacyPath := filepath.Join(legacyRoot, "dir-agent", "config.toml")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("mkdir legacy config dir: %v", err)
	}

	const legacyContent = `[tools.codex]
command = "codex-custom"
default_args = ["--foo"]
`
	if err := os.WriteFile(legacyPath, []byte(legacyContent), 0o644); err != nil {
		t.Fatalf("write legacy config: %v", err)
	}

	configPath, err := EnsureConfigFile()
	if err != nil {
		t.Fatalf("ensure config file failed: %v", err)
	}

	expectedPath := filepath.Join(homeDir, "config.toml")
	if configPath != expectedPath {
		t.Fatalf("expected migrated config path %q, got %q", expectedPath, configPath)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read migrated config: %v", err)
	}
	if !strings.Contains(string(content), `command = "codex-custom"`) {
		t.Fatalf("expected migrated config to preserve legacy content, got: %s", string(content))
	}
}

func TestEnsureConfigFileFallsBackToLegacyWhenMigrationCannotWrite(t *testing.T) {
	brokenTarget := filepath.Join(t.TempDir(), "blocked-target")
	if err := os.WriteFile(brokenTarget, []byte("not-a-directory"), 0o644); err != nil {
		t.Fatalf("write blocking file: %v", err)
	}
	t.Setenv("DIRAGENT_HOME", brokenTarget)

	legacyRoot := t.TempDir()
	t.Setenv("APPDATA", legacyRoot)
	t.Setenv("XDG_CONFIG_HOME", legacyRoot)
	t.Setenv("HOME", legacyRoot)

	legacyPath := filepath.Join(legacyRoot, "dir-agent", "config.toml")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("mkdir legacy dir: %v", err)
	}
	if err := os.WriteFile(legacyPath, []byte("[behavior]\nopen_mode=\"new_window\"\n"), 0o644); err != nil {
		t.Fatalf("write legacy config: %v", err)
	}

	configPath, err := EnsureConfigFile()
	if err != nil {
		t.Fatalf("ensure config file failed: %v", err)
	}
	if configPath != legacyPath {
		t.Fatalf("expected fallback to legacy config path %q, got %q", legacyPath, configPath)
	}
}
