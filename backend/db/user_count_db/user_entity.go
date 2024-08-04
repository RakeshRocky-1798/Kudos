package userCountDb

import "gorm.io/gorm"

type UserCountEntity struct {
	gorm.Model
	UserId        string `gorm:"column:user_id"`
	GivenCount    int    `gorm:"column:given_count"`
	ReceivedCount int    `gorm:"column:received_count"`
	Month         string `gorm:"column:month"`
	Week          string `gorm:"column:week"`
}

type UserMonthlyCount struct {
	UserId string `gorm:"column:user_id"`
	Count  int    `gorm:"column:count__"`
}

type UserGivenCount struct {
	UserId     string `gorm:"column:user_id"`
	GivenCount int    `gorm:"column:given_count"`
}

type LeaderBoardData struct {
	UserId string `gorm:"column:user_id"`
	Count  int    `gorm:"column:count__"`
	Rank   int    `gorm:"column:rank"`
}

func (UserCountEntity) TableName() string {
	return "user_count"
}
