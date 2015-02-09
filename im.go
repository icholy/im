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

var isPing bool

func init() {
	flag.BoolVar(&isPing, "ping", false,
		"notify that the workday is still active",
	)
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
		return "", errors.New("EDITOR environtment variable not set")
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
	return string(data), nil
}

func addTask() error {
	desc, err := getDescription()
	if err != nil {
		return err
	}
	if strings.TrimSpace(desc) == "" {
		return errors.New("description cannot be blank")
	}
	if err := workday.LockDataDir(); err != nil {
		return err
	}
	defer workday.UnlockDataDir()
	return workday.AddTask(desc)
}

func main() {

	flag.Parse()

	if isPing {
		if err := ping(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := addTask(); err != nil {
		log.Fatal(err)
	}
}
