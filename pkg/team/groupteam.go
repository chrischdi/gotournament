package team

type TeamGroupStats struct {
	GroupUUID    string
	GoalsFor     int
	GoalsAgainst int
	Points       int
	PlayedGames  int
}

func NewTeamGroupStats() TeamGroupStats {
	g := TeamGroupStats{
		GroupUUID:    "",
		GoalsFor:     0,
		GoalsAgainst: 0,
		Points:       0,
		PlayedGames:  0,
	}
	return g
}

func (g TeamGroupStats) GoalDifference() int {
	return g.GoalsFor - g.GoalsAgainst
}
