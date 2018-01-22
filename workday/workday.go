package workday

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// DataDir is where the root directory
// where the data is stored
var DataDir = "data"

// Day is a single workday
type Day struct {
	Start time.Time
	End   time.Time
	Tasks []*Task
}

func pathForTime(t time.Time) string {
	year, month, day := t.Date()
	return filepath.Join(
		DataDir,
		"tasks",
		strconv.Itoa(year),
		month.String(),
		strconv.Itoa(day)+".json",
	)
}

func readDayFromFile(fpath string) (*Day, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := new(Day)
	if err := json.NewDecoder(f).Decode(d); err != nil {
		return nil, err
	}
	return d, nil
}

// LoadDay reads the Day corresponding to the supplied time.
// If one is not found, a new one is created.
func LoadDay(t time.Time) (*Day, error) {
	fpath := pathForTime(t)
	if exists, err := fileExists(fpath); err != nil {
		return nil, err
	} else if exists {
		return readDayFromFile(fpath)
	}
	return &Day{t, t, make([]*Task, 0)}, nil
}

// Save writes the Day to the data directory as json
// and creates any required parent directories
func (d *Day) Save() error {
	d.MustBeSane()
	fpath := pathForTime(d.Start)
	if err := createParentDirs(fpath); err != nil {
		return err
	}
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(d); err != nil {
		return err
	}
	return nil
}

// Ping updates the current Day's End time.
// A Day is created if it does not already exist.
func Ping() error {
	t := time.Now()
	d, err := LoadDay(t)
	if err != nil {
		return err
	}
	d.End = t
	return d.Save()
}

// MustBeSane panics if there are Day's data is inconsistent
func (d *Day) MustBeSane() {
	sYear, sMonth, sDay := d.Start.Date()
	eYear, eMonth, eDay := d.End.Date()
	if sYear != eYear {
		panic("years do not match")
	}
	if sMonth != eMonth {
		panic("months do not match")
	}
	if sDay != eDay {
		panic("days do not match")
	}
	if d.End.Before(d.Start) {
		panic("day ends before it starts")
	}
}

func roundDuration(d time.Duration) time.Duration {
	return time.Duration(d.Seconds()) * time.Second
}

func (d *Day) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf,
		"Date: %s\nTime: %s - %s (%s)\n\n",
		d.Start.Format("Mon Jan _2 2006"),
		d.Start.Format(time.Kitchen),
		d.End.Format(time.Kitchen),
		roundDuration(d.End.Sub(d.Start)),
	)
	for _, t := range d.Tasks {
		fmt.Fprintf(&buf, "\t- %s\n", t.Desc)
	}

	return buf.String()
}
