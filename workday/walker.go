package workday

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type ByStartTime []*Day

func (d ByStartTime) Len() int           { return len(d) }
func (d ByStartTime) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByStartTime) Less(i, j int) bool { return d[i].Start.Before(d[j].Start) }

func DaysForMonth(year int, month time.Month) ([]*Day, error) {
	root := filepath.Join(
		DataDir,
		"tasks",
		strconv.Itoa(year),
		month.String(),
	)
	days := []*Day{}
	if err := filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		day, err := readDayFromFile(fpath)
		if err != nil {
			return err
		}
		days = append(days, day)
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Sort(ByStartTime(days))
	return days, nil
}
