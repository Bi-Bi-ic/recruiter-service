package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// TimeLine ...
type TimeLine struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Title    string    `json:"title"`
	SubTitle string    `json:"sub_title"`

	FromTime int64          `json:"from_time"`
	ToTime   int64          `json:"to_time"`
	DeleteAt gorm.DeletedAt `gorm:"index"`

	Description string    `json:"description"`
	UserID      uuid.UUID `json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (timeline *TimeLine) BeforeCreate(scope *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	scope.Statement.SetColumn("ID", id)
	return
}
