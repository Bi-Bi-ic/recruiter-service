package repository

import (
	"github.com/rgrs-x/service/api/models"
	u "github.com/rgrs-x/service/api/utils"
	uuid "github.com/satori/go.uuid"
)

//UserRepository is an interface can be implemented
type UserRepository interface {
	Create(user models.User) (RepoResponse, Status)
	Login(email, password string, user models.User) (u.ResultRepository, int)
	GetByID(id string, user models.User) (map[string]interface{}, int)

	GetInfo(uuid.UUID) (RepoResponse, Status)
	PublicInfo(uuid.UUID) (RepoResponse, Status)
	Update(user models.User) (map[string]interface{}, int)
	UpdateAvatar(user models.User, imageName string) (map[string]interface{}, int)
	UpdateCoverImg(string, string, uuid.UUID) (map[string]interface{}, int)

	CreateTimeLine(models.TimeLine, uuid.UUID) (RepoResponse, Status)
	UpdateTimeLine(models.TimeLine, uuid.UUID) (RepoResponse, Status)
	DeleteTimeLine(uuid.UUID, uuid.UUID) (RepoResponse, Status)
}
