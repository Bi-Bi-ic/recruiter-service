package models

// List Of Unique APIs Code
const (
	// Partner
	PartnerRequest = 20002
	// PartnerAvatarUpload
	PartnerAvatarUpload = 30003
	// TrackingSuccess
	TrackingSuccess = 40004
	// PartnerInfo
	PartnerInfo = 50005

	//Contents
	Contents = 60006

	// ErrValidate ...
	ErrValidate = 99999
)

// Code ...
type Code int

// CodeStatus ...
type CodeStatus struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

const (
	// OK ...
	OK Code = 20000
	// BadRequest ...
	BadRequest Code = 10001
	// ErrMissingField ...
	ErrMissingField Code = 90000
	// ErrEmailType ...
	ErrEmailType Code = 90001
	// ErrPasswdNotStrong ...
	ErrPasswdNotStrong Code = 90002
)

// SendMessage handle message error
func (c Code) SendMessage() CodeStatus {
	return CodeStatus{c, c.asString()}
}

// Return Message for each status_code input
func (c Code) asString() string {
	switch c {
	case BadRequest, ErrMissingField, ErrEmailType, ErrPasswdNotStrong:
		return "Bad Request"
	default:
		return ""
	}
}
