// +build effeinject

package main

import "github.com/GettEngineering/effe"

func A() error {
	effe.BuildFlow(
		effe.Step(step1),
		effe.Wrap(effe.Failure(failureStep), effe.Success(successStep),
			effe.Step(step2),
			effe.Step(step3),
		),
	)
	return nil
}
