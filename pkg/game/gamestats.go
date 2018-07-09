package game

type GameStats struct {
	GoalsA int
	GoalsB int
	Played bool
	AUUID  string
	BUUID  string
}

func NewGameStats(a, b string) GameStats {
	return GameStats{
		GoalsA: 0,
		GoalsB: 0,
		Played: false,
		AUUID:  a,
		BUUID:  b,
	}
}

func (s *GameStats) setGoals(a, b int) {
	s.GoalsA = a
	s.GoalsB = b
}
