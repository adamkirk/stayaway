package v1

import (
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/responses"
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
	"github.com/labstack/echo/v4"
)

func translateErrToHttpErr(err error) HttpError {
	switch t := err.(type) {
	default:
		return nil
	case common.ErrNotFound:
		return ErrNotFound{
			ResourceName: t.ResourceName,
		}
	case common.ErrConflict:
		return ErrConflict{
			Message: t.Message,
		}
	}
}

func handleValidationError(ctx echo.Context, errs validation.ValidationError) {
	respErrors := map[string][]string{}

	for _, err := range errs.Errs {
		respErrors[err.Key] = err.Errors
	}

	respBody := responses.ValidationErrorResponse{
		Errors: respErrors,
	}

	ctx.JSON(422, respBody)
}

func NewErrorHandler(debugErrorsEnabled bool) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			err := next(ctx)
			if err == nil {
				return nil
			}

			if err, ok := err.(validation.ValidationError); ok {
				handleValidationError(ctx, err)
				return nil
			}

			respBody := map[string]any{}

			var httpErr HttpError

			if translated := translateErrToHttpErr(err); translated != nil {
				httpErr = translated
			} else {
				translated, ok := err.(HttpError)

				if ok {
					httpErr = translated
				}
			}

			if httpErr != nil {
				respBody["message"] = err.Error()

				debuggableErr, ok := err.(HttpDebuggableError)

				if ok && debugErrorsEnabled {
					respBody["debug"] = map[string]any{
						"error": debuggableErr.DebugError(),
					}
				}

				ctx.JSON(httpErr.HttpStatusCode(), respBody)

				return nil
			}

			if debugErrorsEnabled {
				respBody["debug"] = map[string]any{
					"error": err.Error(),
				}
			}

			respBody["message"] = "internal server error"
			ctx.JSON(500, respBody)

			return nil
		}
	}
}
