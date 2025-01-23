package validators

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

// https://blog.csdn.net/kuangshp128/article/details/141175492

func Validate(dataStruct interface{}) error {
	zhT := zh.New()
	validate := validator.New()
	// 注册一个函数，获取struct tag里自定义的label作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	uni := ut.New(zhT)
	trans, _ := uni.GetTranslator("zh")
	// 注册自定义结构体字段校验方法
	// 1.注册校验器：校验时间要比当前的晚
	if err := validate.RegisterValidation(
		"checkAfterDate",
		ValidateAfterDate,
	); err != nil {
		return err
	}

	if err := validate.RegisterValidation(
		"checkMobile",
		ValidateMobile,
	); err != nil {
		return err
	}

	if err := validate.RegisterValidation(
		"checkEmail",
		ValidateEmail,
	); err != nil {
		return err
	}

	if err := validate.RegisterValidation(
		"checkDate",
		ValidateDate,
	); err != nil {
		return err
	}

	// 2.注册翻译器：校验时间要比当前的晚
	if err := validate.RegisterTranslation(
		"checkAfterDate",
		trans,
		registerTranslator("checkAfterDate", "{0} 必须要晚于当前日期"),
		translate,
	); err != nil {
		return err
	}

	// 验证器注册翻译器
	if err := zhTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		return err
	}
	err := validate.Struct(dataStruct)
	if err != nil {
		for _, err1 := range err.(validator.ValidationErrors) {
			return errors.New(err1.Translate(trans))
		}
	}
	return nil
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}
