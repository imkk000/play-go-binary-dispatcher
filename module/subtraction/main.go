package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

func main() {
	slog.Info("load module", slog.Any("args", os.Args))

	a, err := strconv.Atoi(os.Args[0])
	if err != nil {
		slog.Error("convert", slog.String("value", os.Args[0]))
		os.Exit(1)
	}
	b, err := strconv.Atoi(os.Args[1])
	if err != nil {
		slog.Error("convert", slog.String("value", os.Args[1]))
		os.Exit(1)
	}

	fmt.Printf("%d - %d = %d\n", a, b, a-b)
}
