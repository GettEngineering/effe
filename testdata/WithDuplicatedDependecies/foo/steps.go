package main

type SuperRepository interface {
	Foo() error
}

func step1(repo SuperRepository) func() error {
	return func() error {
		return repo.Foo()
	}
}

func step2(repo SuperRepository) func() error {
	return func() error {
		return repo.Foo()
	}
}
