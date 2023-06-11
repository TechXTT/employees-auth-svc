package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
	Username string    `gorm:"not null, unique"`
	Email    string    `gorm:"not null, unique"`
	Password string    `gorm:"not null"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	a.ID = uuid

	return
}
