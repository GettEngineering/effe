package main

type StepFunc func() error

func step1() StepFunc {
	return func() error {
		return nil
	}
}
