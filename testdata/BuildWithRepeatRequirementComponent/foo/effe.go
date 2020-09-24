// +build effeinject

package main

import (
	"github.com/GettEngineering/effe"
)

func C() error {
	effe.BuildFlow(
		effe.Step(step2),
		effe.Wrap(effe.Before(step1),
			effe.Step(B),
		),
		effe.Step(step1),
		effe.Wrap(effe.Before(step1),
			effe.Step(B),
		),
	)
	return nil
}

func B() error {
	effe.BuildFlow(
		effe.Step(A),
	)
	return nil
}

func A() error {
	effe.BuildFlow(
		effe.Step(step1),
	)
	return nil
}
