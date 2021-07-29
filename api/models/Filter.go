package models

// Filter ...
type Filter struct {
	Position []string `json:"positions"`
	JobKind  []string `json:"job_kinds"`
	District []string `json:"districts"`

	/* gte, lte, gt, lt*/

	TotalContents int `json:"total_contents"`

	Rows interface{} `json:"rows"`
}
