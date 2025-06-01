package pagination

type Pagination struct {
	Page     int   `json:"page" query:"page"`
	PageSize int   `json:"page_size" query:"page_size"`
	Total    int64 `json:"total"`
}

func NewPagination(page, pageSize int) *Pagination {
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p *Pagination) WithTotal(total int64) *Pagination {
	p.Total = total
	return p
}

func (p *Pagination) Offset() int {
	if p.Page < 1 {
		p.Page = 1 // Default to page 1 if invalid
	}
	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 10 // Default page size if invalid
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) Limit() int {
	if p.PageSize < 1 || p.PageSize > 100 {
		return 10 // Default limit if PageSize is invalid
	}
	return p.PageSize
}
