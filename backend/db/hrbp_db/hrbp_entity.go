package hrbpDb

import "gorm.io/gorm"

type HrbpEntity struct {
	gorm.Model
	Email string `gorm:"column:email"`
	Name  string `gorm:"column:name"`
}

func (HrbpEntity) TableName() string {
	return "hrbp"
}
