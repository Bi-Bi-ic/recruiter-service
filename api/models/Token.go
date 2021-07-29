package models

import (
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

/*
JWT claims struct
*/

// TokenMode ...
type TokenMode string

// UserMode ...
type UserMode string

// Token ...
type Token struct {
	UserId    uuid.UUID
	Email     string
	CreateAt  string
	LimitedAt string
	Type      TokenMode
	UserType  UserMode
	jwt.StandardClaims
}

const (
	// Now ...
	Now TokenMode = "token_now"
	// Refresh ...
	Refresh TokenMode = "token_refresh"
)

// Token mode type
const (
	// User
	UserNormal UserMode = "user_normal"

	// Partner
	PartnerNormal  UserMode = "partner_normal"
	PartnerStandby UserMode = "partner_standby"

	// Admin
	AdminNormal UserMode = "admin_normal"

	// Loction
	LocationNormal UserMode = "location_normal"

	// Post
	PostNormal         UserMode = "post_normal"
	IntroductionNormal UserMode = "introduction_normal"

	// File
	CoverNormal UserMode = "cover_normal"

	// TimeLine
	TimeLineNormal UserMode = "timeline_normal"
)
