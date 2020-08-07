package main

type stepFunc func() error

func step1() stepFunc {
	return func() error {
		return nil
	}
}

func step2() stepFunc {
	return func() error {
		return nil
	}
}
func step3() stepFunc {
	return func() error {
		return nil
	}
}

func step5() stepFunc {
	return func() error {
		return nil
	}
}

func step6() stepFunc {
	return func() error {
		return nil
	}
}

func success() stepFunc {
	return func() error {
		return nil
	}
}

func before() stepFunc {
	return func() error {
		return nil
	}
}

func failure() func(err error) error {
	return func(err error) error {
		return err
	}
}
