package main

import (
	"fmt"
	"time"
)

type Module struct {
	Name    string
	InitFN  func(conf string) error
	Success bool
	Done    chan struct{}
	Deps    []*Module
}

type State struct {
	Name    string
	Require *State
	Modules []*Module
}

func (t *State) RegisterModule(m *Module) {
	t.Modules = append(t.Modules, m)
}

func (t *State) InitSequential() error {
	if t.Require != nil {
		t.Require.InitSequential()
	}
	conf := "seq conf"
	for _, mod := range t.Modules {
		fmt.Println("sraka")
		err := mod.InitFN(conf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *State) InitConcurrent() error {
	if t.Require != nil {
		t.Require.InitConcurrent()
	}
	return nil
}

var visorState *State

func init() {
	initStates()
	initModuleB()
	initModuleA()
}

func initStates() {
	visorState = &State{Name: "visor.State"}
}

var a *Module

func initModuleA() {
	initfn := func(conf string) error {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module a with conf %s\n", conf)
		return nil
	}
	a = &Module{Name: "a", InitFN: initfn}
	visorState.RegisterModule(a)
}

var b *Module

func initModuleB() {
	initfn := func(conf string) error {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module b with conf %s\n", conf)
		return nil
	}
	b = &Module{Name: "b", InitFN: initfn}
	b.Deps = append(b.Deps, a)
	visorState.RegisterModule(b)
}

func main() {
	visorState.InitSequential()
}
