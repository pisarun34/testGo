package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

var Validate *validator.Validate
var Translator ut.Translator

func init() {
	fmt.Println("Initializing Translator...")
	// เตรียม translator
	enTranslator := en.New()
	uni := ut.New(enTranslator, enTranslator)
	var found bool
	Translator, found = uni.GetTranslator("en")
	if !found {
		panic("translator not found")
	}

	// สร้าง validator instance
	Validate = validator.New()

	// ลงทะเบียน translations สำหรับภาษาอังกฤษ
	translations.RegisterDefaultTranslations(Validate, Translator)

	Validate.RegisterTranslation("numeric", Translator, func(ut ut.Translator) error {
		return ut.Add("numeric", "{0} must contain only numbers", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("numeric", fe.Field())
		return t
	})

	// ตั้งค่า custom error messages
	Validate.RegisterTranslation("required", Translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	Validate.RegisterTranslation("len", Translator, func(ut ut.Translator) error {
		return ut.Add("len", "{0} must have a length between {1} and {2}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("len", fe.Field(), fe.Param())
		return t
	})

	// กำหนดให้ใช้ชื่อฟิลด์จาก json tag
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
