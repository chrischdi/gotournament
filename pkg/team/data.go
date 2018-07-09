package team

import "sort"

type Teams map[string]*Team

var TeamsData Teams

func Reset() {
	TeamsData = map[string]*Team{}
}

func (td *Teams) SortedTeams() []*Team {
	teams := []*Team{}
	for _, t := range *td {
		if !t.Dummy {
			teams = append(teams, t)
		}
	}
	sort.Stable(ByKOGroupStats{teams})
	return teams
}

func ResetKOStats() {
	for _, t := range TeamsData {
		t.KOStats = NewTeamKOStats()
	}
}
