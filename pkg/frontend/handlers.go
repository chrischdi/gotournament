package frontend

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/chrischdi/gotournament/pkg/config"
	"github.com/chrischdi/gotournament/pkg/data"
	"github.com/chrischdi/gotournament/pkg/game"
	"github.com/chrischdi/gotournament/pkg/group"
	"github.com/chrischdi/gotournament/pkg/matchday"
	"github.com/chrischdi/gotournament/pkg/schedule"
	"github.com/chrischdi/gotournament/pkg/team"
)

func registerHandlers() {
	http.HandleFunc("/", hIndex)
	http.HandleFunc("/matchplan", hMatchplan)
	http.HandleFunc("/matchplan/set", hMatchplanSet)
	http.HandleFunc("/matchplan/setfoo", hMatchplanSetFoo)
	http.HandleFunc("/setup/", hSetup)
	http.HandleFunc("/setup/reset/", hReset)
	http.HandleFunc("/setup/reset/confirm/", hResetConfirm)
	http.HandleFunc("/setup/add/", hSetupAdd)
	http.HandleFunc("/setup/update/", hSetupUpdate)
	http.HandleFunc("/save/tofile/", hSaveToFile)
	http.HandleFunc("/load/fromfile/", hLoadFromFile)
	http.HandleFunc("/table", hTable)
	http.HandleFunc("/fulltable", hFullTable)
	http.HandleFunc("/ko", hKo)
	http.HandleFunc("/ko/set", hKoSet)
	http.HandleFunc("/ko/setfoo", hKoSetFoo)
	http.HandleFunc("/ko/generate", hKoGenerate)
	http.HandleFunc("/ko/generate/confirm", hKoGenerateConfirm)
}

// Serves: /
func hIndex(w http.ResponseWriter, r *http.Request) {
	head(w, "")
	foot(w)
}

func hSetup(w http.ResponseWriter, r *http.Request) {
	head(w, "Setup")
	defer foot(w)
	setup(w, r)
}

func setup(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	var err error

	fmt.Fprintf(w, "<h1>Setup</h1>\n")
	fmt.Fprintf(w, "<h2>General</h2>\n")
	fmt.Fprintf(w, HtmlSingleForm("/setup/reset", "reset", "Reset"))

	fmt.Fprintf(w, "<h2>Time</h2>\n")
	t = template.New("timeSetup")
	t, err = template.ParseFiles(filepath.Join("tpl", "timeSetup.tmpl"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	t.Execute(w,
		map[string]interface{}{
			"Start":    fmt.Sprintf("%02d:%02d", config.ConfigData.Start.Hour(), config.ConfigData.Start.Minute()),
			"Duration": config.ConfigData.GameDuration.Minutes(),
			"Pause":    config.ConfigData.KOPause.Minutes(),
		})

	writeTemplate(w, "modifyGroups.tmpl", data.GetData())

	fmt.Fprintf(w, "<h2>Add Groups</h2>\n")
	t = template.New("addGroup")
	t, err = template.ParseFiles(filepath.Join("tpl", "addGroup.tmpl"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	t.Execute(w, nil)
}

func hReset(w http.ResponseWriter, r *http.Request) {
	head(w, "Setup - Reset confirmation")
	defer foot(w)

	fmt.Fprintf(w, "<h1>Reset confirmation</h1>\n")
	fmt.Fprintf(w, "<p>Do you really want to reset all data?</p>\n")
	fmt.Fprintf(w, "<p>All currently saved results will be lost!</p>\n")
	fmt.Fprintf(w, HtmlSingleForm("/setup/reset/confirm", "yes", "Yes"))
	fmt.Fprintf(w, HtmlSingleForm("/setup", "no", "No"))
}

// Serves: /setup/reset/confirm
func hResetConfirm(w http.ResponseWriter, r *http.Request) {
	data.Reset()
	http.Redirect(w, r, "/setup", 307)
}

// Serves: /table
func hTable(w http.ResponseWriter, r *http.Request) {
	head(w, "Table")
	defer foot(w)
	table(w)
}

func table(w http.ResponseWriter) {
	writeTemplate(w, "table.tmpl", data.GetData())
}

func hFullTable(w http.ResponseWriter, r *http.Request) {
	head(w, "Table")
	defer foot(w)
	writeTemplate(w, "fulltable.tmpl", data.GetData())
}

func hSetupAdd(w http.ResponseWriter, r *http.Request) {
	head(w, "Setup - Reset confirmation")
	defer foot(w)

	err := setupAdd(w, r)
	if err != nil {
		fmt.Fprintf(w, "<p>%v</p>\n", err)
	}
	setup(w, r)
}

func setupAdd(w http.ResponseWriter, r *http.Request) error {
	object := strings.Split(r.URL.EscapedPath(), "/")[3]
	switch object {
	case "":
	case "group":
		qt, err := strconv.Atoi(r.FormValue("quantityTeams"))
		if err != nil {
			return fmt.Errorf("unable to parse quantityTeams: %v", err)
		}
		qg, err := strconv.Atoi(r.FormValue("quantityGroups"))
		if err != nil {
			return fmt.Errorf("unable to parse quantityGroups: %v", err)
		}
		if qt < 2 {
			return fmt.Errorf("error adding group: Amount of teams needs to be more than 2")
		}
		if qg < 1 {
			return fmt.Errorf("error adding group: Amount of groups needs to be more than 2")
		}

		for g := 0; g < qg; g++ {
			err = group.AddEmptyGroup(qt)
			if err != nil {
				return fmt.Errorf("error adding group: %v", err)
			}
		}
		game.Reset()
		matchday.Reset()
		schedule.CalculateGroupSchedule()
		data.SaveFile()
	default:
		return fmt.Errorf("error: unknown object to add.")
	}
	return nil
}

func hSaveToFile(w http.ResponseWriter, r *http.Request) {
	data.SaveFile()

	s, err := os.Executable()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	dir, _ := filepath.Split(s)

	path := filepath.Join(dir, "db.json")

	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set(`Content-Disposition`, `attachment; filename="db.json"`)

	http.ServeFile(w, r, path)
}

func hLoadFromFile(w http.ResponseWriter, r *http.Request) {
	head(w, "Table")
	if r.Method == "POST" {
		f, _, err := r.FormFile("database")
		defer f.Close()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		b := &bytes.Buffer{}
		b.ReadFrom(f)

		d := data.DataHelper{}
		err = d.Restore(b.Bytes())
		if err != nil {
			fmt.Fprintf(w, "<h2>error loading file:</h2><p>%v</p>", err)
			writeTemplate(w, "loadfromfile.tmpl", nil)
		} else {
			fmt.Fprintf(w, "<h2>file successfully loaded!</h2>")
		}
	} else {
		writeTemplate(w, "loadfromfile.tmpl", nil)
	}
	foot(w)
	data.SaveFile()
}

// Serves: /setup/update/
func hSetupUpdate(w http.ResponseWriter, r *http.Request) {
	head(w, "Setup")
	defer foot(w)

	err := setupUpdate(w, r)
	if err != nil {
		fmt.Fprintf(w, "<p>%v</p>\n", err)
	}
	setup(w, r)
	data.SaveFile()
}

func setupUpdate(w http.ResponseWriter, r *http.Request) error {
	object := strings.Split(r.URL.EscapedPath(), "/")[3]
	switch object {
	case "":
	case "time":
		start := strings.Split(r.FormValue("start"), ":")
		if len(start) != 2 {
			return fmt.Errorf("error updating time: wrong format on start time")
		}
		starth, err := strconv.Atoi(start[0])
		if err != nil {
			return fmt.Errorf("error updating time: %v", err)
		}
		startm, err := strconv.Atoi(start[1])
		if err != nil {
			return fmt.Errorf("error updating time: %v", err)
		}
		duration, err := strconv.Atoi(r.FormValue("duration"))
		if err != nil {
			return fmt.Errorf("error updating time: %v", err)
		}
		pause, err := strconv.Atoi(r.FormValue("pause"))
		if err != nil {
			return fmt.Errorf("error updating time: %v", err)
		}
		config.UpdateTime(starth, startm, duration, pause)
		matchday.UpdateGroupMatchdayTimes()
		data.SaveFile()
		// Tournament().UpdateTime(starth, startm, duration, pause)
		// Tournament().UpdateKOTime()

	case "group":
		for _, v := range strings.Split(r.FormValue("uuids"), ",") {
			err := team.UpdateTeamName(v, r.FormValue(v))
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("error: unknown object to update.")
	}
	return nil
}

// Serves: /matchplan
func hMatchplan(w http.ResponseWriter, r *http.Request) {
	head(w, "Matchplan")
	defer foot(w)
	writeTemplate(w, "matchplan.tmpl", data.GetData())
}

func hMatchplanSetFoo(w http.ResponseWriter, r *http.Request) {
	i := 1
	for _, g := range game.GroupGamesData {
		fmt.Printf("Setting %s %d:0\n", g.UUID, i)
		if g.A().Dummy || g.B().Dummy {
			continue
		}
		g.SetGoals(i, 0)
		i = i + 1
	}
	http.Redirect(w, r, "/matchplan", 307)
	data.SaveFile()
}

func hMatchplanSet(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	goalsA := r.FormValue("GoalsA")
	goalsB := r.FormValue("GoalsB")
	if uuid == "" || goalsA == "" || goalsB == "" {
		errs(w, fmt.Errorf("Values missing! %v %v %v", uuid, goalsA, goalsB))
	}

	a, err := strconv.Atoi(goalsA)
	if err != nil {
		errs(w, err)
	}

	b, err := strconv.Atoi(goalsB)
	if err != nil {
		errs(w, err)
	}

	game.SetGroupGame(uuid, a, b)
	http.Redirect(w, r, "/matchplan", 307)
	data.SaveFile()
}

func hKo(w http.ResponseWriter, r *http.Request) {
	head(w, "KO")
	defer foot(w)

	writeTemplate(w, "ko.tmpl", data.GetData())
}

func hKoSet(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	goalsA := r.FormValue("GoalsA")
	goalsB := r.FormValue("GoalsB")
	if uuid == "" || goalsA == "" || goalsB == "" {
		errs(w, fmt.Errorf("Values missing! %v %v %v", uuid, goalsA, goalsB))
	}

	a, err := strconv.Atoi(goalsA)
	if err != nil {
		errs(w, err)
	}

	b, err := strconv.Atoi(goalsB)
	if err != nil {
		errs(w, err)
	}

	game.SetKOGame(uuid, a, b)
	http.Redirect(w, r, "/ko", 307)
	data.SaveFile()
}

func hKoSetFoo(w http.ResponseWriter, r *http.Request) {
	i := 1
	for _, id := range matchday.KOMatchdaysList {
		md := matchday.KOMatchdaysData[id]
		for _, g := range md.Games() {
			g.SetGoals(i, 0)
			i = i + 1
		}
	}
	http.Redirect(w, r, "/ko", 307)
	data.SaveFile()
}

// Serves: /ko/generate
func hKoGenerate(w http.ResponseWriter, r *http.Request) {
	head(w, "KO Plan generieren")
	defer foot(w)

	writeTemplate(w, "ko-generate.tmpl", nil)
}

// Serves: /ko/generate/confirm
func hKoGenerateConfirm(w http.ResponseWriter, r *http.Request) {
	bestof, err := strconv.Atoi(r.FormValue("bestof"))
	if err != nil {
		errs(w, err)
	}
	_ = bestof

	team.ResetKOStats()
	game.ResetKO()
	matchday.ResetKO()
	data.SaveFile()

	schedule.CalculateKOSchedule(bestof)

	http.Redirect(w, r, "/ko", 307)
	data.SaveFile()
}
