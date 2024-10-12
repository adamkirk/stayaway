package v1

import "fmt"

type HttpError interface {
	HttpStatusCode() int
}

type HttpDebuggableError interface {
	DebugError() string
}

type HttpResponseBuilder interface {
	BuildResponse() map[string]any
}

type ErrNotFound struct {
	ResourceName string
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", err.ResourceName)
}

func (err ErrNotFound) HttpStatusCode() int {
	return 404;
}

type ErrConflict struct {
	Message string
}

func (err ErrConflict) Error() string {
	return fmt.Sprintf("resource conflict: %s", err.Message)
}

func (err ErrConflict) HttpStatusCode() int {
	return 409;
}

type ErrBadRequest struct {
	Message string
	DebugMessage string
}

func (err ErrBadRequest) Error() string {
	return err.Message
}

func (err ErrBadRequest) HttpStatusCode() int {
	return 400;
}

func (err ErrBadRequest) DebugError() string {
	return err.DebugMessage
}
