package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chrischdi/gotournament/pkg/config"
	"github.com/chrischdi/gotournament/pkg/game"
	"github.com/chrischdi/gotournament/pkg/group"
	"github.com/chrischdi/gotournament/pkg/matchday"
	"github.com/chrischdi/gotournament/pkg/team"
)

func Reset() {
	group.Reset()
	team.Reset()
	game.Reset()
	config.Reset()
	matchday.Reset()
}

type DataHelper struct {
	Config             config.Config
	Groups             group.Groups
	Teams              team.Teams
	GroupGames         game.GroupGames
	GroupMatchDays     matchday.GroupMatchDays
	GroupMatchDaysList []string
	KOGames            game.KOGames
	KOMatchDays        matchday.KOMatchDays
	KOMatchDaysList    []string
}

func GetData() *DataHelper {
	return &DataHelper{
		Config:             config.ConfigData,
		Groups:             group.GroupsData,
		Teams:              team.TeamsData,
		GroupGames:         game.GroupGamesData,
		GroupMatchDays:     matchday.GroupMatchDaysData,
		GroupMatchDaysList: matchday.GroupMatchDaysList,
		KOGames:            game.KOGamesData,
		KOMatchDays:        matchday.KOMatchdaysData,
		KOMatchDaysList:    matchday.KOMatchdaysList,
	}
}

func (d *DataHelper) restoreData() {
	config.ConfigData = d.Config
	group.GroupsData = d.Groups
	team.TeamsData = d.Teams
	game.GroupGamesData = d.GroupGames
	matchday.GroupMatchDaysData = d.GroupMatchDays
	matchday.GroupMatchDaysList = d.GroupMatchDaysList
	game.KOGamesData = d.KOGames
	matchday.KOMatchdaysData = d.KOMatchDays
	matchday.KOMatchdaysList = d.KOMatchDaysList
}

func (d *DataHelper) JSON() string {
	b, err := json.Marshal(d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return ""
	}
	return string(b)
}

func Dump() {
	fmt.Printf("%s\n", GetData().JSON())
}

func SaveFile() error {
	d := GetData()
	s, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	dir, _ := filepath.Split(s)

	path := filepath.Join(dir, "db.json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error opening file %s: %v", path, err)
	}
	defer f.Close()

	_, err = f.WriteString(d.JSON())
	if err != nil {
		return fmt.Errorf("Error writing file %s: %v", path, err)
	}

	return nil
}

func RestoreFile() error {
	d := &DataHelper{}
	s, err := os.Executable()
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	dir, _ := filepath.Split(s)

	path := filepath.Join(dir, "db.json")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return err
	}

	return d.Restore(b)
}

func (d *DataHelper) Restore(b []byte) error {
	Reset()
	err := json.Unmarshal(b, &d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
	d.restoreData()
	return nil
}
