package workday

import (
	"regexp"
	"time"
)

type Task struct {
	Time time.Time
	Desc string
	Tags []string
}

func findTags(s string) []string {
	r := regexp.MustCompile(`@\w+`)
	tags := r.FindAllString(s, -1)
	for i, tag := range tags {
		tags[i] = tag[1:]
	}
	return tags
}

func AddTask(desc string) error {
	t := &Task{
		Time: time.Now(),
		Desc: desc,
		Tags: findTags(desc),
	}
	d, err := LoadDay(t.Time)
	if err != nil {
		return err
	}
	d.Tasks = append(d.Tasks, t)
	return d.Save()
}
