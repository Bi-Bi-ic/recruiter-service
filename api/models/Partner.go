package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Partner ... struct
type Partner struct {
	Base
	Address `json:"address"`

	Email       string `json:"email" validate:"required,email"`
	UserName    string `json:"username"`
	PartnerName string `json:"name"`
	MailContact string `json:"mail_contact"`
	Password    string `json:"password,omitempty" validate:"required,min=7"`

	Category    string `json:"category"`
	Description string `json:"description"`

	Token        string `json:"token,omitempty" sql:"-"`
	RefreshToken string `json:"refresh_token,omitempty" sql:"-"`

	PostAvailable string `json:"post_available,omitempty" sql:"-"`
	PostExpired   string `json:"post_expired,omitempty" sql:"-"`

	Link   string `json:"link"`
	Phone  string `json:"phone"`
	Avatar string `json:"avatar"`
	Cover  string `json:"cover"`

	TotalLike uint `json:"total_like,omitempty"`

	WorkSpace `json:"company"`
}

// WorkSpace ...
type WorkSpace struct {
	CompanyID  string `json:"id"`
	Name       string `json:"name" validate:"required"`
	Permission string `json:"permission"`
}

// ValidEmpty ...
func (partner Partner) ValidEmpty() bool {
	if partner.Email != "" || partner.Password != "" {
		return false
	}
	return true
}

// GenerateToken ...
func (partner *Partner) GenerateToken(mode UserMode) {
	//Create new JWT token for the newly registered account
	tk := &Token{
		UserId:    partner.Base.ID,
		Email:     partner.Email,
		Type:      Now,
		UserType:  mode,
		CreateAt:  time.Now().String(),
		LimitedAt: time.Now().Add(24 * time.Hour).String(),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	partner.Token = tokenString

	//Create new JWT refresh token for the newly registered account
	re_tk := &Token{
		UserId:    partner.Base.ID,
		Email:     partner.Email,
		Type:      Refresh,
		UserType:  mode,
		CreateAt:  time.Now().String(),
		LimitedAt: time.Now().AddDate(0, 0, 30).String(),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().AddDate(0, 0, 30).Unix(),
		},
	}

	re_token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), re_tk)
	re_tokenString, _ := re_token.SignedString([]byte(os.Getenv("token_password")))
	partner.RefreshToken = re_tokenString
}
