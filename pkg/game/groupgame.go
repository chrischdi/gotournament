package game

import (
	"fmt"

	"github.com/chrischdi/gotournament/pkg/helper"
	"github.com/chrischdi/gotournament/pkg/team"
)

type GroupGame struct {
	UUID  string
	Place *string
	Stats GameStats
}

func NewGroupGame(A, B *team.Team, place *string) *GroupGame {
	g := &GroupGame{
		UUID:  helper.CreateUUID(),
		Place: place,
		Stats: NewGameStats(A.UUID, B.UUID),
	}

	GroupGamesData[g.UUID] = g
	return g
}

func (g *GroupGame) A() *team.Team { return team.TeamsData[g.Stats.AUUID] }
func (g *GroupGame) B() *team.Team { return team.TeamsData[g.Stats.BUUID] }

func evaluatePoints(a, b int) (int, int) {
	if a > b {
		return 3, 0
	} else if a < b {
		return 0, 3
	} else {
		return 1, 1
	}
}

func (g *GroupGame) SetGoals(a, b int) error {
	pa, pb := evaluatePoints(a, b)

	if g.Stats.Played {
		poa, pob := evaluatePoints(g.Stats.GoalsA, g.Stats.GoalsB)
		pa = pa - poa
		pb = pb - pob
	} else {
		g.A().GroupStats.PlayedGames = g.A().GroupStats.PlayedGames + 1
		g.B().GroupStats.PlayedGames = g.B().GroupStats.PlayedGames + 1
		g.Stats.Played = true
	}

	g.A().GroupStats.Points = g.A().GroupStats.Points + pa
	g.A().GroupStats.GoalsFor = g.A().GroupStats.GoalsFor + a - g.Stats.GoalsA
	g.A().GroupStats.GoalsAgainst = g.A().GroupStats.GoalsAgainst + b - g.Stats.GoalsB

	g.B().GroupStats.Points = g.B().GroupStats.Points + pb
	g.B().GroupStats.GoalsFor = g.B().GroupStats.GoalsFor + b - g.Stats.GoalsB
	g.B().GroupStats.GoalsAgainst = g.B().GroupStats.GoalsAgainst + a - g.Stats.GoalsA

	g.Stats.setGoals(a, b)
	return nil
}

func SetGroupGame(uuid string, a, b int) error {
	g := GroupGamesData[uuid]
	if g == nil {
		return fmt.Errorf("Game (uuid=%s) not found", uuid)
	}

	return g.SetGoals(a, b)
}
