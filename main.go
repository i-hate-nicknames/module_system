package main

import (
	"fmt"
	"time"
)

var visorState *State

func init() {
	initStates()

	// order is fixed here, module registration of a dependent module should be
	// placed below modules it depends upon
	// this can be eliminated by passing around pointers
	regModuleA()
	regModuleB()
}

func initStates() {
	visorState = &State{Name: "visor.State"}
}

var a *Module

func regModuleA() {
	init := func(conf string) error {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module a with conf %s\n", conf)
		return nil
	}
	a = MakeModule("a", init)
	visorState.RegisterModule(a)
}

var b *Module

func regModuleB() {
	init := func(conf string) error {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module b with conf %s\n", conf)
		return nil
	}
	b = MakeModule("b", init, a)
	visorState.RegisterModule(b)
}

func main() {
	visorState.InitSequential()
}
