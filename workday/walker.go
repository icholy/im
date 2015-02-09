package workday

import (
	"sync"
	"time"
)

type Walker struct {
	current time.Time
	end     time.Time
	out     chan *Day
	err     error
	m       sync.Mutex
}

func NewWalker(start, end time.Time) *Walker {
	w := &Walker{
		current: start,
		end:     end,
		out:     make(chan *Day),
	}
	go func() {
		w.m.Lock()
		defer w.m.Unlock()
		w.err = w.loop()
	}()
	return w
}

func (w *Walker) loop() error {
	defer close(w.out)
	for w.current.Before(w.end) {
		fpath := pathForTime(w.current)
		exists, err := fileExists(fpath)
		if err != nil {
			return err
		}
		if exists {
			day, err := readDayFromFile(fpath)
			if err != nil {
				return err
			}
			w.out <- day
		}
		w.current = w.current.Add(24 * time.Hour)
	}
	return nil
}

func (w *Walker) OutCh() chan *Day {
	return w.out
}

func (w *Walker) Err() error {
	w.m.Lock()
	defer w.m.Unlock()
	return w.err
}
