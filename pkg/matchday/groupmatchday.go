package matchday

import (
	"time"

	"github.com/chrischdi/gotournament/pkg/config"
	"github.com/chrischdi/gotournament/pkg/game"
	"github.com/chrischdi/gotournament/pkg/helper"
)

// represents a Matchday in a group f.e. 1. gameday of bundesliga having 18 teams and 9 games at the same weekend.

type GroupMatchday struct {
	UUID      string
	GamesList []string
	NoGame    string
	Time      time.Time
}

func NewGroupMatchday() *GroupMatchday {
	m := &GroupMatchday{
		UUID:      helper.CreateUUID(),
		GamesList: []string{},
	}
	GroupMatchDaysData[m.UUID] = m
	return m
}

func (m *GroupMatchday) Games() []*game.GroupGame {
	r := []*game.GroupGame{}
	for _, i := range m.GamesList {
		r = append(r, game.GroupGamesData[i])
	}
	return r
}

func (m *GroupMatchday) AddGame(g *game.GroupGame) {
	if g.A().Dummy {
		m.NoGame = g.B().UUID
	} else if g.B().Dummy {
		m.NoGame = g.A().UUID
	} else {
		m.GamesList = append(m.GamesList, g.UUID)
	}
}

func UpdateGroupMatchdayTimes() {
	nextTime := config.ConfigData.Start
	for _, id := range GroupMatchDaysList {
		md := GroupMatchDaysData[id]
		md.Time = nextTime
		nextTime = nextTime.Add(config.ConfigData.GameDuration)
	}
	UpdateKOMatchdayTimes()
}
