package repository

import (
	"github.com/rgrs-x/service/api/models"
	u "github.com/rgrs-x/service/api/utils"
	uuid "github.com/satori/go.uuid"
)

//PartnerRepository is an interface can be implemented
type PartnerRepository interface {
	Create(partner models.Partner) (RepoResponse, Status)
	Login(email, password string, partner models.Partner) (RepoResponse, Status)
	GetByID(partnerID string) (u.ResultRepository, int)
	GetDataByID(string) (models.Partner, Status)
	GetInfo(partnerID string) models.Partner
	PublicInfo(uuid.UUID) (RepoResponse, Status)

	Update(partner models.Partner) (u.ResultRepository, int)
	UpdateAvatar(partner models.Partner, imageName string) (u.ResultRepository, int)
	UpdateMentorLike(string) (RepoResponse, Status)
	UpdateCoverImg(string, string, uuid.UUID) (RepoResponse, Status)

	// Using for Contents
	CountPosts(partner *models.Partner) error
	CountPostsExpired(partner *models.Partner) error

	/* Using Company
	--------------------------------------------------------------------------------
	*/
	GetCompanyInfo(companyName string) models.Company
	RequestCompanyList(adminID, companyID string) (u.ResultRepository, int)

	CheckAdminFlag(admin *models.Partner) bool
	CheckMemberFlag(permission string) bool
	IsSameCompany(companyName string, member *models.Partner) bool

	/* Admin
	--------------------------------------------------------------------------------
	*/
	AcceptRequest(adminID, memberID uuid.UUID) (u.ResultRepository, int)
	DeclineRequest(adminID, memberID uuid.UUID) (u.ResultRepository, int)

	/* Partner Standby
	--------------------------------------------------------------------------------
	*/
	JoinRequest(companyName, partnerID string) (u.ResultRepository, int)
	CancelRequest(companyID, partnerID string) (u.ResultRepository, int)
}
