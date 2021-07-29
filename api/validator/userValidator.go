package validator

import (
	"reflect"
	"strings"

	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	govalidator "gopkg.in/go-playground/validator.v9"
)

// userVerify ...
type userVerify struct {
	validate *govalidator.Validate
}

// NewUserValidator ...
func NewUserValidator() UserValidator {
	valid := govalidator.New()
	return &userVerify{validate: valid}
}

// Valid ...
func (userCheck *userVerify) Valid(input models.User) error {

	// Use JSON names rather than Go struct Field names
	userCheck.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := userCheck.validate.Struct(input)
	if err != nil {
		return err
	}

	return nil
}

// Handle ...
func (userCheck *userVerify) Handle(fieldErr []govalidator.FieldError) int {
	for _, err := range fieldErr {
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
