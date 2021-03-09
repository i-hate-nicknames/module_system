package main

import (
	"fmt"
	"time"
)

type Config string

func init() {
	regModuleB()
	regModuleA()
	regVisorModule()
}

var a Module

func regModuleA() {
	init := func(conf Config) error {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module a with conf %s\n", conf)
		return nil
	}
	a = MakeModule("a", init)
}

var b Module

func regModuleB() {
	init := func(conf Config) error {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module b with conf %s\n", conf)
		return nil
	}
	b = MakeModule("b", init, &a)
}

var visor Module

func regVisorModule() {
	init := func(conf Config) error { return nil }
	visor = MakeModule("visor", init, &a, &b)
}

func main() {
	// conf := Config("some config")
	// ctx := context.Background()
	// visor.InitConcurrent(ctx, conf)

	conf := Config("some config")
	visor.InitSequential(conf)
}
