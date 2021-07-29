package partner

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	company_repo "github.com/rgrs-x/service/api/repository/company"
	u "github.com/rgrs-x/service/api/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// partnerStorage is a struct implementing PartnerRepository interface{}
type partnerStorage struct {
	DB *gorm.DB
}

// NewPartnerRepository ... We can implement to use PartnerRepository interface{} there
func NewPartnerRepository(db *gorm.DB) repo.PartnerRepository {
	return &partnerStorage{
		DB: db,
	}
}

// Validate incoming email ...
func (storage *partnerStorage) checkEmailExist(partnerEmail string) repo.RepoResponse {

	var result bool
	//@ check for errors and duplicate emails
	commonDB, _ := storage.DB.DB()
	row := commonDB.QueryRow("SELECT EXISTS(SELECT email FROM partners WHERE email = $1)", partnerEmail)
	row.Scan(&result)
	if result {
		logrus.WithFields(logrus.Fields{
			"email": partnerEmail,
		}).Info("Email is used !!!")

		return repo.RepoResponse{Status: false}
	}

	return repo.RepoResponse{Status: true}
}

//Create make new Partner
func (storage *partnerStorage) Create(partner models.Partner) (result repo.RepoResponse, status repo.Status) {
	if resp := storage.checkEmailExist(partner.Email); !resp.Status {
		return resp, repo.Existed
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(partner.Password), bcrypt.DefaultCost)
	partner.Password = string(hashedPassword)

	//Get Partner username
	components := strings.Split(partner.Email, "@")
	partner.UserName = components[0]

	Find := storage.DB.Create(&partner)
	err := Find.Error
	if err != nil {
		result = repo.RepoResponse{Status: false}
		status = repo.GetError

		return
	}

	// Now Generate JWT token here
	partner.GenerateToken(models.PartnerNormal)

	//@ delete password
	partner.Password = ""

	result = repo.RepoResponse{Status: true, Data: partner}
	status = repo.Created
	return
}

//Login for partner has already account
func (storage *partnerStorage) Login(email, password string, partner models.Partner) (result repo.RepoResponse, status repo.Status) {
	err := storage.DB.Table("partners").Where("email = ?", email).First(&partner).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result = repo.RepoResponse{Status: false}
			status = repo.NotFound
			return
		}
		result = repo.RepoResponse{Status: false}
		status = repo.GetError
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(partner.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		result = repo.RepoResponse{Status: false}
		status = repo.Unauthorized
		return
	}
	//@ Worked! Logged In
	partner.Password = ""

	// Now Generate JWT token here
	partner.GenerateToken(models.PartnerNormal)

	result = repo.RepoResponse{Status: true, Data: partner}
	status = repo.Success
	return
}

//GetByID fetch partner informations
func (storage *partnerStorage) GetByID(partnerID string) (u.ResultRepository, int) {
	partner := storage.GetInfo(partnerID)
	if partner.ID == uuid.Nil {
		return u.ResultRepository{Result: []string{}, Error: ErrPartnerNotFound}, http.StatusForbidden
	}

	err := storage.CountPosts(&partner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			partner.PostAvailable = "0"
		} else {
			return u.ResultRepository{Result: []string{}, Error: ErrPartnerContents}, http.StatusRequestTimeout
		}
	}

	err = storage.CountPostsExpired(&partner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			partner.PostExpired = "0"
		} else {
			return u.ResultRepository{Result: []string{}, Error: ErrPartnerContents}, http.StatusRequestTimeout
		}
	}

	return u.ResultRepository{Result: partner, Message: "Found Partner"}, models.PartnerInfo
}

// GetDataByID ...
func (storage *partnerStorage) GetDataByID(partnerID string) (partner models.Partner, status repo.Status) {
	queryStmt := storage.DB.Where("id = ?", partnerID)
	find := queryStmt.First(&partner)
	err := find.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = repo.NotFound
			return
		}

		status = repo.CannotGet
		return
	}

	err = storage.CountPosts(&partner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			partner.PostAvailable = "0"
		} else {
			status = repo.CannotGet
			return
		}
	}

	err = storage.CountPostsExpired(&partner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			partner.PostExpired = "0"
		} else {
			status = repo.CannotGet
			return
		}
	}

	status = repo.Success
	return
}

// GetInfo ...
func (storage *partnerStorage) GetInfo(partnerID string) models.Partner {
	var partner models.Partner
	storage.DB.Table("partners").Where("id = ?", partnerID).First(&partner)

	partner.Password = ""
	return partner
}

// PublicInfo ...
func (storage *partnerStorage) PublicInfo(partnerID uuid.UUID) (result repo.RepoResponse, status repo.Status) {
	var partner models.Partner
	partner.ID = partnerID

	queryStmt := storage.DB.Model(&partner).First(&partner)
	err := queryStmt.Error
	if err != nil {
		result = repo.RepoResponse{Status: false}
		if err == gorm.ErrRecordNotFound {
			status = repo.NotFound
			return
		}
		status = repo.GetError
		return
	}

	err = storage.CountPosts(&partner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			partner.PostAvailable = "0"
		} else {
			result = repo.RepoResponse{Status: false}
			status = repo.GetError

			return
		}
	}

	err = storage.CountPostsExpired(&partner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			partner.PostExpired = "0"
		} else {
			result = repo.RepoResponse{Status: false}
			status = repo.GetError

			return
		}
	}

	partner.Password = ""

	result = repo.RepoResponse{Status: true, Data: partner}
	status = repo.Success

	return
}

// Get Partner Profile ...
func (storage *partnerStorage) getPartnerProfile(partnerID string) (partner models.Partner, err error) {
	queryStmt := storage.DB.Table("partners").Where("id = ?", partnerID).First(&partner)
	err = queryStmt.Error
	if err != nil {
		return
	}

	return
}

//Update update partner 's informations if they want
func (storage *partnerStorage) Update(partner models.Partner) (u.ResultRepository, int) {

	if partner.PostAvailable != "" || partner.PostExpired != "" {
		return u.ResultRepository{Result: []string{}, Error: errors.New("Invalid Request")}, http.StatusBadRequest
	}

	queryStatement := storage.DB.Model(&partner).Omit(
		"id",
		"create_at", "delete_at",
		"avatar", "company_id",
		"name",
		"permission",
		"user_name",
		"email").
		Updates(partner)

	err := queryStatement.Error
	if err != nil {
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	rowAffected := queryStatement.RowsAffected
	if rowAffected == 0 {
		return u.ResultRepository{Result: []string{}, Error: ErrPartnerNotFound}, http.StatusForbidden
	}

	partnerCard := storage.GetInfo(partner.ID.String())
	if partnerCard.ID == uuid.Nil {
		return u.ResultRepository{Result: []string{}, Error: ErrPartnerNotFound}, http.StatusForbidden
	}

	//Handle some info not wanted
	partnerCard.Token = ""
	partnerCard.RefreshToken = ""

	return u.ResultRepository{Result: partnerCard, Message: "Updated Partner Information"}, models.PartnerInfo
}

//UpdateAvatar set image can be stored on database
func (storage *partnerStorage) UpdateAvatar(partner models.Partner, imageName string) (u.ResultRepository, int) {

	err := storage.DB.Model(&partner).
		Where("id = ?", partner.ID).
		Updates(map[string]interface{}{
			"update_at": time.Now().Unix(),
			"avatar":    imageName,
		}).Error
	if err != nil {
		return u.ResultRepository{Result: []string{}, Error: err}, http.StatusBadRequest
	}

	return u.ResultRepository{Result: partner, Message: "Updated Partner's Avatar"}, models.PartnerAvatarUpload
}

// UpdateMentorLike ...
func (storage *partnerStorage) UpdateMentorLike(mentorID string) (result repo.RepoResponse, status repo.Status) {
	mentor, err := storage.getPartnerProfile(mentorID)
	if err != nil {
		status = repo.NotFound
		return
	}

	mentor.TotalLike++
	queryStmt := storage.DB.Model(&mentor).Updates(map[string]interface{}{
		"update_at":  time.Now().Unix(),
		"total_like": mentor.TotalLike,
	})
	err = queryStmt.Error
	if err != nil {
		status = repo.GetError
		return
	}

	// delete Password
	mentor.Password = ""

	result = repo.RepoResponse{Status: true, Data: mentor}
	status = repo.Liked
	return
}

// UpdateCover ...
func (storage *partnerStorage) UpdateCoverImg(imgName, imgID string, partnerID uuid.UUID) (result repo.RepoResponse, status repo.Status) {
	var partner models.Partner
	var factory factory.FileInfoFactoty

	partner.ID = partnerID
	cover := factory.Create(imgName, imgID)

	queryStmt := storage.DB.Model(&partner).Updates(map[string]interface{}{
		"update_at": time.Now().Unix(),
		"cover":     cover.Link,
	})
	err := queryStmt.Error
	if err != nil {
		result = repo.RepoResponse{Status: false}

		if err == gorm.ErrRecordNotFound {
			status = repo.NotFound
			return
		}
		status = repo.GetError
		return
	}
	result = repo.RepoResponse{Status: true, Data: cover}
	status = repo.Uploaded

	return
}

// CountPosts ...
func (storage *partnerStorage) CountPosts(partner *models.Partner) error {
	var count int64

	queryStatement := storage.DB.Model(&models.Post{}).
		Where("creator_id = ? AND delete_at IS NULL", partner.ID).
		Count(&count)

	partner.PostAvailable = strconv.FormatInt(count, 10)

	result := queryStatement.Find(&models.Post{}).Error

	return result
}

// CountPostExpired ...
func (storage *partnerStorage) CountPostsExpired(partner *models.Partner) error {
	var count int64
	queryStatement := storage.DB.Model(&models.Post{}).
		Where("creator_id = ? AND delete_at IS NOT NULL", partner.ID).
		Count(&count)

	partner.PostExpired = strconv.FormatInt(count, 10)

	result := queryStatement.Find(&models.Post{}).Error

	return result
}

/* For Company
---------------------------------------------------------------------------------------------
*/

// CheckOwnerCompany ...
func (storage *partnerStorage) checkCompanyExist(companyName string) (result bool) {
	commonDB, _ := storage.DB.DB()
	row := commonDB.QueryRow("SELECT EXISTS(SELECT name FROM partners WHERE name = $1)", companyName)
	row.Scan(&result)
	if result {
		logrus.WithFields(logrus.Fields{
			"company": companyName,
		}).Info("Company is owned !!!")
	}

	return
}

// GetCompanyInfo ...
func (storage *partnerStorage) GetCompanyInfo(companyName string) models.Company {
	var company models.Company
	storage.DB.Table("companies").Where("name = ?", companyName).First(&company)

	return company
}

// RequestCompanyList ...
func (storage *partnerStorage) RequestCompanyList(adminID, companyID string) (u.ResultRepository, int) {
	var admin models.Partner
	var err error
	admin.ID, err = uuid.FromString(adminID)
	if err != nil {
		return u.ResultRepository{Result: []string{}, Error: repo.ErrMalformedID}, http.StatusForbidden
	}

	if ok := storage.CheckAdminFlag(&admin); !ok {
		return u.ResultRepository{Result: []string{}, Error: ErrAdminFlag}, http.StatusForbidden
	}
	if ok := strings.Compare(admin.CompanyID, companyID); ok != 0 {
		return u.ResultRepository{Result: []string{}, Error: ErrAdminFlag}, http.StatusForbidden
	}

	var guests []models.Partner
	err = storage.DB.Table("partners").Where("company_id = ? AND permission like 'wait'", companyID).Find(&guests).Error
	if err != nil {
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}
	if len(guests) == 0 {
		return u.ResultRepository{Result: []string{}, Error: errors.New("Nothing request Joining Company to now")}, http.StatusNotFound
	}

	for profileCard := range guests {
		guests[profileCard].Password = ""
	}

	return u.ResultRepository{Result: guests, Message: "Got All Requests"}, models.PartnerRequest
}

// CheckAdminFlag ...
func (storage *partnerStorage) CheckAdminFlag(admin *models.Partner) bool {
	err := storage.DB.Table("partners").Select("company_id ,name, permission").First(&admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		panic(err)
	}
	if admin.Permission != "admin" {
		return false
	}
	return true
}

// CheckMemberFlag ...
func (storage *partnerStorage) CheckMemberFlag(permission string) bool {
	members := [...]string{"admin", "viewer"}
	for _, value := range members {
		if ok := strings.Compare(permission, value); ok == 0 {
			return true
		}
	}

	return false
}

// IsSameCompany ...
func (storage *partnerStorage) IsSameCompany(companyName string, member *models.Partner) bool {
	err := storage.DB.Table("partners").First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		panic(err)
	}
	if member.Name != companyName {
		return false
	}
	return true
}

// AcceptRequest ...
func (storage *partnerStorage) AcceptRequest(adminID, memberID uuid.UUID) (u.ResultRepository, int) {
	var admin, partner models.Partner
	admin.ID = adminID
	partner.ID = memberID

	ok := storage.CheckAdminFlag(&admin)
	if !ok {
		return u.ResultRepository{Result: []string{}, Error: ErrAdminNotFound}, http.StatusForbidden
	}
	ok = storage.IsSameCompany(admin.Name, &partner)
	if !ok {
		return u.ResultRepository{Result: []string{}, Error: ErrNotRequest}, http.StatusForbidden
	}

	if ok := storage.CheckMemberFlag(partner.Permission); ok {
		return u.ResultRepository{Result: []string{}, Error: ErrMemberFlag}, http.StatusForbidden
	}

	if ok := strings.Compare(partner.Permission, "wait"); ok != 0 {
		return u.ResultRepository{Result: []string{}, Error: ErrNotRequest}, http.StatusForbidden
	}

	storage.DB.Model(&partner).Updates(map[string]interface{}{"permission": "viewer"})

	var member models.Member
	member.ID = partner.ID.String()
	member.Permission = partner.Permission

	return u.ResultRepository{Result: member, Message: "Accepted Member's Request"}, models.PartnerRequest
}

// DeclineRequest ...
func (storage *partnerStorage) DeclineRequest(adminID, memberID uuid.UUID) (u.ResultRepository, int) {
	var admin, partner models.Partner
	admin.ID = adminID
	partner.ID = memberID

	ok := storage.CheckAdminFlag(&admin)
	if !ok {
		return u.ResultRepository{Result: []string{}, Error: ErrAdminNotFound}, http.StatusForbidden
	}
	ok = storage.IsSameCompany(admin.Name, &partner)
	if !ok {
		return u.ResultRepository{Result: []string{}, Error: ErrNotRequest}, http.StatusForbidden
	}

	if ok := storage.CheckMemberFlag(partner.Permission); ok {
		return u.ResultRepository{Result: []string{}, Error: ErrMemberFlag}, http.StatusForbidden
	}

	if ok := strings.Compare(partner.Permission, "wait"); ok != 0 {
		return u.ResultRepository{Result: []string{}, Error: ErrNotRequest}, http.StatusForbidden
	}

	storage.DB.Model(&partner).Updates(map[string]interface{}{"permission": "deny"})

	var member models.Member
	member.ID = partner.ID.String()
	member.Permission = partner.Permission

	return u.ResultRepository{Result: member, Message: "Declined Member's Request"}, models.PartnerRequest
}

// JoiningRequest ...
func (storage *partnerStorage) JoinRequest(companyName, partnerID string) (u.ResultRepository, int) {
	company := storage.GetCompanyInfo(companyName)
	if company.ID == uuid.Nil {
		return u.ResultRepository{Result: []string{}, Error: company_repo.ErrCompanyNotFound}, http.StatusForbidden
	}

	partner := storage.GetInfo(partnerID)
	if partner.ID == uuid.Nil {
		return u.ResultRepository{Result: []string{}, Error: ErrPartnerNotFound}, http.StatusForbidden
	}

	if ok := strings.Compare(partner.Permission, "wait"); ok == 0 {
		return u.ResultRepository{Result: []string{}, Error: errors.New("Your Request has been Sent to a Company")}, http.StatusForbidden
	}

	if ok := storage.CheckMemberFlag(partner.Permission); ok {
		return u.ResultRepository{Result: []string{}, Error: ErrMemberRequest}, http.StatusForbidden
	}

	if company.CreatorID == "" {
		storage.DB.Model(&company).Updates(map[string]interface{}{"creator_id": partner.ID.String()})
		storage.DB.Model(&partner).Where("id = ?", partnerID).Updates(map[string]interface{}{"company_id": company.ID, "name": company.Name, "permission": "admin"})

		var member models.Member
		member.ID = partner.ID.String()
		member.Permission = partner.Permission

		return u.ResultRepository{Result: member, Message: "Now You Owned this Company"}, models.PartnerRequest
	}

	storage.DB.Model(&partner).Where("id = ?", partnerID).Updates(map[string]interface{}{"company_id": company.ID, "name": company.Name, "permission": "wait"})

	var member models.Member
	member.ID = partner.ID.String()
	member.Permission = partner.Permission

	return u.ResultRepository{Result: member, Message: "Request Sent"}, models.PartnerRequest
}

// CancleRequest ...
func (storage *partnerStorage) CancelRequest(companyID, partnerID string) (u.ResultRepository, int) {
	company := storage.GetCompanyInfo(companyID)
	if company.ID == uuid.Nil {
		return u.ResultRepository{Result: []string{}, Error: company_repo.ErrCompanyNotFound}, http.StatusForbidden
	}

	partner := storage.GetInfo(partnerID)
	if partner.ID == uuid.Nil {
		return u.ResultRepository{Result: []string{}, Error: ErrPartnerNotFound}, http.StatusForbidden
	}

	if ok := storage.CheckMemberFlag(partner.Permission); ok {
		return u.ResultRepository{Result: []string{}, Error: ErrMemberRequest}, http.StatusForbidden
	}

	if ok := strings.Compare(partner.Permission, "deny"); ok == 0 {
		return u.ResultRepository{Result: []string{}, Error: ErrNotBelongCompany}, http.StatusForbidden
	}

	storage.DB.Model(&partner).Where("id = ?", partnerID).Updates(map[string]interface{}{"permission": "deny"})

	var member models.Member
	member.ID = partner.ID.String()
	member.Permission = partner.Permission

	return u.ResultRepository{Result: member, Message: "Request Canceled"}, models.PartnerRequest
}
