package main

type Foo interface {
	Name() string
}

type foo struct{}

func (f foo) Name() string {
	return "foo"
}

func step1() func() error {
	return func() error {
		return nil
	}
}

func step2() func() (Foo, error) {
	return func() (Foo, error) {
		return foo{}, nil
	}
}
