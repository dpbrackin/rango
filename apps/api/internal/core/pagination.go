package core

type OrderDirection string

const (
	OrderAsc  OrderDirection = "ASC"
	OrderDesc OrderDirection = "DESC"
)

type OrderBy struct {
	Field     string
	Direction OrderDirection
}

type Pagination struct {
	PageSize int
	Page     int
	OrderBy  []OrderBy
}
