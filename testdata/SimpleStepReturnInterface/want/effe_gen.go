// Code generated by Effe. DO NOT EDIT.

//+build !effeinject

package main

import (
	"example.com/foo"
)

func A(service AService) AFunc {
	return func() (Foo, error) {
		err := service.Step1()
		if err != nil {
			return nil, err
		}
		FooVal, err := service.Step2()
		if err != nil {
			return FooVal, err
		}
		return FooVal, nil
	}
}
func NewAImpl() *AImpl {
	return &AImpl{step1FieldFunc: step1(), step2FieldFunc: step2()}
}

type AService interface {
	Step1() error
	Step2() (Foo, error)
}
type AImpl struct {
	step1FieldFunc func() error
	step2FieldFunc func() (Foo, error)
}
type AFunc func() (Foo, error)

func (a *AImpl) Step1() error        { return a.step1FieldFunc() }
func (a *AImpl) Step2() (Foo, error) { return a.step2FieldFunc() }