// +build effeinject

package main

import "github.com/GettEngineering/effe"

func A() error {
	effe.BuildFlow(
		effe.Step(step1),
		effe.Decision(nil),
	)
	return nil
}
