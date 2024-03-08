package atomic

type Error struct {
	UUID   string  `json:"uuid"`
	Module *Module `json:"module"`

	Eng string `json:"eng"`
	Zht string `json:"zht"`
	Zhs string `json:"zhs"`
}
