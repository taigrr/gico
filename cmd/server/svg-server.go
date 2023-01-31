package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/graph/svg"
)

type DayCount [366]int

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weekly.svg", func(w http.ResponseWriter, r *http.Request) {
		author := r.URL.Query().Get("author")
		w.Header().Add("Content-Type", "text/html")
		now := time.Now()
		year := now.Year()
		repoPaths, err := commits.GetMRRepos()
		if err != nil {
			panic(err)
		}
		freq, err := repoPaths.FrequencyChan(year, []string{author})
		if err != nil {
			panic(err)
		}
		today := now.YearDay() - 1
		fmt.Println(today)
		if today < 6 {
			curYear := year - 1
			curFreq, err := repoPaths.FrequencyChan(curYear, []string{author})
			if err != nil {
				panic(err)
			}
			freq = append(curFreq, freq...)
			today += 365
			if curYear%4 == 0 {
				today++
			}
		}
		fmt.Println(freq)

		week := freq[today-6 : today+1]
		svg := svg.GetWeekSVG(week)
		svg.WriteTo(w)
	})
	r.HandleFunc("/yearly.svg", func(w http.ResponseWriter, r *http.Request) {
		year := time.Now().Year()
		yst := r.URL.Query().Get("year")
		author := r.URL.Query().Get("author")
		y, err := strconv.Atoi(yst)
		if err == nil {
			year = y
		}
		repoPaths, err := commits.GetMRRepos()
		if err != nil {
			panic(err)
		}
		freq, err := repoPaths.FrequencyChan(year, []string{author})
		if err != nil {
			panic(err)
		}
		svg := svg.GetYearSVG(freq)
		w.Header().Add("Content-Type", "text/html")
		svg.WriteTo(w)
	})

	http.ListenAndServe(":8080", r)
}
