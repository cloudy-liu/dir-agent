//go:build !windows

package proc

import "os/exec"

func ApplyNoWindow(_ *exec.Cmd) {}
