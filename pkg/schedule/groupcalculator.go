package schedule

import (
	"github.com/chrischdi/gotournament/pkg/config"
	"github.com/chrischdi/gotournament/pkg/game"
	"github.com/chrischdi/gotournament/pkg/group"
	"github.com/chrischdi/gotournament/pkg/matchday"
)

func CalculateGroupSchedule() {
	matchday.GroupMatchDaysList = []string{}

	gameTime := config.ConfigData.Start

	// check if at least one game
	matches := 0
	for _, g := range group.GroupsData.SortedByName() {
		matches += g.Matches()
	}
	if matches == 0 {
		return
	}

	// initialize a calculator per group
	cs := []*GroupMatchDayCalculator{}
	matchDays := 0
	for k, g := range group.GroupsData.SortedByName() {
		cs = append(cs, NewGroupMatchDayCalculator(g.Teams()))
		if r := cs[k].Rounds; r > matchDays {
			matchDays = r
		}
	}

	games := []string{}

	for i := 0; i < matchDays+1; i++ {
		for _, c := range cs {
			gs := c.Iterate()
			if gs != nil {
				games = append(games, gs.GamesList...)
			}
		}
	}

	p := 0
	md := matchday.NewGroupMatchday()
	md.Time = gameTime
	for _, g := range games {
		game.GroupGamesData[g].Place = config.ConfigData.Places[p]
		md.AddGame(game.GroupGamesData[g])

		p++
		if p >= len(config.ConfigData.Places) {
			p = 0
			matchday.GroupMatchDaysList = append(matchday.GroupMatchDaysList, md.UUID)
			md = matchday.NewGroupMatchday()
			gameTime = gameTime.Add(config.ConfigData.GameDuration)
			md.Time = gameTime
		}

	}
	if len(md.GamesList) > 0 {
		matchday.GroupMatchDaysList = append(matchday.GroupMatchDaysList, md.UUID)
	}
}
