package api

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type ValidationMapper struct {}

// Seems rather inefficient to do it this way, but can tidy up later...
// func findStructFieldForValidationField(t reflect.Type, fieldName string, carry string) (string, bool) {
// 	parts := strings.Split(fieldName, ".")

// 	search := parts[0]
// 	var jsonField string
// 	var f reflect.StructField

// 	for i := 0; i < t.NumField(); i++ {
// 		f = t.Field(i)

// 		if f.Tag.Get("validationmap") == search {
// 			if jsonTag := f.Tag.Get("json"); jsonTag != "" {
// 				jsonField = strings.Split(jsonTag, ",")[0]
// 				break
// 			}

// 		}
// 	}

// 	if len(parts) > 0 {
// 		findStructFieldForValidationField()
// 	}

// 	return "", false
// }

func findStructFieldForValidationField(t reflect.Type, fieldName string, carry string, fieldNamePrefix string) (string, bool) {
	fmt.Printf("searching for: %s\n", fieldName)
	parts := strings.Split(fieldName, ".")

	search := parts[0]

	if fieldNamePrefix != "" {
		search = fmt.Sprintf("%s.%s", fieldNamePrefix, search)
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Tag.Get("validationmap") != search {
			continue
		}

		jsonTag := f.Tag.Get("json")
		jsonFieldName := strings.Split(jsonTag, ",")[0]

		if carry == "" {
			carry = jsonFieldName
		} else {
			carry = fmt.Sprintf("%s.%s", carry, jsonFieldName)
		}

		if len(parts) > 1 {
			return findStructFieldForValidationField(f.Type, strings.Join(parts[1:], "."), carry, parts[0])
		}
	}

	return carry, true
}

func (vm *ValidationMapper) Map(err validation.ValidationError, req any) validation.ValidationError {

	fldErrors := []validation.FieldError{}

	t := reflect.TypeOf(req)

	for _, err := range err.Errs {
		fld, found := findStructFieldForValidationField(t, err.Key, "", "")

		if ! found {
			continue
		}

		fldErrors = append(fldErrors, validation.FieldError{
			Key: fld,
			Errors: err.Errors,
		})
	}

	// Ensure they're in a  prefictable order
	sort.Slice(fldErrors, func(i, j int) bool {
		return fldErrors[i].Key < fldErrors[j].Key
	})

	return validation.ValidationError{
		Errs: fldErrors,
	}
}

func NewValidationMapper() *ValidationMapper {
	return &ValidationMapper{}
}