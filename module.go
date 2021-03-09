package main

import "sync"

// InitFN is a function that initializes module from a config
type InitFN func(conf string) error

// Module is a single system unit that represents a part of the system that must
// be initialized. Module can have dependencies, that should be initialized before
// module can start its own initialization
type Module struct {
	Name    string
	init    InitFN
	err     error
	done    <-chan struct{}
	deps    []*Module
	mux     sync.RWMutex
	started bool
}

// MakeModule returns a new module with given init function and dependencies
func MakeModule(name string, init InitFN, deps ...*Module) *Module {
	done := make(chan struct{}, 0)
	return &Module{
		Name: name,
		init: init,
		deps: deps,
		done: done,
	}
}

// InitSequential initializes all module dependencies recursively and sequentially, one by one
// first to last and depth first
// If any of the underlying dependencies, or this module initialize with error, return that error
func (m *Module) InitSequential(conf string) error {
	for _, dep := range m.deps {
		err := dep.InitSequential(conf)
		if err != nil {
			return err
		}
	}
	return m.init(conf)
}

// InitConcurrent initializes all module dependencies recursively and concurrently.
// If module depends on modules a and b, this function will try to run init functions for a and b
// in each in a separate goroutine. It will block and wait on modules whose dependencies are not
// yet fully initialized themselves
// This function should be run in a separate goroutine
func (m *Module) InitConcurrent() {
	// if is in the process of initialization, exit

	// for every dependency:
	// if a dependency is initialized, skip it
	// if a dependency is not initialized, run a goroutine
	// that will try to initialize it

	// for every dependency:
	// wait for finishing initialization
	// check error field of every dependency, if any has it, set its own error
	// field and exit

	// run own initialization, if error during the process, set error field and exit
}
