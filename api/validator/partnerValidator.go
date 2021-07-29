package validator

import (
	"reflect"
	"strings"

	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	govalidator "gopkg.in/go-playground/validator.v9"
)

const (
	// Login mode
	Login = "signin"
)

// PartnerVerify ...
type partnerVerify struct {
	validate *govalidator.Validate
	Mode     string
}

// NewPartnerValidator ...
func NewPartnerValidator() PartnerValidator {
	valid := govalidator.New()
	return &partnerVerify{validate: valid}
}

// SetMode ...
func (partnerCheck *partnerVerify) SetMode(mode string) {
	partnerCheck.Mode = mode
}

// Valid ...
func (partnerCheck *partnerVerify) Valid(input models.Partner) error {

	// Use JSON names rather than Go struct Field names
	partnerCheck.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := partnerCheck.validate.Struct(input)
	if err != nil {
		return err
	}

	return nil
}

func (partnerCheck *partnerVerify) Handle(fieldErr []govalidator.FieldError) int {
	for _, err := range fieldErr {

		// TODO: - Hot-Fix: Ignore company beforr create partner account
		// if err.ActualTag() == "required" {
		// 	if err.Field() == "name" {
		// 		if partnerCheck.Mode == Login {
		// 			continue
		// 		}
		// 	}
		// 	return models.ErrMissingField
		// }

		if err.ActualTag() == "required" {
			switch err.Field() {
			case "email":
				return code.EmailIsEmpty
			case "password":
				return code.PasswordIsEmpty
			}

		}
		switch err.Field() {
		case "email":
			return code.EmailNotFound
		case "password":
			return code.PasswordError
		}
	}

	return code.Ok
}
