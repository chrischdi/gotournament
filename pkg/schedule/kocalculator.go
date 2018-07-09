package schedule

import (
	"fmt"
	"sort"

	"github.com/chrischdi/gotournament/pkg/game"
	"github.com/chrischdi/gotournament/pkg/group"
	"github.com/chrischdi/gotournament/pkg/matchday"
	"github.com/chrischdi/gotournament/pkg/team"
)

func resortTeams(old []*team.Team) []*team.Team {
	new := []*team.Team{}

	counter := map[string]int{}
	for _, t := range old {
		counter[t.GroupStats.GroupUUID] = counter[t.GroupStats.GroupUUID] + 1
	}

	max := 0
	for _, v := range counter {
		if max < v {
			max = v
		}
	}

	groupOrder := []string{}
	for i := max; i > 0; i-- {
		for uuid, v := range counter {
			if i == v {
				groupOrder = append(groupOrder, uuid)
				counter[uuid] = -1
			}
		}
	}

	for _, g := range groupOrder {
		for _, t := range old {
			if g == t.GroupStats.GroupUUID {
				new = append(new, t)
			}
		}
	}

	return new
}

func CalculateKOSchedule(bestof int) error {
	// get the best `bestof` teams

	teamsByPlace := [][]*team.Team{}
	curindex := 0
	curteams := 0

	for {
		tmp := group.GroupsData.GetIndex(curindex)
		if len(tmp)+curteams <= bestof {
			teamsByPlace = append(teamsByPlace, tmp)
			curteams = curteams + len(tmp)
			curindex++
		} else {
			break
		}
	}

	// fill up the rest with the best of the next place, only including the biggest size
	if curteams < bestof {
		tmp := group.GroupsData.GetIndexMaxLen(curindex, bestof-curteams)
		sort.Sort(group.ByStats{tmp})
		teamsByPlace = append(teamsByPlace, tmp[:bestof-curteams])
		curteams = curteams + bestof - curteams
	}

	matchday.ResetKO()
	game.ResetKO()

	GameNameFormat := func(round, nr int) string {
		return fmt.Sprintf("%02dG%03d", round+1, nr+1)
	}

	mds := []*matchday.KOMatchday{matchday.NewKOMatchday()}
	final := game.NewKOGame()
	final.GameName = GameNameFormat(0, 0)
	mds[0].AddGame(final)

	// create empty games and add to the front of mds
	// start at 4 because 2 (final) is already added
	for k := 4; k <= bestof; k = k * 2 {
		md := matchday.NewKOMatchday()
		// create 2 games per previous existing 1 game
		for i, g := range mds[0].Games() {
			newa := game.NewKOGame()
			newa.GameName = GameNameFormat(k, i*2)
			newb := game.NewKOGame()
			newb.GameName = GameNameFormat(k, i*2+1)
			g.PreGAUUID = newa.UUID
			g.PreGBUUID = newb.UUID
			md.AddGame(newa)
			md.AddGame(newb)
		}
		mds = append([]*matchday.KOMatchday{md}, mds...)
	}

	// we now got a empty KO Schedule.
	// Lets fill the first round

	// make two groups of teams the better and the less better half
	bteams := []*team.Team{}
	oteams := []*team.Team{}

	toFill := len(mds[0].GamesList)
	for _, ts := range teamsByPlace {
		if len(ts) < toFill {
			bteams = append(bteams, ts...)
			toFill = toFill - len(ts)
		} else if toFill > 0 {
			sort.Sort(group.ByStats{ts})
			bteams = append(bteams, ts[:toFill]...)
			oteams = append(oteams, ts[toFill:]...)
			toFill = toFill - len(ts[:toFill])
		} else {
			oteams = append(oteams, ts...)
		}
	}

	// resort teams to prevent teams of same group to match up too early
	bteams = resortTeams(bteams)
	oteams = resortTeams(oteams)

	h := NewKOMatchdayHelper(final)
	for i := 0; i < len(bteams); i++ {
		if !h.Insert(bteams[i], true, false) {
			if !h.Insert(bteams[i], true, true) {
				return fmt.Errorf("error: unable to insert team %s\n", bteams[i].Name)
			}
		}
	}
	for i := 0; i < len(oteams); i++ {
		if !h.Insert(oteams[i], false, false) {
			if !h.Insert(oteams[i], false, true) {
				return fmt.Errorf("error: unable to insert team %s\n", bteams[i].Name)
			}
		}
	}

	matchday.UpdateKOMatchdayTimes()

	return nil
}

type KOMatchdayIterator struct {
	treeTop *KOMatchdayHelper
}

func NewKOMatchdayIterator(final *game.KOGame) *KOMatchdayIterator {
	return &KOMatchdayIterator{
		NewKOMatchdayHelper(final),
	}
}

type KOMatchdayHelper struct {
	PreA *KOMatchdayHelper
	PreB *KOMatchdayHelper
	// indexes marker
	PreAIdxs []int
	PreBIdxs []int
	Game     *game.KOGame
}

func NewKOMatchdayHelper(g *game.KOGame) *KOMatchdayHelper {
	h := &KOMatchdayHelper{
		nil,
		nil,
		[]int{},
		[]int{},
		g,
	}

	if h.Game.PreGAUUID != "" {
		h.PreA = NewKOMatchdayHelper(h.Game.PreGA())
	}
	if h.Game.PreGBUUID != "" {
		h.PreB = NewKOMatchdayHelper(h.Game.PreGB())
	}

	return h
}

// makes recursion to all pre games and returns list of Group UUID's
func (h *KOMatchdayHelper) GetPreGroupUUIDsAndCount() (map[string]int, int) {
	ret := map[string]int{}

	if h != nil {
		if h.PreA != nil {
			m, _ := h.PreA.GetPreGroupUUIDsAndCount()
			for k, v := range m {
				ret[k] = ret[k] + v
			}
		} else if h.Game != nil && h.Game.Stats.AUUID != "" {
			guuid := h.Game.A().GroupStats.GroupUUID
			ret[guuid] = ret[guuid] + 1
		}
	}

	if h != nil {
		if h.PreB != nil {
			m, _ := h.PreB.GetPreGroupUUIDsAndCount()
			for k, v := range m {
				ret[k] = ret[k] + v
			}
		} else if h.Game != nil && h.Game.Stats.BUUID != "" {
			guuid := h.Game.B().GroupStats.GroupUUID
			ret[guuid] = ret[guuid] + 1
		}
	}

	i := 0
	for _, v := range ret {
		i = i + v
	}

	return ret, i
}

func (h *KOMatchdayHelper) Insert(t *team.Team, doA, force bool) bool {
	if h == nil {
		return false
	}
	if doA {
		if h.PreA == nil {
			if h.Game.Stats.AUUID != "" {
				return false
			}
			h.Game.Stats.AUUID = t.UUID
			return true
		}
	} else {
		if h.PreB == nil {
			if h.Game.Stats.BUUID != "" {
				return false
			}
			h.Game.Stats.BUUID = t.UUID
			return true
		}
	}

	guuid := t.GroupStats.GroupUUID

	a, alen := h.PreA.GetPreGroupUUIDsAndCount()
	b, blen := h.PreB.GetPreGroupUUIDsAndCount()

	if !force {
		if a[guuid] < b[guuid] {
			return h.PreA.Insert(t, doA, force)
		}
		if a[guuid] > b[guuid] {
			return h.PreB.Insert(t, doA, force)
		}
	}

	if alen <= blen {
		return h.PreA.Insert(t, doA, force)
	} else {
		return h.PreB.Insert(t, doA, force)
	}
}
