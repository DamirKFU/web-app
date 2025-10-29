package core

import "time"

type AbstractModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AbstractNameModel struct {
	AbstractModel
	Name string `gorm:"unique;not null;size:128"`
}
