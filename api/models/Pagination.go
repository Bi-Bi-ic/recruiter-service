package models

// Pagination ...
type Pagination struct {
	Limit         int  `json:"limit"`
	Offset        int  `json:"offset"`
	Sort          Sort `json:"sort"`
	TotalContents int  `json:"total_contents"`

	Rows interface{} `json:"rows"`
}

// Sort ...
type Sort string

// each SortType ...
const (
	Lastest Sort = "now"
	Oldest  Sort = "oldest"
)

// AsString ...
func (sort Sort) AsString() string {
	switch sort {
	case Lastest:
		return "create_at desc"
	case Oldest:
		return "create_at asc"
	default:
		return ""
	}
}
