package validation

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"sort"
	"strings"
)

type ValidationMapper struct {}

func findStructFieldForValidationField(t reflect.Type, fieldName string, carry string, fieldNamePrefix string) (string, bool) {
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

func getJsonNameForField(f reflect.StructField) string {
	jsonTag := f.Tag.Get("json")
	return strings.Split(jsonTag, ",")[0]
}

var errFieldNotFound = errors.New("field not found")

type StructMapMeta map[string]string

func (meta StructMapMeta) ByFieldPath(search string) (string, bool) {
	for k, v := range meta {
		if v == search {
			return k, true
		}
	}

	return search, false
}

func mapJsonFieldsToMapTags(t reflect.Type) (StructMapMeta) {
	props := StructMapMeta{}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		slog.Debug("parsing validationmap", "struct", t.Name(), "field", f.Name)
		validationMap := f.Tag.Get("validationmap")
		jsonName := getJsonNameForField(f)

		if f.Type.Kind() == reflect.Struct {
			sub := mapJsonFieldsToMapTags(f.Type)
			for k, v := range sub {
				props[fmt.Sprintf("%s.%s", jsonName, k)] = v
			}
		} 

		if validationMap != "" {
			props[jsonName] = validationMap
		}

	}

	return props
}

// TODO: add an option to change from using json tag to another tag
// This will support query param validation
func (vm *ValidationMapper) Map(err ValidationError, req any) ValidationError {

	fldErrors := []FieldError{}

	t := reflect.TypeOf(req)

	// TODO cache these results
	meta := mapJsonFieldsToMapTags(t)

	fmt.Printf("%#v\n\n", meta)

	for _, err := range err.Errs {
		k, found := meta.ByFieldPath(err.Key)

		if ! found {
			slog.Warn("did not find validationmap field error", "field", err.Key, "type", t.Name())
		}

		fldErrors = append(fldErrors, FieldError{
			Key: k,
			Errors: err.Errors,
		})
	}

	// Ensure they're in a  prefictable order
	sort.Slice(fldErrors, func(i, j int) bool {
		return fldErrors[i].Key < fldErrors[j].Key
	})

	return ValidationError{
		Errs: fldErrors,
	}
}

func NewValidationMapper() *ValidationMapper {
	return &ValidationMapper{}
}