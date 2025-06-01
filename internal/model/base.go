package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	Id        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = uuid.NewString()
	b.CreatedAt = time.Now().UTC()
	b.UpdatedAt = time.Now().UTC()
	return
}

func (b *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now().UTC()
	return
}
