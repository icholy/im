package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"./workday"
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

func webHandler(w http.ResponseWriter, r *http.Request) {
	year, month, err := pathToYearMonth(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	return http.ListenAndServe(":8080", nil)
}
