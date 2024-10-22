package validation

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

var ErrFieldNotFound = errors.New("field not found")

type StructMapMeta map[string]string

func (meta StructMapMeta) ByFieldPath(search string) (string, bool) {
	for k, v := range meta {
		if v == search {
			return k, true
		}
	}

	return search, false
}

type nullLogger struct{}

func (il *nullLogger) Debug(msg string, args ...any) {}
func (il *nullLogger) Warn(msg string, args ...any)  {}
func (il *nullLogger) Error(msg string, args ...any) {}

type Logger interface {
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type TagFinderFunc func(f reflect.StructField) string

type ValidationMapperOpt func(v *ValidationMapper)

func WithTagFinder(f TagFinderFunc) ValidationMapperOpt {
	return func(v *ValidationMapper) {
		v.tagFinder = f
	}
}

func WithLogger(l Logger) ValidationMapperOpt {
	return func(v *ValidationMapper) {
		v.logger = l
	}
}

func WithCacher(c Cacher) ValidationMapperOpt {
	return func(v *ValidationMapper) {
		v.cacher = c
	}
}

type Cacher interface {
	Add(key string, meta StructMapMeta)
	Get(key string) (StructMapMeta, bool)
}

type internalCache struct {
	store map[string]StructMapMeta
}

func (ic *internalCache) Add(key string, meta StructMapMeta) {
	ic.store[key] = meta
}

func (ic *internalCache) Get(key string) (StructMapMeta, bool) {
	v, found := ic.store[key]

	if !found {
		return StructMapMeta{}, false
	}

	return v, true
}

type ValidationMapper struct {
	tagFinder    TagFinderFunc
	logger       Logger
	cacheEnabled bool

	cacher Cacher
}

func getJsonTagNameForField(f reflect.StructField) string {
	jsonTag := f.Tag.Get("json")

	if jsonTag == "" {
		// Sensible default
		return f.Name
	}

	return strings.Split(jsonTag, ",")[0]
}

func (vm *ValidationMapper) getStructMap(t reflect.Type, useCache bool) StructMapMeta {
	cacheName := fmt.Sprintf("%s/%s", t.PkgPath(), t.Name())

	if useCache {
		v, found := vm.cacher.Get(cacheName)

		if found {
			return v
		}
	}

	props := StructMapMeta{}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		validationMap := f.Tag.Get("validationmap")
		tagName := vm.tagFinder(f)

		if f.Type.Kind() == reflect.Struct {
			sub := vm.getStructMap(f.Type, false)
			for k, v := range sub {
				props[fmt.Sprintf("%s.%s", tagName, k)] = v
			}
		}

		if validationMap != "" {
			props[tagName] = validationMap
		}

	}

	if useCache {
		vm.cacher.Add(cacheName, props)
	}

	return props
}

func (vm *ValidationMapper) Map(err ValidationError, in any) ValidationError {

	fldErrors := []FieldError{}

	t := reflect.TypeOf(in)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// TODO cache these results
	meta := vm.getStructMap(t, true)

	for _, err := range err.Errs {
		k, found := meta.ByFieldPath(err.Key)

		if !found {
			vm.logger.Warn("did not find validationmap field error", "field", err.Key, "type", t.Name())
		}

		fldErrors = append(fldErrors, FieldError{
			Key:    k,
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

func NewValidationMapper(opts ...ValidationMapperOpt) *ValidationMapper {
	v := &ValidationMapper{
		tagFinder: getJsonTagNameForField,
		logger:    &nullLogger{},
		cacher: &internalCache{
			store: map[string]StructMapMeta{},
		},
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}
