package main

import "fmt"

type a string

func failure() func(error) error {
	return func(err error) error {
		return err
	}
}

func step2() func(a) int {
	return func(v a) int {
		fmt.Println(v)
		return 0
	}
}

func step3() func() int {
	return func() int {
		return 0
	}
}

func step1() func() (a, error) {
	return func() (a, error) {
		return "a", nil
	}
}
