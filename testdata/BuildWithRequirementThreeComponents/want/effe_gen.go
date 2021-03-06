// Code generated by Effe. DO NOT EDIT.

//+build !effeinject

package main

func D(service DService) DFunc {
	return func() error {
		err := service.Step2()
		if err != nil {
			return err
		}
		return nil
	}
}
func B(service BService) BFunc {
	return func() error {
		err := service.Step1()
		if err != nil {
			return err
		}
		return nil
	}
}
func C(service CService) CFunc {
	return func() error {
		err := service.B()
		if err != nil {
			return err
		}
		return nil
	}
}
func A(service AService) AFunc {
	return func() error {
		err := service.B()
		if err != nil {
			return err
		}
		err = service.C()
		if err != nil {
			return err
		}
		err = service.D()
		if err != nil {
			return err
		}
		return nil
	}
}
func NewDImpl() *DImpl {
	return &DImpl{step2FieldFunc: step2()}
}
func NewBImpl() *BImpl {
	return &BImpl{step1FieldFunc: step1()}
}
func NewCImpl(service BService) *CImpl {
	return &CImpl{BFieldFunc: B(service)}
}
func NewAImpl(service BService, service1 CService, service2 DService) *AImpl {
	return &AImpl{BFieldFunc: B(service), CFieldFunc: C(service1), DFieldFunc: D(service2)}
}

type DService interface {
	Step2() error
}
type DImpl struct {
	step2FieldFunc func() error
}
type DFunc func() error
type BService interface {
	Step1() error
}
type BImpl struct {
	step1FieldFunc func() error
}
type BFunc func() error
type CService interface {
	B() error
}
type CImpl struct {
	BFieldFunc func() error
}
type CFunc func() error
type AService interface {
	B() error
	C() error
	D() error
}
type AImpl struct {
	BFieldFunc func() error
	CFieldFunc func() error
	DFieldFunc func() error
}
type AFunc func() error

func (d *DImpl) Step2() error { return d.step2FieldFunc() }
func (b *BImpl) Step1() error { return b.step1FieldFunc() }
func (c *CImpl) B() error {
	return c.BFieldFunc()
}
func (a *AImpl) B() error {
	return a.BFieldFunc()
}
func (a *AImpl) C() error {
	return a.CFieldFunc()
}
func (a *AImpl) D() error {
	return a.DFieldFunc()
}
