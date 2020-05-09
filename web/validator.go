// 数据验证插件, 详细的tag请查看
// @see https://godoc.org/gopkg.in/go-playground/validator.v9
// , = 0x2C, | = 0x7C
// 使用方式:
//type Data struct {
//	Value1 `validate:"required"`
//	Value2 `validate:"-"`
//}
// 常用 Tag 如下
//		-				跳过验证
//		|				Or操作符
//		dive			深入子级数据验证
//			[][]string with validation tag "gt=0,dive,len=1,dive,required"
//			gt=0 will be applied to []
//			len=1 will be applied to []string
//			required will be applied to string
//		omitempty		允许空, Usage: omitempty
//		required		必须, Usage: required
//		len				长度, Usage: len=10, 数字验证等于给定的值, 字符串验证长度, 数组验证元素数量
//		max				最大值, Usage: max=10, 同上
//		min				最小值, Usage: min=10, 同上
//		eq				等于, Usage: eq=10
//		ne				不等于, Usage: ne=10
//		oneof			枚举, Usage: oneof=1 2 3
//		gt				大于, Usage: gt=10
//		gte				大于等于, Usage: gte=10
//		lt				小于, Usage: lt=10
//		lte				小于等于, Usage: lte=10
//		eqfield			等于指定字段值, Usage: eqfield=password
//		nefield			不等于指定字段值, Usage: nefield=username
//		unique			For arrays & slices, unique will ensure that there are no duplicates. For maps, unique will ensure that there are no duplicate values.
//		alpha 			仅允许ASCII字符, 即 a-zA-Z
//		alphanum		仅允许ASCII字母和数字, 即 a-zA-Z0-9
//		alphaunicode 	仅允许Unicode字符
//		alphanumunicode	仅允许Unicode字母和数字
//		numeric			仅允许数字类型, 包括浮点和负数
// 		number			仅允许正整数
//		email			Email验证
//		url				URL验证
//		uri				URI验证
//		base64			BASE64验证
//		contains		验证字符串是否包含某个字符串, Usage: contains=@
//		excludes		验证字符串是否不包含某个字符串, Usage: excludes=@
//		ip				IP地址验证
//		ipv4
//		ipv6
package web

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/micro/go-micro/v2/util/log"
	"gopkg.in/go-playground/validator.v9"
	translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"net/http"
	"regexp"
)

type webValidator struct {
	translator ut.Translator
	validator  *validator.Validate
}

func (wv *webValidator) Validate(i interface{}) error {
	if err := wv.validator.Struct(i); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				return echo.NewHTTPError(http.StatusBadRequest, e.Translate(wv.translator))
			}
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	return nil
}

func NewWebValidator() *webValidator {
	valid := validator.New()
	enLocale := en.New()
	zhLocale := zh.New()
	translator := ut.New(enLocale, enLocale, zhLocale)
	trans, _ := translator.GetTranslator("zh")
	// 注册默认转换语言, 中文语言需要自行参考 translations 包实现
	_ = translations.RegisterDefaultTranslations(valid, trans)

	// 自定义验证
	_ = valid.RegisterValidation("mobile", validMobile)
	_ = valid.RegisterTranslation("mobile", trans, regFunc("mobile", "{0}不是有效的手机号", false), tranFunc)

	return &webValidator{validator: valid, translator: trans}
}

func regFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

func tranFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Warnf("error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}
	return t
}

// ----- 自定义验证

// validMobile 验证手机号
func validMobile(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(`^1\d{10}$`)
	return reg.MatchString(fl.Field().String())
}
