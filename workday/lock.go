package workday

import (
	"errors"
	"os"
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

func createLock(fpath string) (lockfile.Lockfile, error) {
	// create the file if it doesn't exist
	if exists, err := fileExists(fpath); err != nil {
		return "", err
	} else if !exists {
		f, err := os.Create(fpath)
		if err != nil {
			return "", err
		}
		if err := f.Close(); err != nil {
			return "", err
		}
	}
	return lockfile.New(fpath)
}

// LockDataDir aquires a lock on the DataDir
// It try for 1 second before giving up.
func LockDataDir() error {
	var (
		err   error
		lpath = filepath.Join(DataDir, "lock")
	)
	lfile, err = createLock(lpath)
	if err != nil {
		return err
	}
	return tryLockFor(lfile, time.Second)
}

// UnlockDataDir releases the lock on the DataDir
func UnlockDataDir() error {
	return lfile.Unlock()
}
