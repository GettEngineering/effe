// Code generated by Effe. DO NOT EDIT.

//+build !effeinject

package main

func A(service AService) AFunc {
	return func() error {
		err := service.Step1()
		if err != nil {
			return err
		}
		err = service.Step2()
		if err != nil {
			return err
		}
		return nil
	}
}
func NewAImpl(repo SuperRepository) *AImpl {
	return &AImpl{step1FieldFunc: step1(repo), step2FieldFunc: step2(repo)}
}

type AService interface {
	Step1() error
	Step2() error
}
type AImpl struct {
	step1FieldFunc func() error
	step2FieldFunc func() error
}
type AFunc func() error

func (a *AImpl) Step1() error { return a.step1FieldFunc() }
func (a *AImpl) Step2() error { return a.step2FieldFunc() }
