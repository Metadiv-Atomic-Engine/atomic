package atomic

type Environment struct {
	Key    string  `json:"key"`
	Type   string  `json:"type"`
	Value  any     `json:"value"`
	Note   string  `json:"note"`
	Module *Module `json:"module"`
}
