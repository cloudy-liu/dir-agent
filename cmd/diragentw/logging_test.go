package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dir-agent/internal/diagnostics"
)

func TestRunWritesLogWhenTargetResolutionFails(t *testing.T) {
	tempDir := t.TempDir()
	currentExe := filepath.Join(tempDir, "diragentw.exe")
	writeEmptyFile(t, currentExe)

	logf, closeFn := diagnostics.OpenLoggerForExecutable(currentExe)
	defer closeFn()

	starterCalled := false
	var stderr bytes.Buffer
	exitCode := runWithLogger([]string{"launch"}, currentExe,
		func(name string) (string, error) {
			return "", errors.New("not found")
		},
		func(name string, args []string) error {
			starterCalled = true
			return nil
		},
		logf,
		&stderr,
	)
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}
	if starterCalled {
		t.Fatalf("starter should not be called when target resolution fails")
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "diragent.log"))
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	if !strings.Contains(string(content), "resolve diragent executable") {
		t.Fatalf("expected diragentw failure to be logged, got %q", string(content))
	}
}
