package workday

import (
	"os"
	"path/filepath"
)

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

func createFileIfNotExists(fpath string) error {
	if exists, err := fileExists(fpath); err != nil {
		return err
	} else if !exists {
		f, err := os.Create(fpath)
		if err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}
