package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func init() {
	regModuleB()
	regModuleC()
	regModuleA()
	regVisorModule()
}

var a Module

func regModuleA() {
	init := func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("initializing module a")
		return nil
	}
	a = MakeModule("a", init)
}

var b Module

func regModuleB() {
	init := func() error {
		time.Sleep(5 * time.Second)
		fmt.Printf("initializing module b\n")
		return nil
	}
	b = MakeModule("b", init)
}

var c Module

func regModuleC() {
	init := func() error {
		time.Sleep(5 * time.Second)
		fmt.Printf("initializing module c\n")
		return nil
	}
	c = MakeModule("c", init)
}

var visor Module

func regVisorModule() {
	init := func() error { return nil }
	visor = MakeModule("visor", init, &a, &b, &c)
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()
	visor.InitConcurrent(ctx)
	if visor.err != nil {
		log.Fatalf("Error init: %s", visor.err)
	}

	// conf := Config("some config")
	// visor.InitSequential(conf)
}
