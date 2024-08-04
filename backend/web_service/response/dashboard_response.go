package service

type DashboardResponse struct {
	MyData              UserData              `json:"myData"`
	KleosMetrics        KleosMetrics          `json:"kleosMetrics"`
	AchievementDropDown []Options             `json:"achievementDropDown"`
	TotalAchievement    []AchievementOptions  `json:"totalAchievement"`
	RecentRecognition   []RecentKleosResponse `json:"recentRecognition"`
}

type KleosMetrics struct {
	GivenCount    string `json:"givenCount"`
	ReceivedCount string `json:"receivedCount"`
}

type Options struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type AchievementOptionOne struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	AEmoji string `json:"aEmoji"`
}

type AchievementOptions struct {
	AchievementName string `json:"achievementName"`
	Count           int    `json:"count"`
	Emoji           string `json:"emoji"`
}

type UserData struct {
	Email      string `json:"email"`
	UserName   string `json:"userName"`
	ProfileUrl string `json:"profileUrl"`
}

type AdvancedAchievementResponse struct {
	ALabel string   `json:"aLabel"`
	AEmoji string   `json:"aEmoji"`
	AFrom  UserData `json:"aFrom"`
}

type RecentKleosResponse struct {
	Message     string                      `json:"message"`
	Achievement AdvancedAchievementResponse `json:"achievement"`
}
