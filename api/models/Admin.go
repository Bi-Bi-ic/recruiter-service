package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	// "github.com/lib/pq"
)

// Admin ...
type Admin struct {
	Base
	GuestMode bool   `json:"guestmode"`
	Avatar    string `json:"avatar"`
	UserName  string `json:"username"`
	Fullname  string `json:"fullname"`

	BirthDay int64 `json:"birth_day"`

	NumberPhone string `json:"tell"`

	Email        string `json:"email" validate:"required,email"`
	MailContact  string `json:"mail_contact"`
	Password     string `json:"password,omitempty" validate:"required,min=6"`
	Token        string `json:"token,omitempty" sql:"-"`
	RefreshToken string `json:"refresh_token,omitempty" sql:"-"`
}

// ValidEmpty ...
func (user Admin) IsEmpty() bool {
	if user.Email != "" && user.Password != "" {
		return false
	}
	return true
}

// GenerateToken ...
func (user *Admin) GenerateToken() {
	//Create new JWT token for the newly registered account
	tk := &Token{
		UserId:    user.Base.ID,
		Email:     user.Email,
		Type:      Now,
		UserType:  AdminNormal,
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
		UserType:  AdminNormal,
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
