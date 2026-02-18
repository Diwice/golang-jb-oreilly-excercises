package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// 1
func formatTime(t time.Time) []byte {
	newTime := t.Format(time.RFC3339)
	return []byte(newTime)
}

// 2
func middlewareLog(l *slog.Logger, h http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		l.Info("User IP", ip)
		h.ServeHTTP(w, r)
	}
}

// 3
type jsonTime struct {
	WeekDay  string `json:"day_of_week"`
	MonthDay int    `json:"day_of_month"`
	Month    string `json:"month"`
	Year     int    `json:"year"`
	Hour     int    `json:"hour"`
	Minute   int    `json:"minute"`
	Second   int    `json:"second"`
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	var formTime []byte
	if val := r.Header.Get("Accept"); val == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		formTime = formatJSONTime(time.Now())
	} else {
		w.Header().Set("Content-Type", "text/plain")
		formTime = formatTime(time.Now())
	}

	w.Write(formTime)
}

func formatJSONTime(t time.Time) []byte {
	weekday := t.Weekday()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()

	timeObj := jsonTime{
		WeekDay:  weekday.String(),
		MonthDay: day,
		Month:    month.String(),
		Year:     year,
		Hour:     hour,
		Minute:   minute,
		Second:   second,
	}
	// Ignoring error, since I can't pass logger up here w/o context
	b, _ := json.Marshal(timeObj)

	return b
}

func main() { // shared
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	newMux := http.NewServeMux()
	newMux.HandleFunc("GET /", middlewareLog(logger, http.HandlerFunc(timeHandler)))
	http.ListenAndServe("localhost:8080", newMux)
}
