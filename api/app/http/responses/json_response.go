package responses

import (
	"blog/app/models"
	"blog/app/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type Code int

type JsonResponse struct {
	Success    bool                `json:"success"`
	Code       Code                `json:"code"`
	Data       any                 `json:"data,omitempty"`
	Cursor     *pagination.Cursor  `json:"cursor,omitempty"`
	Total      int64               `json:"total,omitempty"`
	Message    string              `json:"message"`
	Errors     map[string][]string `json:"errors,omitempty"`
	StatusCode int                 `json:"-"`
}

const (
	CodeOK                    Code = 200_001
	CodeBadRequest            Code = 400_001
	CodeUnauthorized          Code = 401_001
	CodeAccessDenied          Code = 403_001
	CodeNotFound              Code = 404_001
	CodeModelNotFound         Code = 404_002
	CodeMethodNotAllowed      Code = 405_001
	CodeParameterError        Code = 422_001
	CodeTooManyAttemptsCode   Code = 429_001
	CodeSystemError           Code = 500_001
	CodeSystemUnavailableCode Code = 500_002
	CodeLoginError            Code = 600_001
)

func Json(c *gin.Context, params ...any) {
	r, prepare := makeJsonResponse(params...)

	if r.Code == 0 {
		r.Code = CodeOK
	}

	if r.StatusCode == 0 {
		r.StatusCode = http.StatusOK
	}

	if r.Message == "" {
		r.Message = http.StatusText(r.StatusCode)
	}

	// 解析 include 参数加载关系
	if v, ok := r.Data.(models.IncludeParser); ok {
		if include := c.Query("include"); include != "" {
			v.ParseInclude(include)
		}
	}

	if prepare != nil {
		prepare(r)
	}

	c.JSON(r.StatusCode, r)
}

func makeJsonResponse(params ...any) (r *JsonResponse, prepare func(r *JsonResponse)) {
	r = &JsonResponse{
		Success: true,
	}

	for i, param := range params {
		switch i {
		case 0:
			switch v := param.(type) {
			case *JsonResponse:
				r = v
			case JsonResponse:
				r = &v
			default:
				r.Data = param
			}
		default:
			switch v := param.(type) {
			case bool:
				r.Success = v
			case *Code:
				r.Code = *v
			case Code:
				r.Code = v
			case int:
				r.StatusCode = v
			case *int64:
				r.Total = *v
			case int64:
				r.Total = v
			case string:
				r.Message = v
			case *pagination.Cursor:
				r.Cursor = v
			case pagination.Cursor:
				r.Cursor = &v
			case map[string][]string:
				r.Errors = v
			case func(r *JsonResponse):
				prepare = v
			}
		}
	}

	if r.Data == nil {
		return
	}

	if v := reflect.ValueOf(r.Data); v.Kind() != reflect.Ptr {
		if v.CanAddr() {
			r.Data = v.Addr().Interface()
		} else {
			// 对于非指针类型的参数，它们没有地址，因此无法进行获取地址的操作。在这种情况下，可以使用 reflect.New() 方法来创建一个新的值并返回其指针。
			ptr := reflect.New(v.Type())
			ptr.Elem().Set(v)
			r.Data = ptr.Interface()
		}
	}

	return
}
