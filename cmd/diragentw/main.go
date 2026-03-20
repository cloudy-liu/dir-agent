package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"dir-agent/internal/diagnostics"
	"dir-agent/internal/proc"
)

func main() {
	os.Exit(runMain(os.Args[1:], os.Stderr))
}

func runMain(args []string, stderr io.Writer) int {
	currentExe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(stderr, "[diragentw][ERROR] resolve current executable: %v\n", err)
		return 2
	}
	logf, closeFn := diagnostics.OpenLoggerForExecutable(currentExe)
	defer closeFn()
	logf("diragentw start args=%q", args)
	return runWithLogger(args, currentExe, exec.LookPath, startTarget, logf, stderr)
}

func run(
	args []string,
	currentExe string,
	lookPath func(string) (string, error),
	starter func(string, []string) error,
	stderr io.Writer,
) int {
	return runWithLogger(args, currentExe, lookPath, starter, func(string, ...any) {}, stderr)
}

func runWithLogger(
	args []string,
	currentExe string,
	lookPath func(string) (string, error),
	starter func(string, []string) error,
	logf diagnostics.LogFunc,
	stderr io.Writer,
) int {
	targetExe, err := resolveTargetExecutable(currentExe, lookPath)
	if err != nil {
		logf("resolve diragent executable failed: %v", err)
		fmt.Fprintf(stderr, "[diragentw][ERROR] resolve diragent executable: %v\n", err)
		return 2
	}
	logf("resolved diragent executable: %s", targetExe)
	if err := starter(targetExe, args); err != nil {
		logf("start diragent failed: %v", err)
		fmt.Fprintf(stderr, "[diragentw][ERROR] start diragent: %v\n", err)
		return 2
	}
	logf("started diragent target=%s args=%q", targetExe, args)
	return 0
}

func resolveTargetExecutable(currentExe string, lookPath func(string) (string, error)) (string, error) {
	if currentExe != "" {
		candidate := filepath.Join(filepath.Dir(currentExe), "diragent.exe")
		if fileExists(candidate) {
			return candidate, nil
		}
	}

	if path, err := lookPath("diragent.exe"); err == nil {
		return path, nil
	}
	if path, err := lookPath("diragent"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("diragent(.exe) not found near %q or in PATH", currentExe)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func startTarget(targetExe string, args []string) error {
	cmd := exec.Command(targetExe, args...)
	proc.ApplyNoWindow(cmd)
	return cmd.Start()
}
