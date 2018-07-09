package group

import (
	"sort"

	"github.com/chrischdi/gotournament/pkg/team"
)

func (g *Group) SortedGroup() []*team.Team {
	teams := g.Teams()
	sort.Sort(ByStats{teams})
	return teams
}

type SortGroupTeams []*team.Team

func (g SortGroupTeams) Len() int      { return len(g) }
func (g SortGroupTeams) Swap(i, j int) { g[i], g[j] = g[j], g[i] }

type ByStats struct{ SortGroupTeams }

// we got to reverse Less. If A got more points & goals -> return true
func (g ByStats) Less(i, j int) bool {
	a := g.SortGroupTeams[i]
	b := g.SortGroupTeams[j]

	if a.GroupStats.Points != b.GroupStats.Points {
		return a.GroupStats.Points > b.GroupStats.Points
	}

	// Points are equal. GoalDifference
	if a.GroupStats.GoalDifference() != b.GroupStats.GoalDifference() {
		return a.GroupStats.GoalDifference() > b.GroupStats.GoalDifference()
	}

	// GoalDifference is equal, shot only
	if a.GroupStats.GoalsFor != b.GroupStats.GoalsFor {
		return a.GroupStats.GoalsFor > b.GroupStats.GoalsFor
	}

	// Shot is equal, against only
	if a.GroupStats.GoalsAgainst != b.GroupStats.GoalsAgainst {
		return a.GroupStats.GoalsAgainst < b.GroupStats.GoalsAgainst
	}

	return false
}
