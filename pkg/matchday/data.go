package matchday

type GroupMatchDays map[string]*GroupMatchday
type KOMatchDays map[string]*KOMatchday

var GroupMatchDaysData GroupMatchDays
var GroupMatchDaysList []string
var KOMatchdaysData KOMatchDays
var KOMatchdaysList []string

func Reset() {
	GroupMatchDaysData = GroupMatchDays{}
	GroupMatchDaysList = []string{}
	ResetKO()
}

func ResetKO() {
	KOMatchdaysData = KOMatchDays{}
	KOMatchdaysList = []string{}
}
