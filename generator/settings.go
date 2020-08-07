package generator

// Settings for Generator
type Settings interface {
	FlowFuncPostfix() string
	LocalInterfaceVarname() string
	ImplFieldPostfix() string
	NewImplFuncPrefix() string
	ImplPostfix() string
	InterfaceNamePostfix() string
}

type settings struct {
	flowFuncPostfix       string
	localInterfaceVarname string
	implFieldPostfix      string
	newImplFuncPrefix     string
	implPostfix           string
	interfaceNamePostfix  string
}

// Returns a postfix for generated business flow function
func (s settings) FlowFuncPostfix() string {
	return s.flowFuncPostfix
}

// Returns a service object name which uses in business flow function
func (s settings) LocalInterfaceVarname() string {
	return s.localInterfaceVarname
}

// Returns a postifx which Effe adds to generated field
func (s settings) ImplFieldPostfix() string {
	return s.implFieldPostfix
}

// Returns a prefix which Effe adds to generated business flow function
func (s settings) NewImplFuncPrefix() string {
	return s.newImplFuncPrefix
}

// Returns a postfix which Effe adds to generated type for business flow function
func (s settings) ImplPostfix() string {
	return s.implPostfix
}

// Returns a postfix which Effe adds to generated business flow function
func (s settings) InterfaceNamePostfix() string {
	return s.interfaceNamePostfix
}

// Default values for settings
func DefaultSettigs() Settings {
	return settings{
		flowFuncPostfix:       "Func",
		localInterfaceVarname: "service",
		implFieldPostfix:      "FieldFunc",
		newImplFuncPrefix:     "New",
		implPostfix:           "Impl",
		interfaceNamePostfix:  "Service",
	}
}
