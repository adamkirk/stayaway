package validation

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var errorTranslations = []struct{
	rule string
	registerFunc func(ut ut.Translator) error
	translateFunc func(ut ut.Translator, fe validator.FieldError) string
}{
	{
		rule: "required",
		registerFunc: func(ut ut.Translator) error {
			return ut.Add("required", "is required", true)
		},
		translateFunc: func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required")
			
			return t
		},
	},
	{
		rule: "slug",
		registerFunc: func(ut ut.Translator) error {
			return ut.Add("slug", "must contain only alphanumeric and hyphen characters; cannot start with a hyphen", true)
		},
		translateFunc: func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("slug")
			
			return t
		},
	},
	{
		rule: "min",
		registerFunc: func(ut ut.Translator) error {
			return ut.Add("min", "{0}", true)
		},
		translateFunc: func(ut ut.Translator, fe validator.FieldError) string {

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

type Validator struct {}

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

		slog.Info("value", "field", field, "value", validationErr.Value())
		fieldErrors = append(fieldErrors, FieldError{
			Key: field,
			Errors: []string{validationErr.Translate(trans)},
		})
	}

	return ValidationError{
		Errs: fieldErrors,
	}
}

func NewValidator() *Validator {
	return &Validator{}
}

// TODO: add some tests for this, I think it's right
// also see about moving it somewhere so we can keep the compiled regex in memory
func slugValidator(fl validator.FieldLevel) bool {
	r, _ := regexp.Compile("^[a-z0-9]{1}[a-z0-9\\-]*$")

	return r.MatchString(fl.Field().String())
}

func init() {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ = uni.GetTranslator("en")
	
	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	
	validate.RegisterValidation("slug", slugValidator)
	

	for _, errorTranslation := range errorTranslations {
		validate.RegisterTranslation(
			errorTranslation.rule,
			trans,
			errorTranslation.registerFunc,
			errorTranslation.translateFunc,
		)
	}
}