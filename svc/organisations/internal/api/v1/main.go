package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"reflect"

	"github.com/labstack/echo/v4"
)

// RequestWithRawBody allows us to check whether fields were defined at all in
// the json body. Realy just a way to define the different between setting
// something to null, and something not being present
type RequestWithRawBody interface {
	IncludeRawBody(raw map[string]any)
	FieldWasPresent(fld string) bool
}

func bindRequest(req any, ctx echo.Context) error {
	if reflect.ValueOf(req).Kind() != reflect.Ptr {
		slog.Error("cannot bind to non pointer", "path", ctx.Path())

		return errors.New("cannot bind request to non pointer value")
	}

	if reqWithRaw, ok := req.(RequestWithRawBody); ok {
		b, err := io.ReadAll(ctx.Request().Body)

		if err != nil {
			return err
		}

		raw := map[string]any{}
		if err := json.Unmarshal(b, &raw); err != nil {
			return err
		}

		reqWithRaw.IncludeRawBody(raw)

		ctx.Request().Body = io.NopCloser(bytes.NewBuffer(b))
	}

	if err := ctx.Bind(req); err != nil {
		return err
	}

	return nil
}
