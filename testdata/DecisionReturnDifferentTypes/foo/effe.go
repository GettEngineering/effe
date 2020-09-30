// +build effeinject

package main

import "github.com/GettEngineering/effe"

func A() error {
	effe.BuildFlow(
		effe.Step(step1),
		effe.Decision(new(a), effe.Failure(failure),
			effe.Case("a", effe.Step(step2)),
			effe.Case("", effe.Step(step3)),
			effe.Case("b", effe.Step(step4)),
			effe.Case("c", effe.Step(step5)),
			effe.Case("d", effe.Step(step6)),
			effe.Case("e", effe.Step(step6)),
		),
	)
	return nil
}
