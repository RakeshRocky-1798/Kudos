package service

type LeaderBoardResponse struct {
	LeaderBoardData []LeaderBoardData `json:"leaderBoardData"`
}

type LeaderBoardData struct {
	Rank          int      `json:"rank"`
	Count         int      `json:"totalCount"`
	UserMeta      UserData `json:"userMeta"`
	IsCurrentUser bool     `json:"isCurrentUser"`
}
