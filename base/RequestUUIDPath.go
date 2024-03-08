package base

type RequestUUIDPath struct {
	UUID string `uri:"uuid" validate:"required" json:"-"`
}
