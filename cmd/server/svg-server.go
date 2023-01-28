package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/gitgraph/svg"
)

type DayCount [366]int

func init() {
	rand.Seed(time.Now().UnixMilli())
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weekly.svg", func(w http.ResponseWriter, r *http.Request) {
		freq := []int{}
		for i := 0; i < 7; i++ {
			freq = append(freq, rand.Int())
		}
		w.Header().Add("Content-Type", "text/html")
		svg := svg.GetWeekSVG(freq)
		svg.WriteTo(w)
	})
	r.HandleFunc("/yearly.svg", func(w http.ResponseWriter, r *http.Request) {
		freq, err := commits.GlobalFrequency(2022, []string{"Groot"})
		if err != nil {
			panic(err)
		}
		svg := svg.GetYearSVG(freq)
		w.Header().Add("Content-Type", "text/html")
		svg.WriteTo(w)
	})

	http.ListenAndServe(":8080", r)
}
