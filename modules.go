package main

import (
	"fmt"
	"time"
)

type Module struct {
	Name    string
	InitFN  func(conf string)
	Success bool
	Done    chan struct{}
	Deps    []*Module
}

// do I need it?
type State struct {
	Name    string
	Require *State
	Modules []*Module
}

func (t *State) RegisterModule(m *Module) {
	t.Modules = append(t.Modules, m)
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
	initfn := func(conf string) {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module a with conf %s\n", conf)
	}
	a = &Module{Name: "a", InitFN: initfn}
	visorState.RegisterModule(a)
}

var b *Module

func initModuleB() {
	initfn := func(conf string) {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module b with conf %s\n", conf)
	}
	b = &Module{Name: "b", InitFN: initfn}
	b.Deps = append(b.Deps, a)
	visorState.RegisterModule(b)
}

func InitSequential(t *State) {
	if t.Require != nil {
		InitSequential(t.Require)
	}
	conf := "seq conf"
	for _, mod := range t.Modules {
		fmt.Println("sraka")
		mod.InitFN(conf)
	}
}

func InitConcurrent(t *State) {
	if t.Require != nil {
		InitConcurrent(t.Require)
	}
}

func main() {
	InitSequential(visorState)
}
