package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
)

// User information about a User
type User struct {
	Base
	GuestMode bool   `json:"guestmode"`
	Avatar    string `json:"avatar"`
	UserName  string `json:"username"`
	Fullname  string `json:"fullname"`
	Sex       string `json:"sex"` // male || female || others

	BirthDay    int64  `json:"birth_day"`
	Country     string `json:"country"`
	NumberPhone string `json:"tell"`

	Email        string `json:"email" validate:"required,email"`
	MailContact  string `json:"mail_contact"`
	Password     string `json:"password,omitempty" validate:"required,min=7"`
	Token        string `json:"token,omitempty" sql:"-"`
	RefreshToken string `json:"refresh_token,omitempty" sql:"-"`

	Address    `json:"address"`
	JobPurpose string `json:"job_purpose"`
	JobPlace   string `json:"job_place"`

	LevelPurpose  string         `json:"level_purpose"`
	Career        string         `json:"career"`
	SalaryPurpose string         `json:"salary_purpose"`
	JobKind       pq.StringArray `json:"job_kind" gorm:"type:text[]"` // main_personnel || freelace || part_time || intership
	LinkSocical   pq.StringArray `json:"link_social" gorm:"type:text[]"`

	OldJobs    []OldJob    `json:"old_jobs" gorm:"foreignkey:UserID"`
	Educations []Education `json:"educations" gorm:"foreignkey:UserID"`

	SpecializedExplain string     `json:"specialized_explain"`
	Skills             []Skill    `json:"skills" gorm:"foreignkey:UserID"`
	Languages          []Language `json:"languages" gorm:"foreignkey:UserID"`
	TimeLines          []TimeLine `json:"time_line" gorm:"foreignkey:UserID"`
	Cover              string     `json:"cover"`
}

// ValidEmpty ...
func (user User) ValidEmpty() bool {
	if user.Email != "" || user.Password != "" {
		return false
	}
	return true
}

// GenerateToken ...
func (user *User) GenerateToken() {
	//Create new JWT token for the newly registered account
	tk := &Token{
		UserId:    user.Base.ID,
		Email:     user.Email,
		Type:      Now,
		UserType:  UserNormal,
		CreateAt:  time.Now().String(),
		LimitedAt: time.Now().Add(24 * time.Hour).String(),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	//Create new JWT refresh token for the newly registered account
	re_tk := Token{
		UserId:    user.Base.ID,
		Email:     user.Email,
		Type:      Refresh,
		UserType:  UserNormal,
		CreateAt:  time.Now().String(),
		LimitedAt: time.Now().AddDate(0, 0, 30).String(),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().AddDate(0, 0, 30).Unix(),
		},
	}

	re_token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), re_tk)
	re_tokenString, _ := re_token.SignedString([]byte(os.Getenv("token_password")))
	user.RefreshToken = re_tokenString
}
