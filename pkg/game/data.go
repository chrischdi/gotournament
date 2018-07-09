package game

type GroupGames map[string]*GroupGame
type KOGames map[string]*KOGame

var GroupGamesData GroupGames
var KOGamesData KOGames

func Reset() {
	GroupGamesData = GroupGames{}
	ResetKO()
}

func ResetKO() {
	KOGamesData = KOGames{}
}
