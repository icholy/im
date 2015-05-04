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

	"github.com/icholy/im/workday"
)

var (
	isPing  bool
	isWeb   bool
	isToday bool
)

func init() {
	flag.BoolVar(&isPing, "ping", false, "update workday extent")
	flag.BoolVar(&isWeb, "web", false, "start web server")
	flag.BoolVar(&isToday, "today", false, "show tasks from today")

	workday.DataDir = filepath.Join(os.Getenv("HOME"), ".im")
}

func ping() error {
	if err := workday.LockDataDir(); err != nil {
		log.Fatal(err)
	}
	defer workday.UnlockDataDir()
	return workday.Ping()
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

	flag.Parse()

	if isPing {
		handleErr(ping())
		return
	}

	if isToday {
		handleErr(today())
		return
	}

	if isWeb {
		handleErr(web())
		return
	}

	handleErr(addTask())
}
