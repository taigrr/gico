package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/graph/svg"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /weekly.svg", func(w http.ResponseWriter, r *http.Request) {
		author := r.URL.Query().Get("author")
		highlight := r.URL.Query().Get("highlight")
		shouldHighlight := highlight != ""

		repoPaths, err := commits.GetRepos()
		if err != nil {
			http.Error(w, "failed to get repos", http.StatusInternalServerError)
			log.Printf("error getting repos: %v", err)
			return
		}
		week, err := repoPaths.GetWeekFreq([]string{author})
		if err != nil {
			http.Error(w, "failed to get weekly frequency", http.StatusInternalServerError)
			log.Printf("error getting weekly freq: %v", err)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		svgData := svg.GetWeekSVG(week, shouldHighlight)
		svgData.WriteTo(w)
	})

	mux.HandleFunc("GET /stats.json", func(w http.ResponseWriter, r *http.Request) {
		year := time.Now().Year()
		yst := r.URL.Query().Get("year")
		author := r.URL.Query().Get("author")
		if y, err := strconv.Atoi(yst); err == nil {
			year = y
		}
		repoPaths, err := commits.GetRepos()
		if err != nil {
			http.Error(w, "failed to get repos", http.StatusInternalServerError)
			log.Printf("error getting repos: %v", err)
			return
		}
		freq, err := repoPaths.FrequencyChan(year, []string{author})
		if err != nil {
			http.Error(w, "failed to get frequency", http.StatusInternalServerError)
			log.Printf("error getting freq: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(freq); err != nil {
			log.Printf("error encoding response: %v", err)
		}
	})

	mux.HandleFunc("GET /yearly.svg", func(w http.ResponseWriter, r *http.Request) {
		year := time.Now().Year()
		yst := r.URL.Query().Get("year")
		author := r.URL.Query().Get("author")
		highlight := r.URL.Query().Get("highlight")
		shouldHighlight := highlight != ""
		if y, err := strconv.Atoi(yst); err == nil {
			if year != y {
				shouldHighlight = false
			}
			year = y
		}
		repoPaths, err := commits.GetRepos()
		if err != nil {
			http.Error(w, "failed to get repos", http.StatusInternalServerError)
			log.Printf("error getting repos: %v", err)
			return
		}
		freq, err := repoPaths.FrequencyChan(year, []string{author})
		if err != nil {
			http.Error(w, "failed to get frequency", http.StatusInternalServerError)
			log.Printf("error getting freq: %v", err)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		svgData := svg.GetYearSVG(freq, shouldHighlight)
		svgData.WriteTo(w)
	})

	log.Println("gico-server listening on :8822")
	if err := http.ListenAndServe(":8822", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
