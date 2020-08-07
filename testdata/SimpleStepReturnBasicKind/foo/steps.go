package main

import "net/http"

func step1() func() error {
	return func() error {
		return nil
	}
}

func step2() func() ([]string, error) {
	return func() ([]string, error) {
		return []string{}, nil
	}
}

func step3() func() ([]*string, error) {
	return func() ([]*string, error) {
		return []*string{}, nil
	}
}

func step4() func() ([]http.Request, error) {
	return func() ([]http.Request, error) {
		return []http.Request{}, nil
	}
}

type converter func(string) int

func step5() func() (converter, error) {
	return func() (converter, error) {
		return func(string) int {
			return 0
		}, nil
	}
}
