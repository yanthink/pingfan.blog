package storage

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type LocalAdapter struct {
	Root    string `json:"root"`
	BaseUrl string `json:"baseUrl"`
}

func NewLocalAdapter(root, baseUrl string) *LocalAdapter {
	return &LocalAdapter{
		Root:    root,
		BaseUrl: baseUrl,
	}
}

func (a *LocalAdapter) fullPath(path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(a.Root, "/"), strings.TrimPrefix(path, "/"))
}

func (a *LocalAdapter) Write(path string, reader io.Reader) error {
	file, err := os.Create(a.fullPath(path))
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, reader)

	return err
}

func (a *LocalAdapter) Read(path string) (*bytes.Buffer, error) {
	file, err := os.Open(a.fullPath(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)

	return buf, err
}

func (a *LocalAdapter) Exists(path string) (bool, error) {
	_, err := os.Stat(a.fullPath(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (a *LocalAdapter) FileExists(path string) (bool, error) {
	file, err := os.Stat(a.fullPath(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return !file.IsDir(), nil
}

func (a *LocalAdapter) DirExists(path string) (bool, error) {
	file, err := os.Stat(a.fullPath(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return file.IsDir(), nil
}

func (a *LocalAdapter) Del(path string) error {
	return os.Remove(a.fullPath(path))
}

func (a *LocalAdapter) DelDir(path string) error {
	return os.RemoveAll(a.fullPath(path))
}

func (a *LocalAdapter) Mkdir(path string, perm os.FileMode) error {
	return os.MkdirAll(a.fullPath(path), perm)
}

func (a *LocalAdapter) Move(src, dst string) error {
	if exists, _ := a.FileExists(src); exists {
		err := a.Mkdir(filepath.Dir(dst), 0750)
		if err != nil {
			return err
		}
	}

	src = a.fullPath(src)
	dst = a.fullPath(dst)

	return os.Rename(src, dst)
}

func (a *LocalAdapter) Copy(src, dst string) error {
	if exists, _ := a.FileExists(src); exists {
		err := a.Mkdir(filepath.Dir(dst), 0750)
		if err != nil {
			return err
		}
	}

	srcFile, err := os.Open(a.fullPath(src))
	if err != nil {
		return err
	}
	defer srcFile.Close()

	return a.Write(dst, srcFile)
}

func (a *LocalAdapter) Append(path string, reader io.Reader) error {
	file, err := os.OpenFile(a.fullPath(path), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0750)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)

	return err
}

func (a *LocalAdapter) Url(path string) string {
	parsed, err := url.Parse(path)
	if err != nil || parsed.Host != "" || path == "" {
		return path
	}

	return fmt.Sprintf("%s/%s", strings.TrimSuffix(a.BaseUrl, "/"), strings.TrimPrefix(path, "/"))
}

func (a *LocalAdapter) ParsePath(path string) string {
	parsedBaseUrl, _ := url.Parse(a.BaseUrl)
	parsedPath, _ := url.Parse(path)

	parsedPath.Path = strings.TrimPrefix(parsedPath.Path, "/")
	parsedBaseUrl.Path = strings.TrimPrefix(parsedBaseUrl.Path, "/")

	if (parsedPath.Host == "" || parsedPath.Host == parsedBaseUrl.Host) && strings.HasPrefix(parsedPath.Path, parsedBaseUrl.Path) {
		return strings.TrimPrefix(strings.TrimPrefix(parsedPath.Path, parsedBaseUrl.Path), "/")
	}

	return path
}
