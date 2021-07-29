package message

// List Unique Error Messages ...
const (
	BadRequest          = "Bad Request"
	DataInvalid         = "Data is Invalid"
	EmailNotFound       = "Email Not Found"
	EmailIsUsed         = "Email's Address Is already being used"
	QueryError          = "Query has Error"
	PasswordNotMatch    = "Password is Mismatched"
	OK                  = "Ok"
	InternalServerError = "Server has Errors can not Hanle "
	MalFormedID         = "ID is malformed"
	LocationError       = "Location Not Found"
	ResourceNotFound    = "This Resource is Not Found"
	ImageError          = "Image Upload Error"
	NotFound            = "This Resource Not Found"
	IntroductionExited  = "Introduction Existed"
)

// List unique Message ...
const (
	// Login ..
	Login = "Logged In"

	// Created ...
	UserCreated         = "User's Account Created"
	PartnerCreated      = "Partner's Account Created"
	PostCreated         = "Post Created"
	IntroductionCreated = "Introduction Created"
	TimeLineCreated     = "TimeLine Created"
	TimeLineDeleted     = "TimeLine Deleted"

	// Success ..
	IntroductionSuccess = "Found Introduction"
	UserSuccess         = "Found User"
	PartnerSuccess      = "Found Partner"

	// Location Navigation
	LocationSuccess = "Found all locations"

	// LikeSuccess ...
	LikeSuccess = "Updated Like"
	// Update ...
	Uploaded = "Success Uploaded"
	Updated  = "Success Updated"
)
