package diagnostics

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogFunc func(string, ...any)

func OpenLoggerForExecutable(exePath string) (LogFunc, func()) {
	logPath := filepath.Join(filepath.Dir(exePath), "diragent.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return func(string, ...any) {}, func() {}
	}

	var mu sync.Mutex
	logf := func(format string, args ...any) {
		mu.Lock()
		defer mu.Unlock()
		line := fmt.Sprintf(format, args...)
		_, _ = fmt.Fprintf(file, "%s pid=%d %s\n", time.Now().Format(time.RFC3339), os.Getpid(), line)
	}
	closeFn := func() {
		mu.Lock()
		defer mu.Unlock()
		_ = file.Close()
	}
	return logf, closeFn
}
