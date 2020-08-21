// Code generated by Effe. DO NOT EDIT.

//+build !effeinject

package main

func A(service AService) AFunc {
	return func(stringVal string) error {
		err := service.Step1(stringVal)
		if err != nil {
			return err
		}
		err = service.Step2(stringVal)
		if err != nil {
			return err
		}
		return nil
	}
}
func NewAImpl(repo NotificationRepository, repo1 UserRepository) *AImpl {
	return &AImpl{step1FieldFunc: step1(repo1), step2FieldFunc: step2(repo)}
}

type AService interface {
	Step1(id string) error
	Step2(userID string) error
}
type AImpl struct {
	step1FieldFunc func(id string) error
	step2FieldFunc func(userID string) error
}
type AFunc func(stringVal string) error

func (a *AImpl) Step1(id string) error     { return a.step1FieldFunc(id) }
func (a *AImpl) Step2(userID string) error { return a.step2FieldFunc(userID) }
