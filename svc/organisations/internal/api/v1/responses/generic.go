package responses

type ValidationErrorResponse struct {
	// A KV map of errors during validation.
	// The keys in the map correspond to the fields in the request.
	// The keys are using dot notation to flatten any nested fields.
	Errors map[string][]string `json:"errors"`
} // @name	V1.Response.Invalid

type GenericErrorResponse struct {
	// An error message that may give you indication of the problem.
	Message string `json:"message"`
} // @name V1.Response.Error

type SortOptionsResponseMeta struct {
	// Direction that was used when sorting results.
	OrderDirection string `json:"order_dir"`

	// Field by which the results are sorted.
	OrderBy string `json:"order_by"`
}// @name	V1.Response[Meta].SortOptions

type PaginationResponseMeta struct {
	// The page number for the results returned.
	Page int `json:"page"`

	// The amount of results oer page.
	PerPage int `json:"per_page"`

	// The total amount of pages available with the value for per_page.
	TotalPages int `json:"total_pages"`

	// The total amount of results in the system.
	TotalResults int `json:"total_results"`
}// @name	V1.Response[Meta].Pagination

type ListResponseMeta struct {
	PaginationResponseMeta `json:"pagination"`
	SortOptionsResponseMeta `json:"sorting"`
}// @name	V1.Response[Meta].List