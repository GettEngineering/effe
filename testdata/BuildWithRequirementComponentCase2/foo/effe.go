// +build effeinject

package main

import "github.com/GettEngineering/effe"

func A() error {
	effe.BuildFlow(
		effe.Step(B),
		effe.Step(C),
	)
	return nil
}
func B() error {
	effe.BuildFlow(
		effe.Step(step1),
	)
	return nil
}
func C() error {
	effe.BuildFlow(
		effe.Step(B),
	)
	return nil
}
