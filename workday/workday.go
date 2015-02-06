package workday

import (
	"encoding/json"
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

// LoadDay reads the Day corresponding to the supplied time.
// If one is not found, a new one is created.
func LoadDay(t time.Time) (*Day, error) {
	fpath := pathForTime(t)
	if exists, err := fileExists(fpath); err != nil {
		return nil, err
	} else if exists {
		return readFromFile(fpath)
	}
	return &Day{t, t, make([]*Task, 0)}, nil
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
		panic("day starts before it ends")
	}
}

func pathForTime(t time.Time) string {
	year, month, day := t.Date()
	return filepath.Join(
		DataDir,
		strconv.Itoa(year),
		month.String(),
		strconv.Itoa(day)+".json",
	)
}

func readFromFile(fpath string) (*Day, error) {
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

func createParentDirs(fpath string) error {
	parent := filepath.Dir(fpath)
	return os.MkdirAll(parent, 0777)
}

func fileExists(fpath string) (bool, error) {
	if _, err := os.Stat(fpath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
