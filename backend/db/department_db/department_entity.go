package departmentDb

import "gorm.io/gorm"

type DepartmentEntity struct {
	gorm.Model
	Name string `gorm:"column:name"`
}

func (DepartmentEntity) TableName() string {
	return "departments"
}
