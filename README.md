[![Build Status](https://travis-ci.org/chrischdi/gotournament.svg?branch=master)](https://travis-ci.org/chrischdi/gotournament)

# Description

This is a application which serves a http website to be used as a management application for tournaments.
It is a quick and dirty implementation, not supposed to be finished yet.
Lot's of stuff should be refactored.

The playing schedule is created according the algorithm at [wikipedia](https://de.wikipedia.org/wiki/Spielplan_(Sport))
The schedule for the KO part looks forward to matchup teams which had been in the same group as late as possible.

# Build

```
go build -o gotournament cmd/tournament.go
```

# Run

To run this application you will also need the `tpl` directory, which has to be located in the same directory like the executable.

```
./gotournament
```
