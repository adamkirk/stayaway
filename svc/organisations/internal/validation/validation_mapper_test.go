package validation

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CreatePersonRequestAddress struct {
	Line1    *string `json:"line_1" validationmap:"Address.Street"`
	Postcode *string `json:"post_code" validationmap:"Address.Postcode"`
}

type CreatePersonRequest struct {
	Name    *string                    `json:"name" validationmap:"FullName"`
	Email   *string                    `json:"email" validationmap:"EmailAddress"`
	Address CreatePersonRequestAddress `json:"address" validationmap:"Address"`
}

type FlattenedPersonRequestAddress struct {
	Line1    *string `json:"line_1" validationmap:"Street"`
	Postcode *string `json:"post_code" validationmap:"Postcode"`
}

type FlattenedPersonRequest struct {
	Name    *string                       `json:"name" validationmap:"FullName"`
	Email   *string                       `json:"email" validationmap:"EmailAddress"`
	Address FlattenedPersonRequestAddress `json:"address"`
}

func ptr[T any](value T) *T {
	return &value
}

func TestBuildValidationMap(t *testing.T) {
	tests := []struct {
		name      string
		in        any
		expectErr error
		expect    StructMapMeta
	}{
		{
			name:      "nested props",
			in:        CreatePersonRequest{},
			expectErr: nil,
			expect: StructMapMeta{
				"address":           "Address",
				"name":              "FullName",
				"email":             "EmailAddress",
				"address.line_1":    "Address.Street",
				"address.post_code": "Address.Postcode",
			},
		},

		{
			name:      "flattened target",
			in:        FlattenedPersonRequest{},
			expectErr: nil,
			expect: StructMapMeta{
				"address.line_1":    "Street",
				"address.post_code": "Postcode",
				"email":             "EmailAddress",
				"name":              "FullName",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {

			res := mapJsonFieldsToMapTags(reflect.TypeOf(test.in))

			// assert.Equal(tt, test.expectErr, err)

			assert.Equal(tt, test.expect, res)
		})
	}
}
