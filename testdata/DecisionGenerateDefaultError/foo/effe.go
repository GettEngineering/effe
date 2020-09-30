// +build effeinject

package main

import "github.com/GettEngineering/effe"

func A() error {
	effe.BuildFlow(
		effe.Step(step1),
		effe.Decision(new(a), effe.Failure(failure),
			effe.Case("a", effe.Step(step2)),
			effe.Case("", effe.Step(step3)),
		),
	)
	return nil
}
