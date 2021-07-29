package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Company ...
type Company struct {
	ID   uuid.UUID `json:"id" gorm:"primary_key"`
	Name string    `json:"name"`

	CreatorID string   `json:"creator_id,omitempty"`
	Members   []Member `json:"members,omitempty" gorm:"-"`
}

// Member ...
type Member struct {
	ID         string `json:"id"`
	Permission string `json:"permission"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (company *Company) BeforeCreate(scope *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	scope.Statement.SetColumn("ID", id)
	return
}
