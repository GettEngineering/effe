package main

import (
	"net/http"
)

func step1() func() error {
	return func() error {
		return nil
	}
}

func step2() func() (http.Request, error) {
	return func() (http.Request, error) {
		return http.Request{}, nil
	}
}
