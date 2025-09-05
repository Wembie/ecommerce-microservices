package models

type PaginatedResponse[T any] struct {
	Items        []T  `json:"items"`
	Page         int  `json:"page"`
	Size         int  `json:"size"`
	Total        int  `json:"total"`
	Pages        int  `json:"pages"`
	NextPage     *int `json:"next_page,omitempty"`
	PreviousPage *int `json:"previous_page,omitempty"`
}

func NewPaginatedResponse[T any](items []T, page, size, total int) *PaginatedResponse[T] {
	pages := (total + size - 1) / size
	if pages == 0 {
		pages = 1
	}
	
	response := &PaginatedResponse[T]{
		Items: items,
		Page:  page,
		Size:  size,
		Total: total,
		Pages: pages,
	}
	
	if page < pages-1 {
		nextPage := page + 1
		response.NextPage = &nextPage
	}
	
	if page > 0 {
		prevPage := page - 1
		response.PreviousPage = &prevPage
	}
	
	return response
}