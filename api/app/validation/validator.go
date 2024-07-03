package validation

import (
	"blog/app/rules"
	"errors"
	"fmt"
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

func Validate(c *gin.Context, dest any) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case *runtime.TypeAssertionError:
				panic(&Error{Message: fmt.Sprintf("%v", err)})
			default:
				panic(err)
			}
		}
	}()

	var err error

	// json 是通过 c.Request.Body 方法绑定数据，多次绑定需要使用 ShouldBindBodyWith 方法
	if c.Request.Method != http.MethodGet && c.ContentType() == binding.MIMEJSON {
		err = c.ShouldBindBodyWith(dest, binding.JSON)
	} else {
		err = c.ShouldBind(dest)
	}

	if err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			message, errorsMap := Translate(errs, reflect.TypeOf(dest).Elem())
			panic(&Error{
				Message: message,
				Errors:  errorsMap,
				Err:     err,
			})
		}

		panic(fmt.Sprintf("%v", err))
	}
}

// Translate 翻译错误信息
func Translate(errs validator.ValidationErrors, dv reflect.Type) (message string, errorsMap map[string][]string) {
	errorsMap = map[string][]string{}

	messages := map[string]string{
		"username": "用户名格式错误。",
		"phone":    "手机号码格式错误。",
		"captcha":  "验证码错误。",
	}

	for i, err := range errs {
		var msg string
		var tagName string

		if field, ok := dv.FieldByName(err.StructField()); ok {
			if msg = field.Tag.Get("message"); msg == "" {
				msg = messages[err.Tag()]
			}
			if tagName = field.Tag.Get("json"); tagName == "" {
				tagName = field.Tag.Get("form")
			}
		}

		if msg == "" {
			msg = err.Translate(trans)
		}

		if tagName == "" {
			tagName = err.Field()
		}

		if i == 0 {
			message = msg
		}

		errorsMap[tagName] = append(errorsMap[tagName], msg)
	}

	return
}
