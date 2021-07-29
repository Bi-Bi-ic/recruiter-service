package company

import "errors"

var (
	// ErrCompanyNotFound ...
	ErrCompanyNotFound = errors.New("No Companies Found")
	// ErrNoMember occurs when this Company has not had any Members
	ErrNoMember = errors.New("No members of this Company")
)
