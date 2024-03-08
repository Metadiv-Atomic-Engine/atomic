package base

type RequestIDPath struct {
	ID uint `uri:"id" validate:"required" json:"-"`
}
