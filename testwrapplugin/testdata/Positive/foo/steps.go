package main

func step1() func() error {
	return func() error {
		return nil
	}
}

func step2() func() error {
	return func() error {
		return nil
	}
}
