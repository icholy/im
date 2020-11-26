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

var funcMap = template.FuncMap{
	"fmtTask": func(t *workday.Task) string {
		return t.Desc
	},
	"fmtDay": func(d *workday.Day) string {
		layout := "Mon Jan 2 2006"
		return fmt.Sprintf("%s - (%s)", d.Start.Format(layout), d.End.Sub(d.Start))
	},
}

var daysHtmlTemplate = `
	<html>
		<h2>Weeks</h2>
		<ol>
		{{range .Weeks}}
			<li>{{.}}</li>
		{{end}}
		</ol>
		<h2>Days</h2>
		{{range .Days}}
			<h3>{{fmtDay .}}</h3>
			<ul>
			{{range .Tasks}}
				<li>{{fmtTask .}}</li>
			{{end}}
			</ul>
		{{end}}
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
	tmpl, err := template.New("").Funcs(funcMap).Parse(daysHtmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, struct {
		Weeks []time.Duration
		Days  []*workday.Day
	}{
		Days:  days,
		Weeks: workday.WeekTotals(days),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func web(addr string) error {
	http.HandleFunc("/", webHandler)
	log.Printf("Starting server on: %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
