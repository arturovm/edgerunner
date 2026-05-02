package main

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
	parse "github.com/yuin/gopher-lua/parse"
)

//go:embed embed/fennel.lua
var fennelSource string

const fennelTemplate = `
local file, err = io.open("%s", "r")
if not file then
	error(err)
end
local lua = require("fennel").install().compile(file)
return lua
`

func compile(filePath string) (*lua.FunctionProto, error) {
	luaSource, err := compileFennel(filePath)
	if err != nil {
		return nil, fmt.Errorf("compile fennel: %w", err)
	}

	sourceReader := bufio.NewReader(strings.NewReader(luaSource))
	functionProto, err := compileLua(sourceReader, filePath)
	if err != nil {
		return nil, fmt.Errorf("compile lua: %w", err)
	}

	return functionProto, nil
}

func compileFennel(filePath string) (string, error) {
	state := lua.NewState()
	defer state.Close()

	err := state.DoString(fennelSource)
	if err != nil {
		return "", fmt.Errorf("load fennel: %w", err)
	}
	fennelModule := state.Get(-1)
	state.Pop(1)
	pkg := state.GetGlobal("package")
	loaded := state.GetField(pkg, "loaded")
	state.SetField(loaded, "fennel", fennelModule)

	err = state.DoString(fmt.Sprintf(fennelTemplate, filePath))
	if err != nil {
		return "", fmt.Errorf("do string: %w", err)
	}

	lv := state.Get(-1)
	if lv.Type() != lua.LTString {
		return "", errors.New("unexpected value type")
	}

	return string(lv.(lua.LString)), nil
}

func compileLua(sourceReader *bufio.Reader, filePath string) (*lua.FunctionProto, error) {
	chunk, err := parse.Parse(sourceReader, filePath)
	if err != nil {
		return nil, err
	}

	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}

	return proto, nil
}

func runFirstModuleFunction(state *lua.LState, moduleProto *lua.FunctionProto, args ...lua.LValue) (lua.LValue, error) {
	module := state.NewFunctionFromProto(moduleProto)

	state.Push(module)
	state.PCall(0, lua.MultRet, nil)
	function := state.Get(-1)
	state.Pop(state.GetTop())

	state.Push(function)
	for i := range len(args) {
		state.Push(args[i])
	}
	state.PCall(len(args), lua.MultRet, nil)
	function = state.Get(-1)
	state.Pop(state.GetTop())

	return function, nil
}

func luaTableToSlice(table *lua.LTable) []lua.LValue {
	slice := make([]lua.LValue, table.Len(), table.Len())
	key := lua.LNil
	val := lua.LNil
	for i := 0; i < table.Len(); i++ {
		key, val = table.Next(key)
		slice[i] = val
	}
	return slice
}
