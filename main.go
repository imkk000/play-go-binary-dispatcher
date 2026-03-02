package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	rootCmd := buildRootCommand()
	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("command run", slog.Any("err", err))
		os.Exit(1)
	}
}
