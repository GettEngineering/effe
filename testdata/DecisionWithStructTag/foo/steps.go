package main

import "fmt"

type SuperStruct struct {
	Foo bool
}

func step2() func(SuperStruct) error {
	return func(s SuperStruct) error {
		fmt.Println(s.Foo)
		return nil
	}
}

func step3() func() error {
	return func() error {
		return nil
	}
}

func step1() func() (SuperStruct, error) {
	return func() (SuperStruct, error) {
		return SuperStruct{Foo: true}, nil
	}
}
