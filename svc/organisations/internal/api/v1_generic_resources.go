package api

type V1SortOptionsResponseMeta struct {
	// Direction that was used when sorting results.
	OrderDirection string `json:"order_dir"`

	// Field by which the results are sorted.
	OrderBy string `json:"order_by"`
}// @name	V1.Response[Meta].SortOptions

type V1PaginationResponseMeta struct {
	// The page number for the results returned.
	Page int `json:"page"`

	// The amount of results oer page.
	PerPage int `json:"per_page"`

	// The total amount of pages available with the value for per_page.
	TotalPages int `json:"total_pages"`

	// The total amount of results in the system.
	TotalResults int `json:"total_results"`
}// @name	V1.Response[Meta].Pagination

type V1ListResponseMeta struct {
	V1PaginationResponseMeta `json:"pagination"`
	V1SortOptionsResponseMeta `json:"sorting"`
}// @name	V1.Response[Meta].List