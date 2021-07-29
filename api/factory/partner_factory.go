package factory

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/models"
	uuid "github.com/satori/go.uuid"
)

// Partnerable ...
type Partnerable struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	UserName         string    `json:"username"`
	PartnerName      string    `json:"name"`
	Token            string    `json:"token" sql:"-"`
	RefreshToken     string    `json:"refresh_token" sql:"-"`
	Avatar           string    `json:"avatar"`
	models.WorkSpace `json:"company"`
}

// PartnerInfoFactory this object create for anything what if you want about partner
type PartnerInfoFactory struct{}

// Create is a list of 'Partnerable' fixed from Partner entity
func (factory PartnerInfoFactory) Create(partner interface{}) Partnerable {
	customer := models.Partner{}

	err := mapstructure.Decode(partner, &customer)
	if err != nil {
		panic(err)
	}

	return Partnerable{
		ID:           customer.ID,
		Email:        customer.Email,
		UserName:     customer.UserName,
		PartnerName:  customer.PartnerName,
		Token:        customer.Token,
		RefreshToken: customer.RefreshToken,
		Avatar:       URL_SERVER + customer.Avatar,
		WorkSpace:    customer.WorkSpace,
	}
}

// CreateDetail ...
func (factory PartnerInfoFactory) CreateDetail(partner interface{}) models.Partner {
	partnerFactory := models.Partner{}

	err := mapstructure.Decode(partner, &partnerFactory)
	if err != nil {
		panic(err)
	}

	if partnerFactory.Avatar != "" {
		partnerFactory.Avatar  =  URL_SERVER + partnerFactory.Avatar
	}

	if partnerFactory.Cover != "" {
		partnerFactory.Cover  =  URL_SERVER + partnerFactory.Cover
	}

	return partnerFactory
}
