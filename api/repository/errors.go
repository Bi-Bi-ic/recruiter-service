package repository

import "errors"

var (
	// ErrMalformedID occurs when ID is mismatch uuid version
	ErrMalformedID = errors.New("Malformed ID")
	// ErrEmailNotFound ...
	ErrEmailNotFound = errors.New("Email is unaddressable")
	// ErrEmailSyntax occurs when email's format is invalid
	ErrEmailSyntax = errors.New("Email address Format is invalid or typos")
	// ErrEmailUsed ...
	ErrEmailUsed = errors.New("Email is used by another User")
	//ErrPasswordTooShort occurs when Password  < 6 character
	ErrPasswordTooShort = errors.New("Password is invalid or Too Short")
	// ErrRequestTooLong occurs when Repository can not handle Result as excepted
	ErrRequestTooLong = errors.New("Request Time Out. Please retry ... again")
	// ErrLogin occurs when Login's Credentials : Email, Password are mismatched
	ErrLogin = errors.New("Login Credentials is mismatched")
)
