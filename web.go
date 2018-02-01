package main

import (
	"errors"
	"fmt"
	"html/template"
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

var daysHtmlTemplate = `
	<html>
		<head>
			<style>
				table {
					width: 100%;
					border: 1px solid black;
				}
			</style>
		</head>
		<table>
			<tr>
				<th>time</th>
				<th>description</th>
			<tr>
			{{range .}}
				{{range .Tasks}}
					<tr>
						<td>{{.Time}}
						<td>{{.Desc}}
					</tr>
				{{end}}
			{{end}}
		</table>
	</html>
`

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
	tmpl, err := template.New("").Parse(daysHtmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, days); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func web(addr string) error {
	http.HandleFunc("/", webHandler)
	log.Printf("Starting server on: %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
