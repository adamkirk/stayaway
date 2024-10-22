package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/danielgtaylor/huma/v2"
)

func ErrorHandler[Req any, Resp any](debugErrors bool, vm *validation.ValidationMapper, handler func(context.Context, *Req) (*Resp, error)) (func (ctx context.Context, req *Req) (*Resp, error)) {
	return func (ctx context.Context, req *Req) (*Resp, error) {
		resp, err :=  handler(ctx, req)

		if err == nil {
			return resp, nil
		}

		if err, ok := err.(common.ErrNotFound); ok {
			return nil, huma.Error404NotFound(err.Error())
		}

		validationError, ok := err.(validation.ValidationError)

		if ! ok {
			slog.Error("unhandled error", "error", err)

			if debugErrors {
				// This will automatically return the error message in the detail
				return resp, err
			}
			return resp, errors.New("This is an error in our system, please contact us!")
		}

		validationError = vm.Map(validationError, req)

		errors := []*huma.ErrorDetail{}

		for _, err := range validationError.Errs {
			errors = append(errors, &huma.ErrorDetail{
				Message: strings.Join(err.Errors, "|"),
				Location: err.Key,
			})
		}

		return resp, &huma.ErrorModel{
			Detail: "Validation failed",
			Status: http.StatusBadRequest,
			Errors: errors,
		}
	}
}
