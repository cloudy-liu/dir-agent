package launcher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeArgs(t *testing.T) {
	args := MergeArgs([]string{"--model", "gpt-5"}, []string{"--approval", "never"})
	if len(args) != 4 {
		t.Fatalf("expected merged args length 4, got %d", len(args))
	}
	if args[0] != "--model" || args[3] != "never" {
		t.Fatalf("unexpected merge result: %#v", args)
	}
}

func TestResolveTargetDirFromDirectory(t *testing.T) {
	tempDir := t.TempDir()
	dir, err := ResolveTargetDir(tempDir, true)
	if err != nil {
		t.Fatalf("resolve dir: %v", err)
	}
	if dir != tempDir {
		t.Fatalf("expected same dir, got %q", dir)
	}
}

func TestResolveTargetDirFromFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "a.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	dir, err := ResolveTargetDir(filePath, true)
	if err != nil {
		t.Fatalf("resolve file dir: %v", err)
	}
	if dir != tempDir {
		t.Fatalf("expected parent dir %q, got %q", tempDir, dir)
	}
}

func TestResolveTargetDirRejectsMissing(t *testing.T) {
	_, err := ResolveTargetDir(filepath.Join(t.TempDir(), "none"), true)
	if err == nil {
		t.Fatalf("expected error for missing path")
	}
}
