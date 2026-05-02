package main

import (
	"log/slog"
	"os"
	"runtime"
)

var cpus = runtime.NumCPU()

func main() {
	err := runWith("generator.fnl", "task.fnl")
	if err != nil {
		slog.Error("Runner failed", "error", err)
		os.Exit(1)
	}
}
