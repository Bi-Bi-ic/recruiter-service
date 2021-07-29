package models

// Address ...
type Address struct {
	Street   string `json:"street,omitempty"`
	District string `json:"district,omitempty"`
	Ward     string `json:"ward,omitempty"`
	City     string `json:"city,omitempty"`
}
