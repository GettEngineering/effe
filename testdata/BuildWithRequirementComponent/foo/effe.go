// +build effeinject

package main

import (
	"github.com/GettEngineering/effe"
)

func BuildComponent2() error {
	effe.BuildFlow(
		effe.Step(step2),
		effe.Step(BuildComponent3),
		effe.Step(BuildComponent5),
	)
	return nil
}

func BuildComponent5() error {
	effe.BuildFlow(
		effe.Step(step5),
	)
	return nil
}

func BuildComponent6() error {
	effe.BuildFlow(
		effe.Step(step6),
	)
	return nil
}

func BuildComponent3() error {
	effe.BuildFlow(
		effe.Step(step3),
	)
	return nil
}
func BuildComponent4() error {
	effe.BuildFlow(
		effe.Step(BuildComponent3),
		effe.Step(BuildComponent5),
	)
	return nil
}

func BuildComponent1() error {
	effe.BuildFlow(
		effe.Step(step1),
		effe.Wrap(effe.Before(before), effe.Success(success), effe.Failure(failure),
			effe.Step(BuildComponent6),
			effe.Step(BuildComponent6),
			effe.Step(BuildComponent6),
			effe.Step(BuildComponent6),
		),
	)
	return nil
}
