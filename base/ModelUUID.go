package base

import "github.com/Metadiv-Atomic-Engine/nanoid"

type ModelUUID struct {
	UUID string `json:"uuid" csv:"uuid"`
}

func (m *ModelUUID) GetUUID() string {
	return m.UUID
}

func (m *ModelUUID) InitUUID() {
	m.UUID = nanoid.NewSafe()
}
