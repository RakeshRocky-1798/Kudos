package service

type FilteredKleosResponse struct {
	Id          int               `json:"id"`
	Message     string            `json:"message"`
	Achievement AchievementOption `json:"achievementData"`
}

type AchievementOption struct {
	AType      string   `json:"aType"`
	AEmoji     string   `json:"aEmoji"`
	ACreatedAt string   `json:"aCreatedAt"`
	User       UserData `json:"user"`
}
