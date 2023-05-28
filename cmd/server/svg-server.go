package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/graph/svg"
)

type DayCount [366]int

func main() {
	r := mux.NewRouter()
	logger := func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	}
	r.Use(mux.MiddlewareFunc(logger))
	r.HandleFunc("/weekly.svg", func(w http.ResponseWriter, r *http.Request) {
		author := r.URL.Query().Get("author")
		highlight := r.URL.Query().Get("highlight")
		shouldHighlight := highlight != ""

		w.Header().Add("Content-Type", "text/html")
		repoPaths, err := commits.GetRepos()
		if err != nil {
			panic(err)
		}
		week, err := repoPaths.GetWeekFreq([]string{author})
		if err != nil {
			panic(err)
		}
		svg := svg.GetWeekSVG(week, shouldHighlight)
		svg.WriteTo(w)
	})
	r.HandleFunc("/stats.json", func(w http.ResponseWriter, r *http.Request) {
		year := time.Now().Year()
		yst := r.URL.Query().Get("year")
		author := r.URL.Query().Get("author")
		y, err := strconv.Atoi(yst)
		if err == nil {
			year = y
		}
		repoPaths, err := commits.GetRepos()
		if err != nil {
			panic(err)
		}
		freq, err := repoPaths.FrequencyChan(year, []string{author})
		if err != nil {
			panic(err)
		}
		b, _ := json.Marshal(freq)
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	})
	r.HandleFunc("/yearly.svg", func(w http.ResponseWriter, r *http.Request) {
		year := time.Now().Year()
		yst := r.URL.Query().Get("year")
		author := r.URL.Query().Get("author")
		highlight := r.URL.Query().Get("highlight")
		shouldHighlight := highlight != ""
		y, err := strconv.Atoi(yst)
		if err == nil {
			if year != y {
				shouldHighlight = false
			}
			year = y
		}
		repoPaths, err := commits.GetRepos()
		if err != nil {
			panic(err)
		}
		freq, err := repoPaths.FrequencyChan(year, []string{author})
		if err != nil {
			panic(err)
		}
		svg := svg.GetYearSVG(freq, shouldHighlight)
		w.Header().Add("Content-Type", "text/html")
		svg.WriteTo(w)
	})

	err := http.ListenAndServe(":8822", r)
	if err != nil {
		panic(err)
	}
}
