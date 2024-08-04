package managersDb

import "gorm.io/gorm"

type ManagerEntity struct {
	gorm.Model
	Email       string `gorm:"column:email"`
	SlackUserId string `gorm:"column:slack_user_id"`
	Name        string `gorm:"column:user_name"`
	RealName    string `gorm:"column:real_name"`
}
