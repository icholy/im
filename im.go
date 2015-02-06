package main

import (
	"flag"
	"log"
	"os"
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

func main() {

	flag.Parse()

	if isPing {
		if err := workday.Ping(); err != nil {
			log.Fatal(err)
		}
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}
	desc := strings.Join(args, " ")
	if err := workday.AddTask(desc); err != nil {
		log.Fatal(err)
	}
}
