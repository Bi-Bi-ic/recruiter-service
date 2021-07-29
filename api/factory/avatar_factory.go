package factory

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/models"
	uuid "github.com/satori/go.uuid"
)

// UserAvatar ...
type UserAvatar struct {
	ID       uuid.UUID `json:"id"`
	CreateAt int64     `json:"create_at"`
	UpdateAt int64     `json:"update_at"`
	Fullname string    `json:"fullname"`
	Avatar   string    `json:"avatar_image"`
}

// PartnerAvatar ...
type PartnerAvatar struct {
	ID          uuid.UUID `json:"id"`
	CreateAt    int64     `json:"create_at"`
	UpdateAt    int64     `json:"update_at"`
	UserName    string    `json:"username"`
	PartnerName string    `json:"name"`
	Avatar      string    `json:"avatar"`
}

// AvatarInfoFactory this object create for anything what if you want about avatar
type AvatarInfoFactory struct{}

// UserAvatar return UserAvatar-Able
func (factory AvatarInfoFactory) UserAvatar(avatar interface{}) UserAvatar {
	customer := models.User{}

	err := mapstructure.Decode(avatar, &customer)
	if err != nil {
		panic(err)
	}

	userFactory := UserAvatar{
		ID:       customer.ID,
		CreateAt: customer.CreateAt,
		UpdateAt: customer.UpdateAt,
		Fullname: customer.Fullname,
		Avatar:   customer.Avatar,
	}

	return userFactory
}

// PartnerAvatar return PartnerAvatar-Able
func (factory AvatarInfoFactory) PartnerAvatar(avatar interface{}) PartnerAvatar {
	customer := models.Partner{}

	err := mapstructure.Decode(avatar, &customer)
	if err != nil {
		panic(err)
	}

	partnerFactory := PartnerAvatar{
		ID:          customer.ID,
		CreateAt:    customer.CreateAt,
		UpdateAt:    customer.UpdateAt,
		UserName:    customer.UserName,
		PartnerName: customer.PartnerName,
		Avatar:      customer.Avatar,
	}
	return partnerFactory
}
