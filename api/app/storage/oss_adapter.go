package storage

import (
	"blog/app"
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"net/url"
	"os"
	"strings"
)

type OssAdapter struct {
	*oss.Bucket
	Domain string
	Ssl    bool
}

func NewOssAdapter(accessKeyId, accessKeySecret, endpoint, bucket, domain string, ssl bool) *OssAdapter {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		app.Logger.Error(fmt.Sprintf("%v", err))
		panic(err)
	}

	adapter := OssAdapter{
		Domain: domain,
		Ssl:    ssl,
	}

	adapter.Bucket, err = client.Bucket(bucket)
	if err != nil {
		panic(err)
	}

	return &adapter
}

func (a *OssAdapter) Write(path string, reader io.Reader) error {
	return a.PutObject(path, reader)
}

func (a *OssAdapter) Read(path string) (*bytes.Buffer, error) {
	body, err := a.GetObject(path)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, body)

	return buf, err
}

func (a *OssAdapter) Exists(path string) (bool, error) {
	return a.IsObjectExist(path)
}

func (a *OssAdapter) FileExists(path string) (bool, error) {
	return a.IsObjectExist(path)
}

func (a *OssAdapter) DirExists(path string) (bool, error) {
	return a.IsObjectExist(path, oss.Marker(""))
}

func (a *OssAdapter) Del(path string) error {
	return a.DeleteObject(path)
}

func (a *OssAdapter) DelDir(path string) error {
	marker := oss.Marker("")
	// 填写待删除目录的完整路径，完整路径中不包含Bucket名称。
	prefix := oss.Prefix(strings.TrimSuffix(path, "/") + "/")
	count := 0
	for {
		lor, err := a.ListObjects(marker, prefix)
		if err != nil {
			return err
		}

		var objects []string
		for _, object := range lor.Objects {
			objects = append(objects, object.Key)
		}
		// 删除目录及目录下的所有文件。
		// 将oss.DeleteObjectsQuiet设置为true，表示不返回删除结果。
		delRes, err := a.DeleteObjects(objects, oss.DeleteObjectsQuiet(true))
		if err != nil {
			return err
		}

		if len(delRes.DeletedObjects) > 0 {
			return fmt.Errorf("these objects deleted failure, %v", delRes.DeletedObjects)
		}

		count += len(objects)

		prefix = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}

	return nil
}

func (a *OssAdapter) Mkdir(path string, _ os.FileMode) error {
	return a.PutObject(strings.TrimSuffix(path, "/")+"/", bytes.NewReader([]byte("")))
}

func (a *OssAdapter) Move(src, dst string) error {
	if err := a.Copy(src, dst); err != nil {
		return err
	}

	_ = a.Del(src)

	return nil
}

func (a *OssAdapter) Copy(src, dst string) error {
	_, err := a.CopyObject(src, dst)

	return err
}

func (a *OssAdapter) Append(path string, reader io.Reader) error {
	buf, err := a.Read(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(buf, reader)

	return a.Write(path, buf)
}

func (a *OssAdapter) Url(path string) string {
	parsed, err := url.Parse(path)
	if err != nil || parsed.Host != "" || path == "" {
		return path
	}

	parsed.Scheme = "http"
	if a.Ssl {
		parsed.Scheme = "https"
	}

	parsed.Host = a.Domain

	return parsed.String()
}

func (a *OssAdapter) ParsePath(path string) string {
	parsed, _ := url.Parse(path)

	if parsed.Host == "" || parsed.Host == a.Domain {
		return strings.TrimPrefix(parsed.Path, "/")
	}

	return path
}
