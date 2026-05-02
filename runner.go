package main

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	lualibs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

func runWith(generatorPath, taskPath string, numWorkers int) error {
	generatorModuleProto, err := compile(generatorPath)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	taskModuleProto, err := compile(taskPath)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	generatorTable, err := getGeneratorTable(generatorModuleProto)
	if err != nil {
		return fmt.Errorf("get generator table: %w", err)
	}

	var wg sync.WaitGroup
	messaging := make(chan lua.LValue)
	startWorkers(numWorkers, &wg, messaging, taskModuleProto)
	generate(generatorTable, messaging)
	wg.Wait()

	return err
}

func getGeneratorTable(generator *lua.FunctionProto) (*lua.LTable, error) {
	state := lua.NewState()
	defer state.Close()

	lualibs.Preload(state)
	generatorValue, err := runModuleValue(state, generator)
	if err != nil {
		return nil, fmt.Errorf("run: %w", err)
	}
	if generatorValue.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected a table of values but got %s", generatorValue.Type().String())
	}

	return generatorValue.(*lua.LTable), nil
}

func startWorkers(numWorkers int, wg *sync.WaitGroup, inbox <-chan lua.LValue, taskModuleProto *lua.FunctionProto) {
	for range numWorkers {
		wg.Go(func() {
			work(inbox, taskModuleProto)
		})
	}
}

func work(inbox <-chan lua.LValue, taskModuleProto *lua.FunctionProto) {
	state := lua.NewState()
	lualibs.Preload(state)
	for val := range inbox {
		_, err := runModuleValue(state, taskModuleProto, val)
		if err != nil {
			slog.Error("Task failed", "error", err)
		}
	}
	state.Close()
}

func generate(generatorSlice *lua.LTable, outbox chan<- lua.LValue) {
	generatorSlice.ForEach(func(key, val lua.LValue) {
		maybeSleep(sleepValue)
		outbox <- val
	})
	close(outbox)
}

func maybeSleep(sleep int) {
	if sleep == 0 {
		return
	}
	time.Sleep(time.Duration(sleep) * time.Second)
}
