package main

import (
	"github.com/pkg/errors"
)

func step1() func() error {
	return func() error {
		return nil
	}
}

func step2() func() error {
	return func() error {
		return nil
	}
}

func failure() func(error) error {
	return func(err error) error {
		return errors.Wrap(err, "failure call")
	}
}
