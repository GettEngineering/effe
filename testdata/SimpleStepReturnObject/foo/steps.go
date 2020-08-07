package main

type Foo struct {
	Name string
}

func step1() func() error {
	return func() error {
		return nil
	}
}

func step2() func() (Foo, error) {
	return func() (Foo, error) {
		return Foo{}, nil
	}
}
