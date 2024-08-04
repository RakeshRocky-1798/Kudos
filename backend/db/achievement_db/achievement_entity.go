package achiementDb

type AchievementEntity struct {
	ID              uint   `gorm:"column:id"`
	AchievementName string `gorm:"column:achievement_name"`
	DisplayName     string `gorm:"column:display_name"`
	Emoji           string `gorm:"column:emoji"`
	Points          string `gorm:"column:points"`
}

type AchievementName struct {
	ID              uint   `gorm:"column:id"`
	AchievementName string `gorm:"column:achievement_name"`
}

func (AchievementEntity) TableName() string {
	return "achievement"
}
