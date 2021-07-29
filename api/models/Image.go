package models

// Image ...
type Image struct {
	ID     string `json:"name" gorm:"primary_key"`
	Source []byte `gorm:"size: 70000"`
}
