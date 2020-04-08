package jira

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/andygrunwald/go-jira"
)

var BaseURL = "https://jira.accipiterradar.com:8443"

type Issue struct {
	Name    string
	Summary string
}

func InProgress(username, password string, timeout time.Duration) ([]Issue, error) {

	hc := &http.Client{
		Timeout: timeout,
		Transport: &jira.BasicAuthTransport{
			Username: username,
			Password: password,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	client, err := jira.NewClient(hc, BaseURL)
	if err != nil {
		return nil, err
	}

	query := `status = "In Progress" AND assignee in (currentUser())`
	issues, _, err := client.Issue.Search(query, nil)
	if err != nil {
		return nil, err
	}

	var ss []Issue
	for _, issue := range issues {
		ss = append(ss, Issue{Name: issue.Key, Summary: issue.Fields.Summary})
	}
	return ss, nil
}
