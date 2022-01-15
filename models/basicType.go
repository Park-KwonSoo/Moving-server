package models

import "time"

type basicType struct {
	ID        uint
	CreateAt  time.Time
	UpdatedAt time.Time
	DeletedAt uint
}
