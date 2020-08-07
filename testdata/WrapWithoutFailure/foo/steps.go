package main

func step1() func() error {
	return func() error {
		return nil
	}
}

func beforeStep() func() error {
	return func() error {
		return nil
	}
}

func successStep() func() error {
	return func() error {
		return nil
	}
}

func step2() func() error {
	return func() error {
		return nil
	}
}

func step3() func() error {
	return func() error {
		return nil
	}
}
