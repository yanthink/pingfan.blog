package services

import (
	"blog/app"
	"blog/app/helpers"
	"blog/app/resource"
	"blog/app/storage"
	"blog/config"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"path/filepath"
)

type resourceService struct {
}

var resourceSnowflakeNode, _ = snowflake.NewNode(config.Snowflake.Node)

func (s *resourceService) Upload(header *multipart.FileHeader, rType resource.Type) string {
	file, err := header.Open()
	abortIf(err != nil, "无法打开上传文件")
	defer file.Close()

	disk := storage.Disk()
	var path string

	resize := rType.Resize()
	if resize[0] > 0 && resize[1] > 0 {
		image, err := imaging.Decode(file)
		if err != nil {
			panic(err)
		}

		tmpFile, _ := os.CreateTemp("", "resized_*.png")
		defer os.Remove(tmpFile.Name())

		app.Logger.Sugar().Debug(tmpFile.Name())

		newImage := imaging.Resize(image, resize[0], resize[1], imaging.Lanczos)
		err = imaging.Encode(tmpFile, newImage, imaging.PNG)
		abortIf(err != nil, fmt.Sprintf("%v", err))

		// 将文件指针移动到文件开头，disk.Put 才能读取到内容
		tmpFile.Seek(0, 0)

		path = fmt.Sprintf("%s/%s%s", rType.UploadPath(), resourceSnowflakeNode.Generate(), ".png")
		err = disk.Put(path, tmpFile)
		abortIf(err != nil, fmt.Sprintf("%v", err))
	} else {
		path = fmt.Sprintf("%s/%s%s", rType.UploadPath(), resourceSnowflakeNode.Generate(), filepath.Ext(header.Filename))
		err = disk.Put(path, file)
		abortIf(err != nil, fmt.Sprintf("%v", err))
	}

	return disk.Url(path)
}

// CopyToStorePath 将图片复制到存储目录
func (s *resourceService) CopyToStorePath(path string, rType resource.Type) (url string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			app.Logger.Debug(fmt.Sprintf("%v", r), zap.String("path", path))

			url = path
			ok = false
		}
	}()

	disk := storage.Disk()

	src := disk.ParsePath(path)
	if !rType.IsUploadPath(src) {
		panic("path 不是上传路径")
	}
	if exists, _ := disk.FileExists(src); !exists {
		panic("path 文件不存在")
	}

	dst := fmt.Sprintf("%s/%s%s", rType.StorePath(), resourceSnowflakeNode.Generate(), filepath.Ext(src))

	if config.Storage.ResourceStoreDisk != config.Storage.Disk {
		storeDisk := storage.Disk(config.Storage.ResourceStoreDisk)

		buf, err := disk.Read(src)
		if err != nil {
			panic(err)
		}

		err = storeDisk.Put(dst, buf)
		if err != nil {
			panic(err)
		}

		return storeDisk.Url(dst), true
	}

	if err := disk.Copy(src, dst); err != nil {
		return path, false
	}

	return disk.Url(dst), true
}

// Sync 差异比较，将新增的图片复制到存储目录，如果 shouldRemove = true，则删除 oldUrls 存在而 newUrls 不存在的文件 （被删除的文件）
func (s *resourceService) Sync(newUrls []string, oldUrls []string, rType resource.Type, shouldRemove ...bool) []string {
	disk := storage.Disk()

	oldUrlsMap := helpers.KeyBy(oldUrls, func(url string) string {
		return disk.ParsePath(url)
	})

	newUrls = helpers.Map(helpers.Unique(newUrls), func(_ int, url string) string {
		path := disk.ParsePath(url)

		if _, ok := oldUrlsMap[path]; ok {
			delete(oldUrlsMap, path)
			return disk.Url(path)
		}

		if newUrl, ok := s.CopyToStorePath(path, rType); ok {
			return newUrl
		}

		return url
	})

	if len(shouldRemove) > 0 && shouldRemove[0] {
		storeDisk := storage.Disk(config.Storage.ResourceStoreDisk)

		for path, _ := range oldUrlsMap {
			_ = storeDisk.Del(storeDisk.ParsePath(path))
		}
	}

	return newUrls
}
