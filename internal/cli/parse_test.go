package cli

import "testing"

func TestParseLaunchArgsWithPassthrough(t *testing.T) {
	opts, err := ParseLaunchArgs([]string{"--tool", "codex", "--path", "C:/repo", "--", "--model", "gpt-5"})
	if err != nil {
		t.Fatalf("parse launch args: %v", err)
	}
	if opts.Tool != "codex" {
		t.Fatalf("expected tool codex, got %q", opts.Tool)
	}
	if opts.Path != "C:/repo" {
		t.Fatalf("expected path C:/repo, got %q", opts.Path)
	}
	if len(opts.ExtraArgs) != 2 || opts.ExtraArgs[1] != "gpt-5" {
		t.Fatalf("expected extra args to keep passthrough, got %#v", opts.ExtraArgs)
	}
}

func TestParseLaunchArgsRequiresToolAndPath(t *testing.T) {
	_, err := ParseLaunchArgs([]string{"--tool", "codex"})
	if err == nil {
		t.Fatalf("expected missing path error")
	}
}
