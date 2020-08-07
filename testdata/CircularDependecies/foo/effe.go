// +build effeinject

package main

import "github.com/GettEngineering/effe"

func A() error {
	effe.BuildFlow(
		effe.Step(B),
	)
	return nil
}
func B() error {
	effe.BuildFlow(
		effe.Step(step1),
		effe.Step(A),
	)
	return nil
}
