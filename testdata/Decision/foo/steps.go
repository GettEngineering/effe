package main

import "fmt"

type a string

func failure() func(error) error {
	return func(err error) error {
		return err
	}
}

func step2() func(a) error {
	return func(v a) error {
		fmt.Println(v)
		return nil
	}
}

func step3() func() error {
	return func() error {
		return nil
	}
}

func step1() func() (a, error) {
	return func() (a, error) {
		return "a", nil
	}
}
