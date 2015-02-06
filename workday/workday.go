package workday

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var DataDir = "data"

func Ping() error {
	t := time.Now()
	d, err := LoadDay(t)
	if err != nil {
		return err
	}
	d.End = t
	return d.Save()
}

type Day struct {
	Start time.Time
	End   time.Time
	Tasks []*Task
}

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

func LoadDay(t time.Time) (*Day, error) {
	fpath := pathForTime(t)
	if exists, err := fileExists(fpath); err != nil {
		return nil, err
	} else if exists {
		return readFromFile(fpath)
	}
	return &Day{t, t, make([]*Task, 0)}, nil
}

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
	if d.End.After(d.Start) {
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
