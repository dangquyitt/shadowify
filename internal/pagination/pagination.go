package pagination

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func NewPagination(page, pageSize, total int) *Pagination {
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) Process() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 10
	}
}
