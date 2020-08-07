// +build effeinject

package main

import "github.com/GettEngineering/effe"

func A() error {
	effe.BuildFlow(
		effe.Step(step1),
		effe.Decision(SuperStruct{}.Foo,
			effe.Case(true, effe.Step(step2)),
			effe.Case(false, effe.Step(step3)),
		),
	)
	return nil
}
