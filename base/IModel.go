package base

import "time"

type IModel interface {
	GetID() uint
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

func (m *Model) GetID() uint {
	return m.ID
}

func (m *Model) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Model) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
