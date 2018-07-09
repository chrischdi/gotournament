package team

import (
	"fmt"

	"github.com/chrischdi/gotournament/pkg/helper"
)

type Team struct {
	UUID       string         `json: "uuid"`
	Name       string         `json: "name"`
	Dummy      bool           `json: "dummy"`
	GroupStats TeamGroupStats `json: "GroupStats"`
	KOStats    TeamKOStats    `json: "KOStats"`
}

func NewTeam(name string) *Team {
	t := &Team{
		UUID:       helper.CreateUUID(),
		Name:       name,
		Dummy:      false,
		GroupStats: NewTeamGroupStats(),
		KOStats:    NewTeamKOStats(),
	}
	TeamsData[t.UUID] = t
	return t
}

func DummyTeam() *Team {
	t := &Team{
		UUID:  helper.CreateUUID(),
		Dummy: true,
	}
	TeamsData[t.UUID] = t
	return t
}

func UpdateTeamName(uuid, name string) error {
	a, b := TeamsData[uuid]
	if !b {
		return fmt.Errorf("team with uuid=%s not found, %s", uuid, name)
	}
	a.Name = name
	return nil
}
