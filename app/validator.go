package app

import (
	"github.com/go-playground/validator/v10"
)

func GetError(data interface{}) (errBool bool, errMsg string) {
	validate := validator.New()
	errValidate := validate.Struct(data)
	if errValidate != nil {
		for _, err := range errValidate.(validator.ValidationErrors) {
			if err != nil {
				switch err.Tag() {
				case "required":
					return false, err.Field() + " 是必须的"
				case "len":
					return false, err.Field() + " 长度为 " + err.Param()
				case "max":
					return false, err.Field() + " 最大长度长度为 " + err.Param()
				case "min":
					return false, err.Field() + "最小长度为 " + err.Param()
				case "gt":
					return false, err.Field() + " 必须要大于 " + err.Param()
				case "eq":
					return false, err.Field() + " 必须要等于 " + err.Param()
				case "gte":
					return false, err.Field() + " 必须要大于等于 " + err.Param()
				case "lt":
					return false, err.Field() + " 必须要小于 " + err.Param()
				case "lte":
					return false, err.Field() + " 必须要小于等于 " + err.Param()
				}
			}
		}
	}
	return true, ""
}
