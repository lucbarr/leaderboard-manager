package leaderboard

type Score struct {
	PlayerID      string `json:"playerID"`
	LeaderboardID string `json:"leaderboardID"`
	Score         int    `json:"score"`
	Rank          int    `json:"rank"`
}
