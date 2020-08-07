package main

func step1() {
	return
}

func step2() func() error {
	return func() error {
		return nil
	}
}
