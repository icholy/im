package main

import (
	"flag"
	"log"
	"strings"

	"./workday"
)

var (
	defaultRootDir = "data"
	isPing         bool
)

func init() {
	flag.BoolVar(&isPing, "ping", false, "notify that the workday is still active")
	flag.StringVar(&workday.DataDir, "dir", defaultRootDir, "data directory")
}

func main() {

	flag.Parse()

	if isPing {
		if err := workday.Ping(); err != nil {
			log.Fatal(err)
		}
	} else {
		desc := strings.Join(flag.Args(), " ")
		if err := workday.AddTask(desc); err != nil {
			log.Fatal(err)
		}
	}
}
