package api

type V1ValidationErrorResponse struct {
	// A KV map of errors during validation.
	// The keys in the map correspond to the fields in the request.
	// The keys are using dot notation to flatten any nested fields.
	Errors map[string][]string `json:"errors"`
} // @name	V1.Response.Invalid

type V1GenericErrorResponse struct {
	// An error message that may give you indication of the problem.
	Message string `json:"message"`
} // @name V1.Response.Error