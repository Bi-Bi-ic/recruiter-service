package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Education ...
type Education struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	University    string    `json:"university"`
	Specialized   string    `json:"specialized"`
	Classname     string    `json:"classname"`
	Graduated     bool      `json:"graduated"`
	GraduatedDate int64     `json:"graduated_date"`

	UserID uuid.UUID `json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (education *Education) BeforeCreate(scope *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	scope.Statement.SetColumn("ID", id)
	return
}
