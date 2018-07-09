package matchday

import (
	"github.com/chrischdi/gotournament/pkg/config"
	"github.com/chrischdi/gotournament/pkg/game"
	"github.com/chrischdi/gotournament/pkg/helper"
)

type KOMatchday struct {
	UUID      string
	GamesList []string
}

func NewKOMatchday() *KOMatchday {
	m := &KOMatchday{
		UUID:      helper.CreateUUID(),
		GamesList: []string{},
	}
	KOMatchdaysData[m.UUID] = m
	KOMatchdaysList = append([]string{m.UUID}, KOMatchdaysList...)
	return m
}

func (m *KOMatchday) Games() []*game.KOGame {
	r := []*game.KOGame{}
	for _, i := range m.GamesList {
		r = append(r, game.KOGamesData[i])
	}
	return r
}

func (m *KOMatchday) AddGame(g *game.KOGame) {
	m.GamesList = append(m.GamesList, g.UUID)
}

func UpdateKOMatchdayTimes() {
	lastGroupMatchdayUUID := GroupMatchDaysList[len(GroupMatchDaysList)-1]
	nextTime := GroupMatchDaysData[lastGroupMatchdayUUID].Time.Add(config.ConfigData.KOPause)
	// now set the time

	nextPlaceIdx := 0

	for _, md := range KOMatchdaysList {
		for _, g := range KOMatchdaysData[md].GamesList {
			if len(KOMatchdaysData[md].GamesList) > 2 {
				game.KOGamesData[g].Place = *config.ConfigData.Places[nextPlaceIdx]
				game.KOGamesData[g].Time = nextTime
				nextPlaceIdx++
				if nextPlaceIdx >= len(config.ConfigData.Places) {
					nextTime = nextTime.Add(config.ConfigData.GameDuration)
					nextPlaceIdx = 0
				}
			} else {
				game.KOGamesData[g].Place = *config.ConfigData.Places[0]
				game.KOGamesData[g].Time = nextTime
				nextTime = nextTime.Add(config.ConfigData.GameDuration)

			}
		}
	}
}
