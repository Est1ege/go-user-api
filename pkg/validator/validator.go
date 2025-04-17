package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Настройка кастомных валидаторов
func SetupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Регистрация кастомных валидаторов
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		
		// Добавление кастомных валидаций при необходимости
		// например: v.RegisterValidation("custom_validation", customValidationFunc)
	}
}

// ValidateStruct валидирует структуру и возвращает ошибки в удобном формате
func ValidateStruct(obj interface{}) map[string]string {
	errors := make(map[string]string)
	
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.Struct(obj)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				field := err.Field()
				tag := err.Tag()
				errors[field] = fmt.Sprintf("Поле не соответствует правилу: %s", tag)
			}
		}
	}
	
	return errors
}