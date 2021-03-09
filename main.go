package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func init() {
	regModuleB()
	regModuleA()
	regVisorModule()
}

var a Module

func regModuleA() {
	init := func() error {
		time.Sleep(1 * time.Second)
		fmt.Println("initializing module a")
		return nil
	}
	a = MakeModule("a", init)
}

var b Module

func regModuleB() {
	init := func() error {
		time.Sleep(1 * time.Second)
		fmt.Printf("initializing module b\n")
		return nil
	}
	b = MakeModule("b", init, &a)
}

var visor Module

func regVisorModule() {
	init := func() error { return nil }
	visor = MakeModule("visor", init, &a, &b)
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	visor.InitConcurrent(ctx)
	err := a.Wait(ctx)
	if err != nil {
		log.Fatalf("Error init: %s", err)
	}

	// conf := Config("some config")
	// visor.InitSequential(conf)
}
