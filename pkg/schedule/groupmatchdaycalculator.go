package schedule

import (
	"github.com/chrischdi/gotournament/pkg/game"
	"github.com/chrischdi/gotournament/pkg/matchday"
	"github.com/chrischdi/gotournament/pkg/team"
)

type GroupMatchDayCalculator struct {
	Teams  []*team.Team
	I      int
	Rounds int
	Helper groupPlanHelper
}

func NewGroupMatchDayCalculator(teams []*team.Team) *GroupMatchDayCalculator {
	t := []*team.Team{}

	// add teams and dummy to t
	for _, team := range teams {
		t = append(t, team)
	}
	if len(t)%2 == 1 {
		t = append(t, team.DummyTeam())
	}

	h := groupPlanHelper{
		x: t[len(t)/2-1],
		y: t[len(t)-1],
	}
	for i := 0; i < len(t)/2-1; i++ {
		h.upper = append(h.upper, t[i])
		h.lower = append(h.lower, t[i+len(t)/2])
	}

	return &GroupMatchDayCalculator{
		Teams:  t,
		I:      0,
		Helper: h,
		Rounds: len(t) - 1,
	}
}

func (c *GroupMatchDayCalculator) Iterate() *matchday.GroupMatchday {
	if c.I >= len(c.Teams)-1 {
		return nil
	}
	m := matchday.NewGroupMatchday()

	if (c.I % 2) == 0 {
		m.AddGame(game.NewGroupGame(c.Helper.x, c.Helper.y, nil))
	}

	for i := 0; i < len(c.Helper.lower); i++ {
		m.AddGame(game.NewGroupGame(c.Helper.upper[i], c.Helper.lower[i], nil))
	}

	if (c.I % 2) == 1 {
		m.AddGame(game.NewGroupGame(c.Helper.a, c.Helper.b, nil))
	}
	c.I++
	c.Helper.rotatePlan()
	return m
}

// Helper for calculating matchdays
type groupPlanHelper struct {
	a     *team.Team
	b     *team.Team
	upper []*team.Team
	lower []*team.Team
	x     *team.Team
	y     *team.Team
}

// Rotate to the next matchday setup
func (h *groupPlanHelper) rotatePlan() {
	if h.x != nil && h.y != nil {
		h.a = h.y
		h.b = h.upper[0]
		h.upper = h.upper[1:]
		h.upper = append(h.upper, h.x)
		h.x = nil
		h.y = nil
	} else {
		h.y = h.a
		h.x = h.lower[len(h.lower)-1]
		newLow := []*team.Team{h.b}
		h.lower = append(newLow, h.lower...)
		h.lower = h.lower[:len(h.lower)-1]
		h.a = nil
		h.b = nil
	}
}
