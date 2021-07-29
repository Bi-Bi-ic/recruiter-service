package code

// List of Unique Error Codes ...
const (
	// Like ...
	LikeSuccess int = 50000

	// Post ...
	PostError int = 50005

	// UploadImageError
	UploadError int = 90003

	//IsEmpty ...
	EmailIsEmpty    int = 90004
	PasswordIsEmpty int = 90004
	DataIsEmpty     int = 90004

	// EmailNotFound ...
	EmailNotFound int = 90005

	// QueryError ...
	QueryError int = 90006

	// PasswordError ...
	PasswordError int = 90007

	// EmailIsUsed ...
	EmailIsUsed int = 90008

	// Resource ...
	ResourceError int = 90009

	// LocationError ...
	LocationError int = 91001

	// MalFormedID ...
	MalFormedID int = 99000

	// InternalServerError ...
	InternalServerError int = 99999
)

// List Of Unique Success Code ...
const (
	Ok       int = 90000
	Created  int = 90001
	Accepted int = 90002

	Uploaded int = 90009

	Deleted int = 99998

	LocationSuccess int = 91000
)
