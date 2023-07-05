package config

import (
	"fmt"
	"time"
)

type Config struct {
	Start        time.Time
	GameDuration time.Duration
	KOPause      time.Duration
	Places       []*string
}

var ConfigData Config

func Reset() {
	n := time.Now()
	ConfigData = Config{
		time.Date(n.Year(), n.Month(), n.Day(), 18, 0, 0, 0, time.UTC),
		time.Duration(7) * time.Minute,
		time.Duration(15) * time.Minute,
		[]*string{},
	}
	for i := 0; i < 2; i++ {
		s := fmt.Sprintf("Tor %d", i+1)
		ConfigData.Places = append(ConfigData.Places, &s)
	}
}

func UpdateTime(starth, startm, duration, pause int) {
	n := time.Now()
	start := time.Date(n.Year(), n.Month(), n.Day(), starth, startm, 0, 0, time.UTC)
	ConfigData.Start = start
	ConfigData.GameDuration = time.Duration(duration) * time.Minute
	ConfigData.KOPause = time.Duration(pause) * time.Minute
}

func UpdatePlaces(nr int) {
	fmt.Printf("updating places to %d\n", nr)
	places := []*string{}
	for i := 1; i <= nr; i++ {
		s := fmt.Sprintf("Tor %d", i)
		fmt.Printf("adding places %s\n", s)
		places = append(places, &s)
	}
	ConfigData.Places = places
}