package base

import (
	"time"

	"gorm.io/gorm"
)

/*
Model is following the GORM Model convention.
*/
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id" csv:"id"`
	CreatedAt time.Time      `json:"created_at" csv:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" csv:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" csv:"-"`
}
