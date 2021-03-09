package main

import (
	"context"
	"fmt"
	"sync"
)

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
	mux     sync.Mutex
	running bool
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

func (m *Module) setRunning(val bool) bool {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.running != val {
		return false
	}
	m.running = val
	return true
}

func (m *Module) isInitDone() bool {
	select {
	case <-m.done:
		return true
	default:
		return false
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
func (m *Module) InitConcurrent(ctx context.Context, conf string) {
	// don't do anything if we already started
	ok := m.setRunning(true)
	if !ok {
		return
	}
	// start init in every dependency
	for _, dep := range m.deps {
		if !m.isInitDone() {
			go dep.InitConcurrent(ctx, conf)
		}
	}

	// wait for every dependency to finish
	// collect error status for each, and set own error in case
	// any dependency errored
	// when cancelled return immediately
	for _, dep := range m.deps {
		select {
		case <-dep.done:
			if dep.err != nil {
				m.err = dep.err
				return
			}
		case <-ctx.Done():
			m.err = context.Canceled
			return
		}
	}
	// init the module itself
	err := m.init(conf)
	if err != nil {
		m.err = err
	}
	ok = m.setRunning(false)
	// this should never happen
	if !ok {
		panic(fmt.Sprintf("double initialization of module %s", m.Name))
	}
}
