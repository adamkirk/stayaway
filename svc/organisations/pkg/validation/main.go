package validation

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var errorTranslations = []Translation{
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
}

var customRules = []CustomRule{
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
		// Pretty absic but covers the standard format of a postcode
		Handler: func(fl validator.FieldLevel) bool {
			r, _ := regexp.Compile("(?i)^[a-z]{1,2}\\d[a-z\\d]?\\s*\\d[a-z]{2}$")

			return r.MatchString(fl.Field().String())
		},
	},
}

var validate *validator.Validate
var trans ut.Translator

type FieldError struct {
	Key string
	Errors []string
}

type ValidationError struct {
	Errs []FieldError
}

func (err ValidationError) Error() string {
	return "invalid data"
}

type Validator struct {
	validate *validator.Validate
	trans ut.Translator
}

type Translation struct  {
	Rule string
	RegisterFunc func(ut ut.Translator) error
	TranslateFunc func(ut ut.Translator, fe validator.FieldError) string
}

type CustomRule struct {
	Rule string
	Handler func (fl validator.FieldLevel) bool
}

type Extension interface {
	Translations() []Translation
	Rules() []CustomRule
}

func (v *Validator) Validate(in any) error {
	
	err := validate.Struct(in)

	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)

	if ! ok {
		return errors.New("unable to handle validation")
	}

	// return err
	fieldErrors := []FieldError{}

	for _, validationErr := range errs {
		// The StructNamespace includes the type of the struct we're validating
		// we don't actually care about the top-level so we remove it
		// e.g. MyStructType.MyField.SubField becomes MyField.SubField
		field := strings.Join(
			strings.Split(validationErr.StructNamespace(), ".")[1:],
			".",
		)

		fieldErrors = append(fieldErrors, FieldError{
			Key: field,
			Errors: []string{validationErr.Translate(trans)},
		})
	}

	return ValidationError{
		Errs: fieldErrors,
	}
}

func NewValidator(extensions []Extension) *Validator {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ = uni.GetTranslator("en")
	
	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	
	for _, rule := range customRules {
		validate.RegisterValidation(rule.Rule, rule.Handler)
	}

	for _, errorTranslation := range errorTranslations {
		validate.RegisterTranslation(
			errorTranslation.Rule,
			trans,
			errorTranslation.RegisterFunc,
			errorTranslation.TranslateFunc,
		)
	}

	for _, ext := range extensions {
		for _, rule := range ext.Rules() {
			validate.RegisterValidation(rule.Rule, rule.Handler)
		}

		for _, translation := range ext.Translations() {
			validate.RegisterTranslation(
				translation.Rule,
				trans,
				translation.RegisterFunc,
				translation.TranslateFunc,
			)
		}
	}

	return &Validator{
		validate: validate,
		trans: trans,
	}
}


