package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/icholy/im/jira"
	"github.com/icholy/im/workday"
)

var (
	isPing   bool
	isWeb    bool
	isToday  bool
	isMonth  bool
	isUndo   bool
	isTest   time.Duration
	webAddr  string
	jiraUser string
	jiraPass string
	timeout  time.Duration
)

func init() {
	flag.BoolVar(&isPing, "ping", false, "update workday extent")
	flag.DurationVar(&timeout, "ping.timeout", 10*time.Second, "ping timeout")
	flag.BoolVar(&isWeb, "web", false, "start web server")
	flag.BoolVar(&isToday, "today", false, "show tasks from today")
	flag.BoolVar(&isMonth, "month", false, "show tasks from month")
	flag.BoolVar(&isUndo, "undo", false, "undo last task for today")
	flag.DurationVar(&isTest, "test", 0, "exit with 0 if there are tasks")
	flag.StringVar(&webAddr, "web.addr", ":8081", "web address to listen on")
	flag.StringVar(&jiraUser, "jira.username", os.Getenv("JIRA_USERNAME"), "jira username")
	flag.StringVar(&jiraPass, "jira.password", os.Getenv("JIRA_API_TOKEN"), "jira password or token")
	flag.StringVar(&jira.BaseURL, "jira.url", os.Getenv("JIRA_BASE_URL"), "jira base url")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	workday.DataDir = filepath.Join(home, ".im")

	flag.Parse()
}

func ping() error {
	if err := workday.LockDataDir(); err != nil {
		log.Fatal(err)
	}
	defer workday.UnlockDataDir()
	defer workday.Ping()

	if jiraUser == "" && jiraPass == "" {
		return nil
	}
	issues, err := jira.InProgress(jiraUser, jiraPass, timeout)
	if err != nil {
		return err
	}
	for _, issue := range issues {
		desc := fmt.Sprintf("%s: %s", issue.Name, issue.Summary)
		if err := workday.AddTask(desc); err != nil {
			return err
		}
	}
	return nil
}

func today() error {
	if err := workday.LockDataDir(); err != nil {
		log.Fatal(err)
	}
	defer workday.UnlockDataDir()
	day, err := workday.LoadDay(time.Now())
	if err != nil {
		return err
	}
	fmt.Println(day.String())
	return nil
}

func month() error {
	if err := workday.LockDataDir(); err != nil {
		return err
	}
	defer workday.UnlockDataDir()
	now := time.Now()
	days, err := workday.DaysForMonth(now.Year(), now.Month())
	if err != nil {
		return err
	}
	for _, d := range days {
		fmt.Println(d.String())
	}
	return nil
}

func getDescription() (string, error) {
	// use args as message if there are any
	if flag.NArg() > 0 {
		return strings.Join(flag.Args(), " "), nil
	}
	// create temp file
	f, err := ioutil.TempFile(os.TempDir(), "im.")
	if err != nil {
		return "", err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	// open file in editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "", errors.New("EDITOR environment variable not set")
	}
	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	// read file contents
	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func test() error {
	now := time.Now()
	day, err := workday.LoadDay(now)
	if err != nil {
		return err
	}
	if len(day.Tasks) == 0 {
		fmt.Println("You don't have any tasks recorded for the day")
		return nil
	}
	var (
		latest  = day.Tasks[len(day.Tasks)-1]
		idleFor = now.Sub(latest.Time)
	)
	if idleFor > isTest {
		fmt.Printf(
			"You haven't recorded any tasks for %s\n",
			time.Duration(idleFor.Seconds())*time.Second, // rounded
		)
	}
	return nil
}

func undo() error {
	return workday.Undo()
}

func addTask() error {
	desc, err := getDescription()
	if err != nil {
		return err
	}
	if desc == "" {
		return errors.New("description cannot be blank")
	}
	if err := workday.LockDataDir(); err != nil {
		return err
	}
	defer workday.UnlockDataDir()
	return workday.AddTask(desc)
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	switch {
	case isPing:
		handleErr(ping())
	case isToday:
		handleErr(today())
	case isMonth:
		handleErr(month())
	case isTest != 0:
		handleErr(test())
	case isUndo:
		handleErr(undo())
	case isWeb:
		handleErr(web(webAddr))
	default:
		handleErr(addTask())
	}
}
