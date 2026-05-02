package main

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/spf13/pflag"
)

var (
	numWorkers    int
	generatorPath string
	taskPath      string
)

func main() {
	parseFlags()
	err := runWith(generatorPath, taskPath)
	if err != nil {
		slog.Error("Runner failed", "error", err)
		os.Exit(1)
	}
}

func parseFlags() {
	pflag.IntVarP(&numWorkers, "workers", "w", runtime.NumCPU(), "number of task runners")
	pflag.StringVarP(&generatorPath, "generator", "g", "generator.fnl", "path of generator spec file")
	pflag.StringVarP(&taskPath, "task", "t", "task.fnl", "path of task spec file")

	pflag.Parse()
}
