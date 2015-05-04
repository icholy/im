package workday

import "time"

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
