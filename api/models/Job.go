package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// OldJob ...
type OldJob struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	PostitionJob string    `json:"position_job"`
	Company      string    `json:"company"`
	FromTime     int64     `json:"from_time"` // time stamp
	ToTime       int64     `json:"to_time"`
	JobExplain   string    `json:"job_explain"`
	UserID       uuid.UUID `json:"-"`
}

// Salary ...
type Salary struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

//Skill ...
type Skill struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name   string    `json:"name"`
	Level  int       `json:"level"` //0 - 100
	UserID uuid.UUID `json:"-"`
}

// Language ...
type Language struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name   string    `json:"name"`
	Level  int       `json:"level"`
	UserID uuid.UUID `json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (oldJob *OldJob) BeforeCreate(scope *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	scope.Statement.SetColumn("ID", id)
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (skill *Skill) BeforeCreate(scope *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	scope.Statement.SetColumn("ID", id)
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (language *Language) BeforeCreate(scope *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	scope.Statement.SetColumn("ID", id)
	return
}
