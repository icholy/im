package workday

import (
	"errors"
	"time"
)

type Task struct {
	Time time.Time
	Desc string
}

func AddTask(desc string) error {
	t := &Task{
		Time: time.Now(),
		Desc: desc,
	}
	d, err := LoadDay(t.Time)
	if err != nil {
		return err
	}
	d.Tasks = append(d.Tasks, t)
	d.End = t.Time
	return d.Save()
}

func Undo() error {
	d, err := LoadDay(time.Now())
	if err != nil {
		return err
	}
	if len(d.Tasks) == 0 {
		return errors.New("there are no tasks for today")
	}
	d.Tasks = d.Tasks[:len(d.Tasks)-1]
	return d.Save()
}
