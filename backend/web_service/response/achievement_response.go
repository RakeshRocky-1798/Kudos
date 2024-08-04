package service

import (
	achiementDb "kleos/db/achievement_db"
)

type AchievementResponse struct {
	ID              uint   `json:"id"`
	AchievementName string `json:"achievementName"`
	Points          string `json:"points"`
}

func ConvertToAchievementResponse(achievementEntity achiementDb.AchievementEntity) AchievementResponse {
	return AchievementResponse{
		ID:              achievementEntity.ID,
		AchievementName: achievementEntity.AchievementName,
		Points:          achievementEntity.Points,
	}
}
