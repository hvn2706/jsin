package model

import "time"

type Base struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
