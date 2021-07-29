package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"
)

//Post declared Post's base informations
type Post struct {
	Base
	Title       string `json:"title"`
	TimeExpired int64  `json:"time_expired"`
	Type        string `json:"type"`
	Creator

	Language     string         `json:"language"`
	Position     string         `json:"position"`
	Description  string         `json:"descriptions"`
	Requirements string         `json:"requirements"`
	JobKind      pq.StringArray `json:"job_kind" gorm:"type:text[]"`
	Salary       `json:"salary"`

	LinkSocialMedia string `json:"link_social_media"`
	Benefits        string `json:"benefits"`
	TotalLike       int
	TotalView       int

	Cover   string         `json:"cover"`
	Tags    pq.StringArray `json:"tags" gorm:"type:text[]"`
	Feature int
}

//Creator show creator's informations
type Creator struct {
	CreatorID   string `json:"id"`
	UserName    string `json:"username" gorm:"-"`
	PartnerName string `json:"name" gorm:"-"`

	CompanyID   string `json:"company_id" gorm:"-"`
	CompanyName string `json:"company_name" gorm:"-"`

	Address     `json:"address" gorm:"-"`
	MailContact string `json:"mail_contact" gorm:"-"`
	Phone       string `json:"phone" gorm:"-"`

	CreatorAvatar string `json:"avatar" gorm:"-"`
	Link          string `json:"link" gorm:"-"`
}

// BeforeCreate generate uuid type for ID-Post
func (post *Post) BeforeCreate(scope *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	scope.Statement.SetColumn("ID", id)
	return
}
