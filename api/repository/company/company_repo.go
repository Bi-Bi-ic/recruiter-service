package company

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	u "github.com/rgrs-x/service/api/utils"
	"gorm.io/gorm"
)

// partnerStorage is a struct implementing PartnerRepository interface{}
type companyStorage struct {
	Db *gorm.DB
}

// NewCompanyRepository ... We can implement to use CompanyRepository interface{} there
func NewCompanyRepository(db *gorm.DB) repo.CompanyRepository {
	return &companyStorage{
		Db: db,
	}
}

//Check if field is blank
func (storage *companyStorage) ValidateBlank(details ...string) bool {
	for _, detail := range details {
		re := regexp.MustCompile(`(?m)^\s*$`)
		result := re.MatchString(detail)
		if result == true {
			return false
		}
	}
	return true
}

// CreateCompany ...
func (storage *companyStorage) Create(company *models.Company) (u.ResultRepository, int) {
	if result := storage.ValidateBlank(company.Name); !result {
		return u.ResultRepository{Result: []string{}, Error: errors.New("Company's Name is invalid or typos")}, http.StatusBadRequest
	}

	var workSpace models.Company
	err := storage.Db.Table("companies").Where("name = ?", company.Name).First(&workSpace).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = storage.Db.Create(&company).Error
			if err != nil {
				return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
			}
			return u.ResultRepository{Result: company}, http.StatusCreated
		}
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	return u.ResultRepository{Result: []string{}, Error: errors.New("Sorry this Company was created")}, http.StatusForbidden
}

// GetCompanyMemberList ...
func (storage *companyStorage) GetCompanyList() (u.ResultRepository, int) {
	var companies []models.Company
	// Check if Companies is exist
	err := storage.Db.Model(&models.Company{}).Select("id, name").Find(&companies).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.ResultRepository{Result: []string{}, Error: ErrCompanyNotFound}, http.StatusNotFound
		}
	}
	return u.ResultRepository{Result: companies}, http.StatusOK
}

// checkMemberFlag ...
func (storage *companyStorage) checkMemberFlag(permission string) bool {
	members := [...]string{"admin", "viewer"}
	for _, value := range members {
		if ok := strings.Compare(permission, value); ok == 0 {
			return true
		}
	}

	return false
}

// GetCompanyMember ...
func (storage *companyStorage) GetMembers(company models.Company) (u.ResultRepository, int) {
	//Check Company is exist
	err := storage.Db.Table("companies").Select("name").First(&company).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.ResultRepository{Result: []string{}, Error: ErrCompanyNotFound}, http.StatusNotFound
		}
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	// Fetch all Members from this Company
	var members []models.Member
	err = storage.Db.Table("partners").Select("id, permission").Where("company_id = ?", company.ID).Find(&members).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.ResultRepository{Result: []string{}, Error: errors.New("Company is not Owned")}, http.StatusAccepted
		}
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	if len(members) == 0 {
		return u.ResultRepository{Result: []string{}, Error: ErrNoMember}, http.StatusNotFound
	}

	for _, value := range members {
		if ok := storage.checkMemberFlag(value.Permission); ok {
			company.Members = append(company.Members, value)
		}
	}
	return u.ResultRepository{Result: company, Message: "Found all Member"}, http.StatusOK
}
