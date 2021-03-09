package main

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
		err := mod.InitSequential(conf)
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
