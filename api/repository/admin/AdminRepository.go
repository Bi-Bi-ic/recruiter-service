package admin

import (
	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type adminStorage struct {
	Db *gorm.DB
}

// New ... We can implement to use AdminRepository interface{} there
func New(db *gorm.DB) repo.AdminRepository {
	return &adminStorage{
		Db: db,
	}
}

func (storage *adminStorage) Create(email string, password string) (repo.RepoResponse, repo.Status) {

	var current = &models.Admin{}

	if email != "root.huc.admin.bt@gmail.com" || password != "1234567" {
		return repo.RepoResponse{Status: false, Data: []string{}}, repo.GetError
	}

	err := storage.Db.Table("admins").Where("email = ?", email).First(&current).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			current.Email = email
			storage.Db.Create(&current)
		} else {
			return repo.RepoResponse{Status: false, Data: []string{}}, repo.GetError
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return repo.RepoResponse{Status: false, Data: []string{}}, repo.GetError
	}
	//@ Worked! Logged In
	current.Password = ""

	current.GenerateToken()
	return repo.RepoResponse{Status: false, Data: current}, repo.Success
}
