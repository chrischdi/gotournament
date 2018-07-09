package team

type TeamKOStats struct {
	GoalsFor     int
	GoalsAgainst int
	Points       int
	PlayedGames  int
}

func NewTeamKOStats() TeamKOStats {
	g := TeamKOStats{
		GoalsFor:     0,
		GoalsAgainst: 0,
		Points:       0,
		PlayedGames:  0,
	}
	return g
}

func (g TeamKOStats) GoalDifference() int {
	return g.GoalsFor - g.GoalsAgainst
}
