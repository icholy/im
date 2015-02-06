package workday

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/nightlyone/lockfile"
)

var (
	lfile          lockfile.Lockfile
	ErrLockTimeout = errors.New("lock timeout")
)

func tryLockFor(l lockfile.Lockfile, d time.Duration) error {
	start := time.Now()
	for time.Now().Sub(start) < d {
		if err := l.TryLock(); err == nil {
			return nil
		}
		time.Sleep(1)
	}
	return ErrLockTimeout
}

// LockDataDir aquires a lock on the DataDir
// It try for 1 second before giving up.
func LockDataDir() error {
	var (
		err   error
		lpath = filepath.Join(DataDir, "lock")
	)
	if err := createFileIfNotExists(lpath); err != nil {
		return err
	}
	lfile, err = lockfile.New(lpath)
	if err != nil {
		return err
	}
	return tryLockFor(lfile, time.Second)
}

// UnlockDataDir releases the lock on the DataDir
func UnlockDataDir() error {
	return lfile.Unlock()
}
