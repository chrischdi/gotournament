package team

type SortTeams []*Team

func (g SortTeams) Len() int      { return len(g) }
func (g SortTeams) Swap(i, j int) { g[i], g[j] = g[j], g[i] }

type ByKOGroupStats struct{ SortTeams }

// we got to reverse Less. If A got more points & goals -> return true
func (g ByKOGroupStats) Less(i, j int) bool {
	a := g.SortTeams[i]
	b := g.SortTeams[j]

	if a.KOStats.PlayedGames > b.KOStats.PlayedGames {
		return true
	} else if a.KOStats.PlayedGames < b.KOStats.PlayedGames {
		return false
	}

	if a.KOStats.Points < b.KOStats.Points {
		return false
	} else if a.KOStats.Points > b.KOStats.Points {
		return true
	}

	// Points in KO system are equal. Points in GroupStats
	if a.GroupStats.Points != b.GroupStats.Points {
		return a.GroupStats.Points > b.GroupStats.Points
	}

	// Points are equal. GoalDifference
	if a.GroupStats.GoalDifference()+a.KOStats.GoalDifference() != b.GroupStats.GoalDifference()+b.GroupStats.GoalDifference() {
		return a.GroupStats.GoalDifference()+a.KOStats.GoalDifference() > b.GroupStats.GoalDifference()+b.KOStats.GoalDifference()
	}

	// GoalDifference is equal, shot only
	if a.GroupStats.GoalsFor+a.KOStats.GoalsFor != b.GroupStats.GoalsFor+b.KOStats.GoalsFor {
		return a.GroupStats.GoalsFor+a.KOStats.GoalsFor > b.GroupStats.GoalsFor+b.KOStats.GoalsFor
	}

	// Shot is equal, against only
	if a.GroupStats.GoalsAgainst+a.KOStats.GoalsAgainst != b.GroupStats.GoalsAgainst+b.KOStats.GoalsAgainst {
		return a.GroupStats.GoalsAgainst+a.KOStats.GoalsAgainst < b.GroupStats.GoalsAgainst+b.KOStats.GoalsAgainst
	}
	return false
}
