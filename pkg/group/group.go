package group

import (
	"fmt"
	"sort"

	"github.com/chrischdi/gotournament/pkg/helper"
	"github.com/chrischdi/gotournament/pkg/team"
)

type Groups map[string]*Group

var GroupsData Groups

type Group struct {
	UUID          string
	Name          string
	TeamsList     []string
	MatchdaysList []string
}

func Reset() {
	GroupsData = map[string]*Group{}
}

func Dump() {
	for _, g := range GroupsData {
		fmt.Printf("Group >%s< (%s)\n", g.Name, g.UUID)
		for _, t := range g.Teams() {
			fmt.Printf("  Team >%s< (%s)\n", t.Name, t.UUID)
		}
	}
}

type ByName []*Group

func (s ByName) Len() int      { return len(s) }
func (s ByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (g *Groups) SortedByName() []*Group {
	groups := []*Group{}
	for _, g := range GroupsData {
		groups = append(groups, g)
	}
	sort.Sort(ByName(groups))

	return groups
}

func newGroup(name string) *Group {
	g := &Group{
		UUID:          helper.CreateUUID(),
		Name:          name,
		TeamsList:     []string{},
		MatchdaysList: []string{},
	}
	GroupsData[g.UUID] = g
	return g
}

func AddEmptyGroup(q int) error {
	letter := string('A' + int32(len(GroupsData)))
	name := "Group " + letter
	teams := []string{}

	for i := 1; i <= q; i++ {
		teams = append(teams, fmt.Sprintf("Team %s%d", letter, i))
	}
	AddGroup(name, teams)

	// Data.SaveFile()
	return nil
}

func AddGroup(name string, teams []string) {
	g := newGroup(name)
	for _, n := range teams {
		t := team.NewTeam(n)
		g.TeamsList = append(g.TeamsList, t.UUID)
		t.GroupStats.GroupUUID = g.UUID
	}
}

func (g Group) Teams() []*team.Team {
	teams := []*team.Team{}
	for _, i := range g.TeamsList {
		teams = append(teams, team.TeamsData[i])
	}
	return teams
}

func (g *Group) Matches() int {
	return len(g.TeamsList) * (len(g.TeamsList) - 1) / 2
}

// returns the array of teams which are currently placed at index i in each group
func (gs *Groups) GetIndex(i int) []*team.Team {
	teams := []*team.Team{}

	for _, g := range gs.SortedByName() {
		t := g.SortedGroup()[i]
		if t != nil {
			teams = append(teams, g.SortedGroup()[i])
		}
	}

	return teams
}

// returns the array of teams which are currently placed at index i in each group of maximum length
func (gs *Groups) GetIndexMaxLen(i, min int) []*team.Team {
	teams := []*team.Team{}

	sizes := map[int]int{}
	sizesArr := []int{}

	for _, g := range gs.SortedByName() {
		s, ok := sizes[len(g.TeamsList)]
		if !ok {
			sizes[len(g.TeamsList)] = 1
			sizesArr = append(sizesArr, len(g.TeamsList))
		} else {
			sizes[len(g.TeamsList)] = s + 1
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sizesArr)))

	for _, s := range sizesArr {
		if min < len(teams) {
			return teams
		}
		for _, g := range gs.SortedByName() {
			if len(g.TeamsList) == s {
				t := g.SortedGroup()[i]
				if t != nil {
					teams = append(teams, g.SortedGroup()[i])
				}
			}
		}
	}

	return teams
}
