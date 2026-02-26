package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestResolveTargetExecutablePrefersSiblingDiragentExe(t *testing.T) {
	tempDir := t.TempDir()
	currentExe := filepath.Join(tempDir, "diragentw.exe")
	targetExe := filepath.Join(tempDir, "diragent.exe")
	writeEmptyFile(t, currentExe)
	writeEmptyFile(t, targetExe)

	lookPathCalled := false
	resolved, err := resolveTargetExecutable(currentExe, func(name string) (string, error) {
		lookPathCalled = true
		return filepath.Join(tempDir, "from-path", name), nil
	})
	if err != nil {
		t.Fatalf("resolve target executable: %v", err)
	}
	if resolved != targetExe {
		t.Fatalf("expected sibling target %q, got %q", targetExe, resolved)
	}
	if lookPathCalled {
		t.Fatalf("lookPath should not be called when sibling diragent.exe exists")
	}
}

func TestResolveTargetExecutableFallsBackToPath(t *testing.T) {
	tempDir := t.TempDir()
	currentExe := filepath.Join(tempDir, "diragentw.exe")
	writeEmptyFile(t, currentExe)

	expected := `C:\tools\diragent.exe`
	resolved, err := resolveTargetExecutable(currentExe, func(name string) (string, error) {
		if name == "diragent.exe" {
			return expected, nil
		}
		return "", errors.New("not found")
	})
	if err != nil {
		t.Fatalf("resolve target executable: %v", err)
	}
	if resolved != expected {
		t.Fatalf("expected PATH fallback %q, got %q", expected, resolved)
	}
}

func TestRunReturnsErrorWhenTargetMissing(t *testing.T) {
	tempDir := t.TempDir()
	currentExe := filepath.Join(tempDir, "diragentw.exe")
	writeEmptyFile(t, currentExe)

	starterCalled := false
	var stderr bytes.Buffer
	exitCode := run([]string{"launch"}, currentExe,
		func(name string) (string, error) {
			return "", errors.New("not found")
		},
		func(name string, args []string) error {
			starterCalled = true
			return nil
		},
		&stderr,
	)
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}
	if starterCalled {
		t.Fatalf("starter should not be called when target resolution fails")
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected error output when target is missing")
	}
}

func TestRunForwardsArgsToStarter(t *testing.T) {
	tempDir := t.TempDir()
	currentExe := filepath.Join(tempDir, "diragentw.exe")
	targetExe := filepath.Join(tempDir, "diragent.exe")
	writeEmptyFile(t, currentExe)
	writeEmptyFile(t, targetExe)

	var gotName string
	var gotArgs []string
	var stderr bytes.Buffer
	inputArgs := []string{"launch", "--tool", "codex", "--path", `C:\work\repo`}
	exitCode := run(inputArgs, currentExe,
		func(name string) (string, error) {
			return "", errors.New("lookPath should not be called when sibling exists")
		},
		func(name string, args []string) error {
			gotName = name
			gotArgs = append([]string(nil), args...)
			return nil
		},
		&stderr,
	)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q", exitCode, stderr.String())
	}
	if gotName != targetExe {
		t.Fatalf("expected starter target %q, got %q", targetExe, gotName)
	}
	if !reflect.DeepEqual(gotArgs, inputArgs) {
		t.Fatalf("expected args %#v, got %#v", inputArgs, gotArgs)
	}
}

func writeEmptyFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("write empty file %s: %v", path, err)
	}
}
