package location

import (
	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// locationStorage is a struct implementing PartnerRepository interface{}
type locationStorage struct {
	Db *gorm.DB
}

// NewLocationRepository ... We can implement to use PartnerRepository interface{} there
func NewLocationRepository(db *gorm.DB) repo.LocationRepository {
	return &locationStorage{
		Db: db,
	}
}

// GetProvinceList ...
func (storage *locationStorage) GetProvinceList() (result repo.RepoResponse, status repo.Status) {
	// Check if List of Provinces is initialized
	var provinces models.LocationList
	queryStmt := storage.Db.Where("length(code) = 2").Find(&provinces)
	err := queryStmt.Error
	if err != nil {
		result = repo.RepoResponse{Status: false}
		status = repo.GetError
		return
	}

	result = repo.RepoResponse{Status: true, Data: provinces}
	status = repo.Success
	return
}

// GetDistrictList ...
func (storage *locationStorage) GetDistrictList(provinceID int) (result repo.RepoResponse, status repo.Status) {
	var coordinates models.LocationList
	queryStmt := storage.Db.Where("parent_code = ?", provinceID).Find(&coordinates)
	err := queryStmt.Error
	if err != nil {
		result = repo.RepoResponse{Status: false}
		status = repo.GetError
		return
	}

	if len(coordinates) == 0 {
		result = repo.RepoResponse{Status: false}
		status = repo.NotFound
		return
	}

	result = repo.RepoResponse{Status: true, Data: coordinates}
	status = repo.Success
	return
}

// GetWardList ...
func (storage *locationStorage) GetWardList(cityID, districtID int) (result repo.RepoResponse, status repo.Status) {
	result = storage.checkDistrictExist(cityID, districtID)
	if !result.Status {
		status = repo.NotFound
		return
	}

	var coordinates models.LocationList
	queryStmt := storage.Db.Where("parent_code = ?", districtID).Find(&coordinates)
	err := queryStmt.Error
	if err != nil {
		result = repo.RepoResponse{Status: false}
		status = repo.GetError
		return
	}

	result = repo.RepoResponse{Status: true, Data: coordinates}
	status = repo.Success
	return
}

// checkCity
func (storage *locationStorage) checkDistrictExist(cityID, districtID int) repo.RepoResponse {

	var result bool
	//@ check for errors and duplicate emails
	commonDB, _ := storage.Db.DB()
	row := commonDB.QueryRow("SELECT EXISTS(SELECT * FROM locations WHERE parent_code = $1 AND code = $2)", cityID, districtID)
	row.Scan(&result)
	if !result {
		logrus.WithFields(logrus.Fields{
			"cityID":     cityID,
			"districtID": districtID,
		}).Info("Location Not Found !!!")

		return repo.RepoResponse{Status: false}
	}

	return repo.RepoResponse{Status: true}
}

// FindAddress ...
func (storage *locationStorage) FindAddress(locationID string) (result repo.RepoResponse, status repo.Status) {
	var coordinates models.LocationList
	queryStmt := storage.Db.Table("locations").Where("parent_code = ?", locationID).Find(&coordinates)
	err := queryStmt.Error
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

	if len(coordinates) == 0 {
		result = repo.RepoResponse{Status: false}
		status = repo.NotFound
		return
	}

	result = repo.RepoResponse{Status: true, Data: coordinates}
	status = repo.Success
	return
}
