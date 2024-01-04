package storage

import (
	"blog/config"
	"io"
	"path/filepath"
)

type Storage struct {
	Adapter
}

func (f *Storage) Put(path string, reader io.Reader) error {
	if err := f.Adapter.Mkdir(filepath.Dir(path), 0750); err != nil {
		return err
	}

	return f.Adapter.Write(path, reader)
}

var (
	manager = NewManager()
	disk    *Storage
)

func Disk(name ...string) *Storage {
	if len(name) > 0 {
		disk, err := manager.Disk(name[0])
		if err != nil {
			panic(err)
		}

		return disk
	}

	if disk == nil {
		var err error
		if disk, err = manager.Disk(config.Storage.Disk); err != nil {
			panic(err)
		}
	}

	return disk
}
