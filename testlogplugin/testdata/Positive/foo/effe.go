// +build effeinject

package main

import (
	"github.com/GettEngineering/effe"
)

func BuildComponent2() error {
	effe.BuildFlow(
		effe.Step(step1),
		effe.Step(step2),
	)
	return nil
}
