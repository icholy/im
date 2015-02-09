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

	"./workday"
	"github.com/jinzhu/now"
)

var (
	isPing    bool
	queryDate string
)

func init() {
	flag.BoolVar(&isPing, "ping", false,
		"notify that the workday is still active",
	)
	flag.StringVar(&queryDate, "query", "", "query date")
	workday.DataDir = filepath.Join(
		os.Getenv("HOME"), ".im",
	)
}

func ping() error {
	if err := workday.LockDataDir(); err != nil {
		log.Fatal(err)
	}
	defer workday.UnlockDataDir()
	return workday.Ping()
}

func getDescription() (string, error) {
	if flag.NArg() > 0 {
		return strings.Join(flag.Args(), " "), nil
	}
	f, err := ioutil.TempFile(os.TempDir(), "im.")
	if err != nil {
		return "", err
	}
	defer f.Close()
	defer os.Remove(f.Name())
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
	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
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

func query(t time.Time) error {
	w := workday.NewWalker(
		now.New(t).BeginningOfMonth(),
		now.New(t).EndOfMonth(),
	)
	for d := range w.OutCh() {
		fmt.Println(d.String())
	}
	return w.Err()
}

func main() {

	flag.Parse()

	if isPing {
		if err := ping(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if queryDate != "" {
		t, err := now.Parse(queryDate)
		if err != nil {
			log.Fatal(err)
		}
		if err := query(t); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := addTask(); err != nil {
		log.Fatal(err)
	}
}
