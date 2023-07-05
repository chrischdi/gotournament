package frontend

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/chrischdi/gotournament/pkg/matchday"
	"github.com/chrischdi/gotournament/tpl"

	"github.com/chrischdi/gotournament/pkg/config"
	"github.com/chrischdi/gotournament/pkg/game"
)

var funcMap = template.FuncMap{
	"add":               add,
	"multipl":           multipl,
	"lastMatchdaygame":  lastMatchdaygame,
	"lastPlace":         lastPlace,
	"getMatchDay":       getMatchDay,
	"allGroupGamesDone": allGroupGamesDone,
}

func add(a, b int) int {
	return a + b
}
func multipl(a, b int) int {
	return a * b
}
func lastMatchdaygame(elems []*game.GroupGame, k int) bool {
	return len(elems)-1 == k
}
func lastPlace(k int) bool {
	return len(config.ConfigData.Places)-1 == k
}
func getMatchDay(i string) *matchday.GroupMatchday {
	return matchday.GroupMatchDaysData[i]
}
func allGroupGamesDone() bool {
	for _, g := range game.GroupGamesData {
		if !g.Stats.Played {
			if g.A().Dummy || g.B().Dummy {
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func writeTemplate(w http.ResponseWriter, tplName string, vars interface{}) error {
	t, err := template.New(tplName).Funcs(funcMap).ParseFS(tpl.Content, tplName)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return fmt.Errorf("error opening template %s: %v", tplName, err)
	}

	err = t.Execute(w, vars)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return err
}

func HtmlInput(t, name, value, add string) string {
	if t == "number" {
		add = add + " style=\"text-align: center; width: 3em;\""
	}
	return fmt.Sprintf("<input type=\"%s\" name=\"%s\" value=\"%v\" %s>", t, name, value, add)
}

func HtmlSingleForm(action, name, value string) string {
	button := HtmlInput("submit", name, value, "")
	return fmt.Sprintf(`<form action="%s">%s</form>`, action, button)
}
