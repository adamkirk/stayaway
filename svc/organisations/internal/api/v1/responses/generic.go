package responses

type SortOptionsResponseMeta struct {
	OrderDirection string `json:"order_dir" doc:"Direction that was used when sorting results."`
	OrderBy string `json:"order_by" doc:"Field by which the results are sorted."`
}

type PaginationResponseMeta struct {
	Page int `json:"page" doc:"The page number for the results returned."`
	PerPage int `json:"per_page" doc:"The amount of results per page."`
	TotalPages int `json:"total_pages" doc:"The total amount of pages available with the value for per_page."`
	TotalResults int `json:"total_results" doc:"The total amount of results in the system."`
}

type ListResponseMeta struct {
	PaginationResponseMeta  `json:"pagination"`
	SortOptionsResponseMeta `json:"sorting"`
}

type GenericResponseBody[M any, D any] struct {
	Data D `json:"data"`
	Meta M `json:"meta"`
}

type NoMeta struct {}

type GenericResponse[M any, D any] struct {
	Body GenericResponseBody[M, D]
}

type NoContent struct {
	Status int
}