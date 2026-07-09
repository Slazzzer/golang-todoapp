package core_pagination

// Page — страница списка с метаданными пагинации.
type Page[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func NewPage[T any](items []T, total, limit, offset int) Page[T] {
	if items == nil {
		items = []T{}
	}

	return Page[T]{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}
