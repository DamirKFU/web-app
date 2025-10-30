package core

import "time"

type AbstractModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AbstractNameModel struct {
	AbstractModel
	Name string `gorm:"type:varchar(150);unique;not null" json:"name"`
}

func (a AbstractNameModel) String() string {
	return a.Name
}
