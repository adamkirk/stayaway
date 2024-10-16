package common

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/adamkirk-stayaway/organisations/internal/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationExtension struct{}

func (ve *ValidationExtension) Translations() []validation.Translation {
	return []validation.Translation{
		{
			Rule: "required",
			RegisterFunc: func(ut ut.Translator) error {
				return ut.Add("required", "is required", true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("required")

				return t
			},
		},
		{
			Rule: "slug",
			RegisterFunc: func(ut ut.Translator) error {
				return ut.Add("slug", "must contain only alphanumeric and hyphen characters; cannot start with a hyphen", true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("slug")

				return t
			},
		},
		{
			Rule: "postcode",
			RegisterFunc: func(ut ut.Translator) error {
				// TODO give a better message here, need to account for different options
				return ut.Add("postcode", "must be a valid postcode", true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("postcode")

				return t
			},
		},
		{
			Rule: "min",
			RegisterFunc: func(ut ut.Translator) error {
				return ut.Add("min", "{0}", true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {

				minValue := fe.Param()
				var msg string

				k := fe.Type().Kind()

				if fe.Type().Kind() == reflect.Pointer {
					k = fe.Type().Elem().Kind()
				}

				switch k {
				case reflect.Array, reflect.Slice:
					msg = "must contain more than %s items"
				case
					reflect.Float32,
					reflect.Float64,
					reflect.Int,
					reflect.Int8,
					reflect.Int16,
					reflect.Int32,
					reflect.Int64,
					reflect.Uint,
					reflect.Uint8,
					reflect.Uint16,
					reflect.Uint32,
					reflect.Uint64:
					msg = "must be larger than %s"
				case reflect.String:
					msg = "must be more than %s characters long"
				}

				t, _ := ut.T("min", fmt.Sprintf(msg, minValue))

				return t
			},
		},
		{
			Rule: "max",
			RegisterFunc: func(ut ut.Translator) error {
				return ut.Add("max", "{0}", true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {

				minValue := fe.Param()
				var msg string

				k := fe.Type().Kind()

				if fe.Type().Kind() == reflect.Pointer {
					k = fe.Type().Elem().Kind()
				}

				switch k {
				case reflect.Array, reflect.Slice:
					msg = "cannot contain more than %s items"
				case
					reflect.Float32,
					reflect.Float64,
					reflect.Int,
					reflect.Int8,
					reflect.Int16,
					reflect.Int32,
					reflect.Int64,
					reflect.Uint,
					reflect.Uint8,
					reflect.Uint16,
					reflect.Uint32,
					reflect.Uint64:
					msg = "must be smaller than %s"
				case reflect.String:
					msg = "cannot be more than %s characters long"
				}

				t, _ := ut.T("min", fmt.Sprintf(msg, minValue))

				return t
			},
		},
	}
}

func (ve *ValidationExtension) Rules() []validation.CustomRule {
	return []validation.CustomRule{
		{
			Rule: "slug",
			// TODO: add some tests for this, I think it's right
			// also see about moving it somewhere so we can keep the compiled regex in memory
			Handler: func(fl validator.FieldLevel) bool {
				r, _ := regexp.Compile("^[a-z0-9]{1}[a-z0-9\\-]*$")

				return r.MatchString(fl.Field().String())
			},
		},
		{
			Rule: "postcode",
			// Pretty basic but covers the standard format of a postcode
			Handler: func(fl validator.FieldLevel) bool {
				r, _ := regexp.Compile("(?i)^[a-z]{1,2}\\d[a-z\\d]?\\s*\\d[a-z]{2}$")

				return r.MatchString(fl.Field().String())
			},
		},
	}
}

func NewValidationExtension() *ValidationExtension {
	return &ValidationExtension{}
}
