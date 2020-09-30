package entities

type Cmd interface {
	Type() string
}

type Foo struct {
	A string
}
