package diagnostics

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenLoggerForExecutableWritesNextToExecutable(t *testing.T) {
	exeDir := t.TempDir()
	exePath := filepath.Join(exeDir, "diragent.exe")
	if err := os.WriteFile(exePath, []byte(""), 0o644); err != nil {
		t.Fatalf("write exe: %v", err)
	}

	logf, closeFn := OpenLoggerForExecutable(exePath)
	defer closeFn()

	logf("launch failed: %s", "boom")

	content, err := os.ReadFile(filepath.Join(exeDir, "diragent.log"))
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	if !strings.Contains(string(content), "launch failed: boom") {
		t.Fatalf("expected log line to be written, got %q", string(content))
	}
}
