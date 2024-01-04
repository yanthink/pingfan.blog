package validation

import (
	"blog/app/rules"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

var (
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	validate *validator.Validate
)

func init() {
	locale := zh.New()
	uni = ut.New(locale, locale)
	trans, _ = uni.GetTranslator("zh")
	// 获取gin的校验器
	validate = binding.Validator.Engine().(*validator.Validate)
	// 注册翻译器
	if err := zhTrans.RegisterDefaultTranslations(validate, trans); err != nil {
		return
	}
	// 注册一个函数，获取 struct tag 里自定义的 label 作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// 获取结构体字段的名称
		name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]

		switch name {
		case "-":
			return ""
		case "":
			return fld.Tag.Get("form")
		default:
			return name
		}
	})

	_ = validate.RegisterValidation("username", rules.ValidateUsername)
	_ = validate.RegisterValidation("phone", rules.ValidatePhone)
	_ = validate.RegisterValidation("captcha", rules.ValidateCaptcha)
}

func Validate(c *gin.Context, data any) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case *runtime.TypeAssertionError:
				panic(&Error{Message: "请求参数格式错误"})
			default:
				panic(err)
			}
		}
	}()

	var err error

	// json 是通过 c.Request.Body 方法绑定数据，多次绑定需要使用 ShouldBindBodyWith 方法
	if c.Request.Method != http.MethodGet && c.ContentType() == binding.MIMEJSON {
		err = c.ShouldBindBodyWith(data, binding.JSON)
	} else {
		err = c.ShouldBind(data)
	}

	if err != nil {
		message, errors := Translate(err, data)
		panic(&Error{
			Message: message,
			Errors:  errors,
			Err:     err,
		})
	}
}

// Translate 翻译错误信息
func Translate(err error, data ...any) (message string, errors map[string][]string) {
	message = ""
	errors = map[string][]string{}

	var t reflect.Type

	if len(data) == 1 {
		t = reflect.TypeOf(data[0]).Elem()
	}

	for i, err := range err.(validator.ValidationErrors) {
		var msg string
		var fieldName string

		if t != nil {
			if field, ok := t.FieldByName(err.StructField()); ok {
				msg = field.Tag.Get("message")
				fieldName = field.Tag.Get("form")
			}
		}

		if msg == "" {
			msg = err.Translate(trans)
		}

		if fieldName == "" {
			fieldName = err.Field()
		}

		if i == 0 {
			message = msg
		}

		errors[fieldName] = append(errors[fieldName], msg)
	}

	return
}
