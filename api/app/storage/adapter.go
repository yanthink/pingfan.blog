package storage

import (
	"bytes"
	"io"
	"os"
)

type Adapter interface {
	Write(path string, reader io.Reader) error
	Read(path string) (*bytes.Buffer, error)
	Exists(path string) (bool, error)
	FileExists(path string) (bool, error)
	DirExists(path string) (bool, error)
	Del(path string) error
	DelDir(path string) error
	Mkdir(path string, perm os.FileMode) error
	Move(src, dst string) error
	Copy(src, dst string) error
	Append(path string, reader io.Reader) error
	Url(path string) string
	ParsePath(path string) string
}
