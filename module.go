package main

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// InitFN is a function that initializes module from a config
type InitFN func(conf Config) error

// Module is a single system unit that represents a part of the system that must
// be initialized. Module can have dependencies, that should be initialized before
// module can start its own initialization
type Module struct {
	Name    string
	init    InitFN
	err     error
	done    chan struct{}
	deps    []*Module
	mux     *sync.Mutex
	running bool
}

// MakeModule returns a new module with given init function and dependencies
func MakeModule(name string, init InitFN, deps ...*Module) Module {
	done := make(chan struct{}, 0)
	var mux sync.Mutex
	return Module{
		Name: name,
		init: init,
		deps: deps,
		done: done,
		mux:  &mux,
	}
}

func (m *Module) setRunning(val bool) bool {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.running == val {
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
func (m *Module) InitSequential(conf Config) error {
	for _, dep := range m.deps {
		if dep.isInitDone() {
			continue
		}
		err := dep.InitSequential(conf)
		if err != nil {
			return err
		}
	}
	err := m.init(conf)
	close(m.done)
	return err
}

func (m *Module) Wait(ctx context.Context) error {
	select {
	case <-m.done:
		if m.err != nil {
			return m.err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// InitConcurrent initializes all module dependencies recursively and concurrently.
// If module depends on modules a and b, this function will try to run init functions for a and b
// in each in a separate goroutine. It will block and wait on modules whose dependencies are not
// yet fully initialized themselves
// This function should be run in a separate goroutine
func (m *Module) InitConcurrent(ctx context.Context, conf Config) {
	// don't do anything if we already started
	ok := m.setRunning(true)
	if !ok {
		return
	}
	defer func() {
		log.Printf("mod %s: finishing init", m.Name)
		close(m.done)
		ok = m.setRunning(false)
		// this should never happen
		if !ok {
			panic(fmt.Sprintf("double initialization of module %s", m.Name))
		}
	}()
	// start init in every dependency
	for _, dep := range m.deps {
		log.Printf("mod %s: init dep %s", m.Name, dep.Name)
		if !m.isInitDone() {
			go dep.InitConcurrent(ctx, conf)
		}
	}

	// wait for every dependency to finish
	// collect error status for each, and set own error in case
	// any dependency errored
	// when cancelled return immediately
	for _, dep := range m.deps {
		log.Printf("mod %s: wait dep %s", m.Name, dep.Name)
		err := dep.Wait(ctx)
		if err != nil {
			m.err = err
		}
	}
	log.Printf("mod %s: init self", m.Name)
	// init the module itself
	err := m.init(conf)
	if err != nil {
		m.err = err
	}
}
