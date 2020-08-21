// Code generated by Effe. DO NOT EDIT.

//+build !effeinject

package main

import (
	"github.com/pkg/errors"
)

func BuildComponent2(service BuildComponent2Service) BuildComponent2Func {
	return func() error {
		err := service.Step1()
		if err != nil {
			return errors.Wrap(err, "failure call Step1")
		}
		err = service.Step2()
		if err != nil {
			return errors.Wrap(err, "failure call Step2")
		}
		return nil
	}
}
func NewBuildComponent2Impl() *BuildComponent2Impl {
	return &BuildComponent2Impl{step1FieldFunc: step1(), step2FieldFunc: step2()}
}

type BuildComponent2Service interface {
	Step1() error
	Step2() error
}
type BuildComponent2Impl struct {
	step1FieldFunc func() error
	step2FieldFunc func() error
}
type BuildComponent2Func func() error

func (b *BuildComponent2Impl) Step1() error { return b.step1FieldFunc() }
func (b *BuildComponent2Impl) Step2() error { return b.step2FieldFunc() }
