package factory

import (
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/models"
	uuid "github.com/satori/go.uuid"
)

// Postable is struct to return
type Postable struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`

	CreateAt    int64  `json:"created_at"`
	UpdateAt    int64  `json:"updated_at"`
	DeleteAt    *int64 `json:"deleted_at" sql:"index"`
	TimeExpired int64  `json:"time_expired"`

	Type string `json:"type"`

	models.Creator `json:"creator"`

	Language     string `json:"language"`
	Position     string `json:"position"`
	Description  string `json:"descriptions"`
	Requirements string `json:"requirements"`

	JobKind         pq.StringArray `json:"job_kind"`
	models.Salary   `json:"salary"`
	LinkSocialMedia string `json:"link_social_media"`
	Benefits        string `json:"benefits"`

	TotalLike int `json:"total_like"`
	TotalView int `json:"total_view,omitempty"`

	Cover string `json:"cover"`

	Tags pq.StringArray `json:"tags" gorm:"type:text[]"`
}

// PostAble show post's details after go through factory
type PostAble struct {
	ID       uuid.UUID `json:"id"`
	CreateAt int64     `json:"created_at"`
	UpdateAt int64     `json:"updated_at"`
	DeleteAt *int64    `json:"deleted_at,omitempty"`

	Title       string      `json:"title"`
	TimeExpired int64       `json:"time_expired"`
	Type        string      `json:"type"`
	Creator     creatorAble `json:"creator"`

	Language      string         `json:"language"`
	Position      string         `json:"position"`
	Description   string         `json:"descriptions"`
	Requirements  string         `json:"requirements"`
	JobKind       pq.StringArray `json:"job_kind"`
	models.Salary `json:"salary"`

	LinkSocialMedia string `json:"link_social_media,omitempty"`
	Benefits        string `json:"benefits,omitempty"`
	TotalLike       int    `json:"total_like"`
	TotalView       int    `json:"total_view"`

	Cover   string         `json:"cover,omitempty"`
	Tags    pq.StringArray `json:"tags,omitempty"`
	Feature int            `json:"feature"`
}

type creatorAble struct {
	ID          string `json:"id"`
	UserName    string `json:"username"`
	PartnerName string `json:"partnername,omitempty"`

	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name"`

	models.Address `json:"address"`
	MailContact    string `json:"mail_contact,omitempty"`
	Phone          string `json:"phone,omitempty"`

	Avatar string `json:"avatar,omitempty"`
	Link   string `json:"link,omitempty"`
}

// PostInfoFactoty this object create for anything what if you want about post
type PostInfoFactoty struct{}

// CreatedWithCreator ...
func (factory PostInfoFactoty) CreatedWithCreator(post models.Post, creator models.Partner) PostAble {

	var deletedAt *int64
	if !post.DeleteAt.Time.IsZero() {
		deletedTime := post.DeleteAt.Time.Unix()
		deletedAt = &deletedTime
	}

	return PostAble{
		ID:       post.ID,
		CreateAt: post.CreateAt,
		UpdateAt: post.UpdateAt,
		DeleteAt: deletedAt,

		Title: post.Title,
		Type:  post.Type,

		Creator: creatorAble{
			ID:          creator.ID.String(),
			UserName:    creator.UserName,
			PartnerName: creator.PartnerName,
			CompanyID:   creator.CompanyID,
			CompanyName: creator.Name,
			Address:     creator.Address,
			MailContact: creator.MailContact,
			Phone:       creator.Phone,
			Avatar:      creator.Avatar,
			Link:        creator.Link,
		},

		Language:     post.Language,
		Position:     post.Position,
		Description:  post.Description,
		Requirements: post.Requirements,
		JobKind:      post.JobKind,
		Salary:       post.Salary,

		LinkSocialMedia: post.LinkSocialMedia,
		Benefits:        post.Benefits,
		TotalLike:       post.TotalLike,
		TotalView:       post.TotalView,

		Cover:   post.Cover,
		Tags:    post.Tags,
		Feature: post.Feature,
	}
}

// NewListPost is a list of 'Postable' fixed from Post entity
func (factory PostInfoFactoty) NewListPost(posts interface{}) []Postable {
	listPosts := []Postable{}

	mapstructure.Decode(posts, &listPosts)

	return listPosts
}

//NewPost is a 'Postable' factory fixed from Post entity
func (factory PostInfoFactoty) NewPost(post interface{}) Postable {
	postDetails := Postable{}

	err := mapstructure.Decode(post, &postDetails)
	if err != nil {
		panic(err)
	}

	return postDetails
}
