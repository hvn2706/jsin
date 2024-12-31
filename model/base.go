package model

import "time"

type Base struct {
	ID        int64      `gorm:"primaryKey;column:id"`
	CreatedAt *time.Time `gorm:"created_at"`
	UpdatedAt *time.Time `gorm:"updated_at"`
}
