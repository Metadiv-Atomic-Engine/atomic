package atomic

/*
Init a new module with this function.
*/
func (e *engine) NewModule(symbol string, dependencies ...*Module) *Module {
	return &Module{
		Symbol:       symbol,
		Dependencies: dependencies,
	}
}

/*
Register the module to installed modules.
Use this function in the Install function.
*/
func (e *engine) InstallModule(m *Module) {
	if _, ok := e.Modules[m.Symbol]; ok {
		panic("Module already exists: " + m.Symbol)
	}
	e.Modules[m.Symbol] = m
}
