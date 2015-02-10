package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"./workday"
)

var (
	isPing bool
	isWeb  bool
)

func init() {
	flag.BoolVar(&isPing, "ping", false, "update workday extent")
	flag.BoolVar(&isWeb, "web", false, "start web server")

	workday.DataDir = filepath.Join(os.Getenv("HOME"), ".im")
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

	if isWeb {
		handleErr(web())
		return
	}

	handleErr(addTask())
}
