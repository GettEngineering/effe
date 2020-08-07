package main

type stepFunc func() error

func step1() stepFunc {
	return func() error {
		return nil
	}
}
