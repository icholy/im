package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/icholy/im/workday"
)

func pathToYearMonth(p string) (int, time.Month, error) {
	r := regexp.MustCompile(`/(\d+)/(\d+)`)
	matches := r.FindStringSubmatch(p)
	if matches == nil {
		return 0, 0, errors.New("failed to parse url")
	}
	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, err
	}
	month, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, err
	}
	return year, time.Month(month), nil
}

func redirectToNow(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	nowPath := fmt.Sprintf("/%d/%d", now.Year(), now.Month())
	http.Redirect(w, r, nowPath, http.StatusSeeOther)
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	year, month, err := pathToYearMonth(r.URL.Path)
	if err != nil {
		redirectToNow(w, r)
		return
	}
	days, err := workday.DaysForMonth(year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, d := range days {
		fmt.Fprintln(w, d.String())
	}
}

func web() error {
	http.HandleFunc("/", webHandler)
	log.Println("Starting server on port 8080")
	return http.ListenAndServe(":8080", nil)
}
