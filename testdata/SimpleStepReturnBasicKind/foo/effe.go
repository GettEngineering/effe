// +build effeinject

package main

import "github.com/GettEngineering/effe"

func A() error {
	effe.BuildFlow(
		effe.Step(step1),
		effe.Step(step2),
		effe.Step(step3),
		effe.Step(step4),
		effe.Step(step5),
	)
	return nil
}
