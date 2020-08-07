// +build effeinject

package main

import (
	"github.com/GettEngineering/effe"
	"github.com/GettEngineering/effe/testcustomization"
)

func C() error {
	effe.BuildFlow(
		effe.Step(step1),
		testcustomization.POST(
			"http://example.com",
		),
	)
	return nil
}
