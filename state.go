package main

// type State struct {
// 	Name    string
// 	Require *State
// 	Modules []*Module
// }

// func (t *State) RegisterModule(m *Module) {
// 	t.Modules = append(t.Modules, m)
// }

// func (t *State) InitSequential(conf Config) error {
// 	if t.Require != nil {
// 		t.Require.InitSequential(conf)
// 	}
// 	for _, mod := range t.Modules {
// 		err := mod.InitSequential(conf)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (t *State) InitConcurrent(ctx context.Context, conf Config) error {
// 	if t.Require != nil {
// 		t.Require.InitConcurrent(ctx, conf)
// 	}
// 	for _, mod := range t.Modules {
// 		go mod.InitConcurrent(ctx, conf)
// 	}
// 	for _, mod := range t.Modules {
// 		<-mod.done
// 		if mod.err != nil {
// 			return mod.err
// 		}
// 	}
// 	return nil
// }
