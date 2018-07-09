package game

import (
	"fmt"
	"time"

	"github.com/chrischdi/gotournament/pkg/helper"
	"github.com/chrischdi/gotournament/pkg/team"
)

type KOGame struct {
	GameName  string
	Stats     GameStats
	UUID      string
	Place     string
	PreGAUUID string
	PreGBUUID string
	Time      time.Time
}

func NewKOGame() *KOGame {
	g := &KOGame{
		UUID:      helper.CreateUUID(),
		GameName:  "",
		Place:     "",
		Stats:     NewGameStats("", ""),
		PreGAUUID: "",
		PreGBUUID: "",
	}

	KOGamesData[g.UUID] = g
	return g
}

func (g *KOGame) PreGA() *KOGame {
	return KOGamesData[g.PreGAUUID]
}
func (g *KOGame) PreGB() *KOGame {
	return KOGamesData[g.PreGBUUID]
}

func (g *KOGame) A() *team.Team {
	if g.Stats.AUUID != "" {
		return team.TeamsData[g.Stats.AUUID]
	}
	return g.PreGA().Winner()
}
func (g *KOGame) B() *team.Team {
	if g.Stats.BUUID != "" {
		return team.TeamsData[g.Stats.BUUID]
	}
	return g.PreGB().Winner()
}

func (g *KOGame) NameA() string {
	a := g.A()
	if a != nil {
		return a.Name
	}

	return fmt.Sprintf("Sieger des Spiels %v", g.PreGA().GameName)
}

func (g *KOGame) NameB() string {
	b := g.B()
	if b != nil {
		return b.Name
	}

	return fmt.Sprintf("Sieger des Spiels %v", g.PreGB().GameName)
}

func (g *KOGame) Winner() *team.Team {
	if g.Stats.Played {
		if g.Stats.GoalsA > g.Stats.GoalsB {
			return g.A()
		} else if g.Stats.GoalsA < g.Stats.GoalsB {
			return g.B()
		}
	}
	return nil
}

func (g *KOGame) IsPlayable() bool {
	if g.A() != nil && g.B() != nil {
		return true
	}
	return false
}

func (g *KOGame) SetGoals(a, b int) error {
	pa, pb := evaluatePoints(a, b)

	if g.Stats.Played {
		poa, pob := evaluatePoints(g.Stats.GoalsA, g.Stats.GoalsB)
		pa = pa - poa
		pb = pb - pob
	} else {
		g.A().KOStats.PlayedGames = g.A().KOStats.PlayedGames + 1
		g.B().KOStats.PlayedGames = g.B().KOStats.PlayedGames + 1
		g.Stats.Played = true
	}

	g.A().KOStats.Points = g.A().KOStats.Points + pa
	g.A().KOStats.GoalsFor = g.A().KOStats.GoalsFor + a - g.Stats.GoalsA
	g.A().KOStats.GoalsAgainst = g.A().KOStats.GoalsAgainst + b - g.Stats.GoalsB

	g.B().KOStats.Points = g.B().KOStats.Points + pb
	g.B().KOStats.GoalsFor = g.B().KOStats.GoalsFor + b - g.Stats.GoalsB
	g.B().KOStats.GoalsAgainst = g.B().KOStats.GoalsAgainst + a - g.Stats.GoalsA

	g.Stats.setGoals(a, b)
	return nil
}

func SetKOGame(uuid string, a, b int) error {
	g := KOGamesData[uuid]
	if g == nil {
		return fmt.Errorf("Game (uuid=%s) not found", uuid)
	}

	return g.SetGoals(a, b)
}
