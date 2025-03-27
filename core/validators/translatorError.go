package validators

import (
	"fmt"
	"regexp"
	"strings"
)

func TranslatorError(err error) string {
	reg := regexp.MustCompile(`"(.*)?"`)
	match := reg.FindString(err.Error())
	if strings.Contains(err.Error(), "is not set") {
		if match != "" {
			return fmt.Sprintf("%s 为必填字段xx!", strings.Replace(match, `"`, "", -1))
		}
	} else if strings.Contains(err.Error(), "mismatch for field") {
		if match != "" {
			return fmt.Sprintf("%s 数据类型不对!", strings.Replace(match, `"`, "", -1))
		}
	}
	return ""
}

//func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
//	return func(ctx context.Context, err error) (int, any) {
//		fmt.Println(err, "错误信息")
//		causeErr := errors.Cause(err)
//		fmt.Println(causeErr, "111-->")
//		var errMessage = err.Error()
//		// 翻译错误
//		translatorError := TranslatorError(err)
//		if translatorError != "" {
//			errMessage = translatorError
//		}
//		// 日志记录
//		logx.WithContext(ctx).Errorf("【%s】 err %v", name, errMessage)
//		return http.StatusBadRequest, errMessage
//	}
//}
