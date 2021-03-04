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
type Target struct {
	Name    string
	Require *Target
	Modules []*Module
}

func (t *Target) RegisterModule(m *Module) {
	t.Modules = append(t.Modules, m)
}

var visorTarget *Target

func init() {
	initTargets()
	initModuleB()
	initModuleA()
}

func initTargets() {
	visorTarget = &Target{Name: "visor.target"}
}

var a *Module

func initModuleA() {
	initfn := func(conf string) {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module a with conf %s\n", conf)
	}
	a = &Module{Name: "a", InitFN: initfn}
	visorTarget.RegisterModule(a)
}

var b *Module

func initModuleB() {
	initfn := func(conf string) {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module b with conf %s\n", conf)
	}
	b = &Module{Name: "b", InitFN: initfn}
	b.Deps = append(b.Deps, a)
	visorTarget.RegisterModule(b)
}

func InitSequential(t *Target) {
	if t.Require != nil {
		InitSequential(t.Require)
	}
	conf := "seq conf"
	for _, mod := range t.Modules {
		fmt.Println("sraka")
		mod.InitFN(conf)
	}
}

func main() {
	InitSequential(visorTarget)
}
