package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// expected command: app add 1 1 -> 2
	// expected command: app sub 1 1 -> 0

	slog.Info("get pid", slog.Int("pid", os.Getpid()))

	// caution: path injection, for POC is fine
	fileName := fmt.Sprintf("app-%s", os.Args[1])
	path, err := exec.LookPath(fileName)
	if err != nil {
		slog.Error("look path", slog.Any("err", err))
		os.Exit(1)
	}
	slog.Info("found", slog.String("path", path))
	slog.Info("pass", slog.Any("args", os.Args[2:]))

	// exec with same PID (caller as proxy terminated)
	if err := syscall.Exec(path, os.Args[2:], os.Environ()); err != nil {
		slog.Error("exec", slog.Any("err", err))
	}
}
