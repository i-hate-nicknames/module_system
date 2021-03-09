package main

// InitFN is a function that initializes module from a config
type InitFN func(conf string) error

// Module is a single system unit that represents a part of the system that must
// be initialized. Module can have dependencies, that should be initialized before
// module can start its own initialization
type Module struct {
	Name string
	init InitFN
	err  error
	done chan struct{}
	deps []*Module
}

// MakeModule returns a new module with given init function and dependencies
func MakeModule(name string, init InitFN, deps ...*Module) *Module {
	return &Module{
		Name: name,
		init: init,
		deps: deps,
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

}
