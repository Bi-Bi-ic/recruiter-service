package repository

import (
	"github.com/rgrs-x/service/api/models"
	u "github.com/rgrs-x/service/api/utils"
	uuid "github.com/satori/go.uuid"
)

// RepoResponse ...
type RepoResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

// Status ...
type Status string

// Unique Status Message
const (
	Created      Status = "CREATED"
	Success      Status = "SUCCESS"
	Accepted     Status = "ACCEPTED"
	Uploaded     Status = "UPLOADED"
	Deleted      Status = "DELETED"
	Unauthorized Status = "UNAUTHORIZED"
	Forbidden    Status = "FORBIDDEN"
	Liked        Status = "LIKED"

	NotFound Status = "NOT_FOUND"
	GetError Status = "GET_ERROR"
	Existed  Status = "EXISTED"

	CanNotCreate Status = "CAN_NOT_CREATE"
	CanNotUpdate Status = "CAN_NOT_UPDATE"
	CanNotDelete Status = "CAN_NOT_DELETE"
	Updated      Status = "UPDATED"

	CannotGetAll Status = "CAN_NOT_GET_ALL"
	CannotGet    Status = "CAN_NOT_GET"
)

// AsString return status as String
func (status Status) AsString() string {

	switch status {
	case CanNotCreate:
		return "Can not create"
	case CanNotUpdate:
		return "Can not update"
	case CanNotDelete:
		return "Can not delete"
	case Updated:
		return "updated"
	case CannotGetAll:
		return "can not get all"
	case Success:
		return "success"
	case Created:
		return "created"
	case Deleted:
		return "deleted"

	case CannotGet:
		return "can not get"
	case NotFound:
		return "Not Found"
	default:
		return ""
	}
}

// AsStatus return status as bool
func (status Status) AsStatus() bool {

	switch status {
	case CanNotCreate:
		return false
	case CannotGet:
		return false
	case CanNotUpdate:
		return false
	case CanNotDelete:
		return false
	case CannotGetAll:
		return false
	case Created:
		return true
	case Updated:
		return true
	case Success:
		return true
	case Deleted:
		return true
	case Accepted:
		return true
	case NotFound:
		return false
	default:
		return false
	}
}

//PostRepository is an interface can be implemented
type PostRepository interface {
	Validate(post models.Post) bool
	ValidateBlank(details ...string) bool

	Create(models.Post, string) (models.Post, Status)
	CreateIntroduction(models.Post, string) (models.Post, Status)

	GetPostDetails(string) (models.Post, Status)
	GetIntroductionDetails(string) (models.Post, Status)
	GetPost(id uuid.UUID) (u.ResultRepository, int)
	GetIntroduction(string) (RepoResponse, Status)

	GetAllPosts() []models.Post
	GetPartnerContents(id uuid.UUID) (u.ResultRepository, int)
	GetCompanyContents(id uuid.UUID) (u.ResultRepository, int)

	FetchCreator(creatorID string, post *models.Post) (u.ResultRepository, int)

	UpdatePost(post models.Post, postID string, creator string) (map[string]interface{}, int)

	DeletePost(id string, creator string) (map[string]interface{}, int)
	GetAllTags() (map[string]interface{}, int)

	//New implement Pagination
	CountContents(pagintion *models.Pagination) error
	Pagination(pagination *models.Pagination) u.ResultRepository

	/* Filter*/
	Filter(filter *models.Filter) u.ResultRepository
	CheckFilter(input string, test []string) bool
	RemoveSamePost(input string, test []models.Post) []models.Post

	/*Like Post*/
	UpdatePostLike(id uuid.UUID) (u.ResultRepository, int)
	/* Tracking */
	UpdatePostReview(id uuid.UUID) (u.ResultRepository, int)
}

//CompanyRepository is an interface can be implemented
type CompanyRepository interface {
	ValidateBlank(details ...string) bool
	Create(company *models.Company) (u.ResultRepository, int)
	GetCompanyList() (u.ResultRepository, int)
	GetMembers(company models.Company) (u.ResultRepository, int)
}
