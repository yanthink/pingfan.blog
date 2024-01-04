package resource

import (
	"blog/app/helpers"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/duke-git/lancet/v2/formatter"
	"mime/multipart"
	"regexp"
	"strings"
)

type Dimensions struct {
	accept    string
	width     int
	height    int
	minWidth  int
	minHeight int
	maxWidth  int
	maxHeight int
	ratio     float64
	minSize   int64
	maxSize   int64
}

func NewDimensions() *Dimensions {
	return &Dimensions{}
}

func (d *Dimensions) Accept(accept string) *Dimensions {
	d.accept += "," + accept

	return d
}

func (d *Dimensions) Width(width int) *Dimensions {
	d.width = width

	return d
}

func (d *Dimensions) Height(height int) *Dimensions {
	d.height = height

	return d
}

func (d *Dimensions) MinWidth(minWidth int) *Dimensions {
	d.minWidth = minWidth

	return d
}

func (d *Dimensions) MinHeight(minHeight int) *Dimensions {
	d.minHeight = minHeight

	return d
}

func (d *Dimensions) MaxWidth(maxWidth int) *Dimensions {
	d.maxWidth = maxWidth

	return d
}

func (d *Dimensions) MaxHeight(maxHeight int) *Dimensions {
	d.maxWidth = maxHeight

	return d
}

func (d *Dimensions) Ratio(ratio float64) *Dimensions {
	d.ratio = ratio

	return d
}

func (d *Dimensions) MinSize(minSize int64) *Dimensions {
	d.minSize = minSize

	return d
}

func (d *Dimensions) MaxSize(maxSize int64) *Dimensions {
	d.maxSize = maxSize

	return d
}

func (d *Dimensions) Validate(file *multipart.FileHeader) {
	var errs []string

	contentType := file.Header.Get("Content-Type")

	ok := d.accept == "" || helpers.ContainsBy(strings.Split(d.accept, ","), func(item string) bool {
		item = strings.TrimSpace(item)
		if item == "" {
			return false
		}

		regexPattern := fmt.Sprintf("(?i)%s", strings.ReplaceAll(item, "*", ".*"))

		if match, _ := regexp.MatchString(regexPattern, contentType); match {
			return true
		}

		return false
	})

	if !ok {
		message := "图片格式不正确"
		panic(&Error{
			Message: message,
			Errors: map[string][]string{
				"file": {message},
			},
			Err: fmt.Errorf(message),
		})
	}

	image, err := file.Open()
	if err != nil {
		message := "无法打开图片文件"
		panic(&Error{
			Message: message,
			Errors: map[string][]string{
				"file": {message},
			},
			Err: fmt.Errorf(message),
		})
	}
	defer image.Close()

	decodedImage, err := imaging.Decode(image)
	if err != nil {
		message := "图片格式不正确"
		panic(&Error{
			Message: message,
			Errors: map[string][]string{
				"file": {message},
			},
			Err: fmt.Errorf(message),
		})
	}

	width := decodedImage.Bounds().Dx()
	height := decodedImage.Bounds().Dy()

	if d.width > 0 && width != d.width {
		errs = append(errs, fmt.Sprintf("图像宽度必须等于%d像素", d.width))
	}

	if d.height > 0 && height != d.height {
		errs = append(errs, fmt.Sprintf("图像高度必须等于%d像素", d.height))
	}

	if d.minWidth > 0 && width < d.minWidth {
		errs = append(errs, fmt.Sprintf("图像宽度不能小于%d像素", d.minWidth))
	}

	if d.minHeight > 0 && height < d.minHeight {
		errs = append(errs, fmt.Sprintf("图像高度不能小于%d像素", d.minWidth))
	}

	if d.maxWidth > 0 && width > d.maxWidth {
		errs = append(errs, fmt.Sprintf("图像宽度不能大于%d像素", d.maxWidth))
	}

	if d.maxHeight > 0 && height > d.maxHeight {
		errs = append(errs, fmt.Sprintf("图像高度不能大于%d像素", d.maxHeight))
	}

	if d.ratio > 0 && float64(width)/float64(height) != d.ratio {
		errs = append(errs, fmt.Sprintf("图像宽高比必须等于%.1f", d.ratio))
	}

	if d.minSize > 0 && file.Size < d.minSize {
		errs = append(errs, fmt.Sprintf("图像大小不能小于%s", formatter.BinaryBytes(float64(d.minSize))))
	}

	if d.maxSize > 0 && file.Size > d.maxSize {
		errs = append(errs, fmt.Sprintf("图像大小不能大于%s", formatter.BinaryBytes(float64(d.maxSize))))
	}

	if len(errs) > 0 {
		panic(&Error{
			Message: errs[0],
			Errors: map[string][]string{
				"file": errs,
			},
			Err: fmt.Errorf(errs[0]),
		})
	}
}
