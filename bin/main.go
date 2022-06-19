package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/taigrr/gitgraph/graph"
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
		svg := graph.GetWeekSVG(freq)
		svg.WriteTo(w)

	})
	r.HandleFunc("/yearly.svg", func(w http.ResponseWriter, r *http.Request) {
		freq := []int{}
		for i := 0; i < 365; i++ {
			freq = append(freq, rand.Int())
		}
		svg := graph.GetYearSVG(freq)
		w.Header().Add("Content-Type", "text/html")
		svg.WriteTo(w)

	})

	http.ListenAndServe("0.0.0.0:5578", r)
}
