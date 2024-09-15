package api

type V1SortOptionsResponseMeta struct {
	OrderDirection string `json:"order_dir"`
	OrderBy string `json:"order_by"`
}

type V1PaginationResponseMeta struct {
	Page int `json:"page"`
	PerPage int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type V1ListResponseMeta struct {
	V1PaginationResponseMeta `json:"pagination"`
	V1SortOptionsResponseMeta `json:"sorting"`
}