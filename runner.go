package main

import (
	"fmt"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

func runWith(generatorPath, taskPath string) error {
	generatorModuleProto, err := compile(generatorPath)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	taskModuleProto, err := compile(taskPath)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	generatorSlice, err := gatherGeneratorSlice(generatorModuleProto)
	if err != nil {
		return fmt.Errorf("gather generator slice: %w", err)
	}

	var wg sync.WaitGroup
	messaging := make(chan lua.LValue)
	startWorkers(numWorkers, &wg, messaging, taskModuleProto)
	generator(generatorSlice, messaging)
	wg.Wait()

	return err
}

func startWorkers(numWorkers int, wg *sync.WaitGroup, inbox <-chan lua.LValue, taskModuleProto *lua.FunctionProto) {
	for range numWorkers {
		wg.Go(func() {
			worker(inbox, taskModuleProto)
		})
	}
}

func worker(inbox <-chan lua.LValue, taskModuleProto *lua.FunctionProto) {
	state := lua.NewState()
	for val := range inbox {
		runFirstModuleFunction(state, taskModuleProto, val)
	}
	state.Close()
}

func generator(generatorSlice *lua.LTable, outbox chan<- lua.LValue) {
	generatorSlice.ForEach(func(key, val lua.LValue) {
		outbox <- val
	})
	close(outbox)
}

func gatherGeneratorSlice(generator *lua.FunctionProto) (*lua.LTable, error) {
	tmpState := lua.NewState()
	defer tmpState.Close()

	generatorValue, err := runFirstModuleFunction(tmpState, generator)
	if err != nil {
		return nil, fmt.Errorf("run: %w", err)
	}
	if generatorValue.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected a table of values but got %s", generatorValue.Type().String())
	}

	return generatorValue.(*lua.LTable), nil
}
