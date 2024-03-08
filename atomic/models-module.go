package atomic

type Module struct {
	Symbol       string    `json:"symbol"`       // Symbol is the unique identifier for the module.
	Dependencies []*Module `json:"dependencies"` // Dependencies is a list of modules that this module depends on.
}
