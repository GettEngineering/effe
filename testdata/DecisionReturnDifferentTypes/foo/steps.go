package main

import (
	"fmt"

	"example.com/foo/entities"
)

type a string

func failure() func(error) error {
	return func(err error) error {
		return err
	}
}

func step2() func(a) (entities.Cmd, []int, [3]int) {
	return func(v a) (entities.Cmd, []int, [3]int) {
		fmt.Println(v)
		return nil, []int{}, [3]int{}
	}
}

func step3() func() (entities.Cmd, *entities.Foo) {
	return func() (entities.Cmd, *entities.Foo) {
		return nil, nil
	}
}

func step4() func() entities.Foo {
	return func() entities.Foo {
		return entities.Foo{}
	}
}

func step5() func() int {
	return func() int {
		return 0
	}
}

func step6() func() string {
	return func() string {
		return ""
	}
}

func step7() func() bool {
	return func() bool {
		return false
	}
}

func step1() func(entities.Cmd) (a, error) {
	return func(c entities.Cmd) (a, error) {
		return "a", nil
	}
}
