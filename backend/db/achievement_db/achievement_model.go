package achiementDb

type AchievementRequest struct {
	AchievementName string `gorm:"column:achievement_name"`
	Points          int    `gorm:"column:points"`
}
