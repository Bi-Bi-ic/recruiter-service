package partner

import "errors"

var (
	// ErrAdminNotFound ...
	ErrAdminNotFound = errors.New("Admin does not Exit")
	// ErrPartnerNotFound ...
	ErrPartnerNotFound = errors.New("Partner does not Exit")
	// ErrorCompanyOwned occurs when Company is Belong to another Partner
	ErrorCompanyOwned = errors.New("This Company is owned")
	// ErrNotBelongCompany occurs when Partner has not became Member any Companies yet
	ErrNotBelongCompany = errors.New("From now You are have not Joined any Companies yet")
	// ErrRequestWaiting occurs when Partner's Request is not accepted by Admin
	ErrRequestWaiting = errors.New("Your Request Joining Company is not permitted by Admin")
	// ErrAdminFlag occurs when Partner is not the Admin: Owner of Company
	ErrAdminFlag = errors.New("You are not the Admin of Company")
	// ErrNotRequest occurs when Partner not request joining this Company
	ErrNotRequest = errors.New("This Partner does not request Joining this Company")
	// ErrMemberRequest occurs when Partner is a member but wanna send another request Joining
	ErrMemberRequest = errors.New("This Partner became Member of a Company ! Request did not Succeed")
	// ErrMemberFlag occurs when Partner is a Member of Company
	ErrMemberFlag = errors.New("This Partner was already Member of a Company")
	// ErrPartnerContents occurs when Getting Partner's Contents has Errors
	ErrPartnerContents = errors.New("Some errors occurs with Partner's Contents")
)
