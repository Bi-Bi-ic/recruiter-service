package user

import (
	"net/http"
	"strings"
	"time"

	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	u "github.com/rgrs-x/service/api/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userStorage is a struct implementing PartnerRepository interface{}
type userStorage struct {
	DB *gorm.DB
}

// NewUserRepository ... We can implement to use UserRepository interface{} there
func NewUserRepository(db *gorm.DB) repo.UserRepository {
	return &userStorage{
		DB: db,
	}
}

// Validate incoming email ...
func (storage *userStorage) checkEmailExist(userEmail string) repo.RepoResponse {

	var result bool
	//@ check for errors and duplicate emails
	commonDB, _ := storage.DB.DB()
	row := commonDB.QueryRow("SELECT EXISTS(SELECT email FROM users WHERE email = $1)", userEmail)
	row.Scan(&result)
	if result {
		logrus.WithFields(logrus.Fields{
			"email": userEmail,
		}).Info("Email is used !!!")

		return repo.RepoResponse{Status: false}
	}

	return repo.RepoResponse{Status: true}
}

//Create make new user
func (storage *userStorage) Create(user models.User) (result repo.RepoResponse, status repo.Status) {
	if resp := storage.checkEmailExist(user.Email); !resp.Status {
		return resp, repo.Existed
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	//Get User username
	components := strings.Split(user.Email, "@")
	user.UserName = components[0]

	Find := storage.DB.Create(&user)
	err := Find.Error
	if err != nil {
		result = repo.RepoResponse{Status: false}
		status = repo.GetError

		return
	}

	//Now generate JWT token here
	user.GenerateToken()

	//@ delete password
	user.Password = ""

	result = repo.RepoResponse{Status: true, Data: user}
	status = repo.Created
	return
}

//Login for user has already account
func (storage *userStorage) Login(email, password string, user models.User) (u.ResultRepository, int) {
	err := storage.DB.Table("users").Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.ResultRepository{Result: []string{}, Error: repo.ErrEmailNotFound}, http.StatusNotFound
		}
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.ResultRepository{Result: []string{}, Error: repo.ErrLogin}, http.StatusUnauthorized
	}
	//@ Worked! Logged In
	user.Password = ""

	// Now make JWT token here
	user.GenerateToken()

	return u.ResultRepository{Result: user, Message: "Logged in"}, http.StatusOK
}

//GetByID fetch user informations
func (storage *userStorage) GetByID(id string, user models.User) (map[string]interface{}, int) {
	storage.DB.Preload("Skills").Preload("Languages").Preload("OldJobs").Preload("Educations").Preload("TimeLines", "delete_at IS NULL").Where("id = ?", id).First(&user)
	if user.Email == "" { //User not found!
		return u.Message(false, "User does not exist"), http.StatusNotFound
	}

	user.Password = ""

	response := u.Message(true, "Found User")
	response["data"] = user

	return response, http.StatusOK
}

// GetInfo ...
func (storage *userStorage) GetInfo(userID uuid.UUID) (result repo.RepoResponse, status repo.Status) {
	user, err := storage.getUserInfo(userID)
	if err != nil {
		result = repo.RepoResponse{Status: false}
		if err != gorm.ErrRecordNotFound {
			status = repo.NotFound
			return
		}
		status = repo.GetError
		return
	}

	result = repo.RepoResponse{Status: true, Data: user}
	status = repo.Success

	return
}

// get user info ...
func (storage *userStorage) getUserInfo(userID uuid.UUID) (user models.User, err error) {
	user.ID = userID
	queryStmt := storage.DB.Model(&user).Preload("Skills").Preload("Languages").Preload("OldJobs").Preload("Educations").Preload("TimeLines", "delete_at IS NULL").First(&user)
	err = queryStmt.Error
	if err != nil {
		return
	}

	user.Password = ""
	return
}

// PublicInfo ...
func (storage *userStorage) PublicInfo(userID uuid.UUID) (result repo.RepoResponse, status repo.Status) {
	var user models.User
	user.ID = userID

	queryStmt := storage.DB.Model(&user).Preload("Skills").Preload("Languages").Preload("OldJobs").Preload("Educations").Preload("TimeLines", "delete_at IS NULL").First(&user)
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

	user.Password = ""

	result = repo.RepoResponse{Status: true, Data: user}
	status = repo.Success

	return
}

//Update update user 's informations if they want
func (storage *userStorage) Update(user models.User) (map[string]interface{}, int) {

	if user.MailContact != "" {
		if !strings.Contains(user.MailContact, "@") {
			return u.Message(false, "Mail Contact address is invalid"), http.StatusBadRequest
		}
	}
	storage.DB.Model(&user).Omit("id", "email", "create_at", "delete_at", "avatar").Updates(user)

	response := u.Message(true, "Updated user infomation")

	//Handle some info not wanted
	user.Token = ""
	user.RefreshToken = ""
	user.Password = ""

	response["data"] = user
	return response, http.StatusOK
}

//UpdateAvatar set image can be stored on database
func (storage *userStorage) UpdateAvatar(user models.User, imageName string) (map[string]interface{}, int) {

	err := storage.DB.Model(&user).Updates(map[string]interface{}{"avatar": imageName}).Error
	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	response := u.Message(true, "Updated Avatar")
	response["data"] = user

	return response, http.StatusOK
}

// UpdateCover ...
func (storage *userStorage) UpdateCoverImg(imgName, imgID string, userID uuid.UUID) (map[string]interface{}, int) {
	var user models.User
	var factory factory.FileInfoFactoty

	user.ID = userID

	errGetUserInfo := storage.DB.Table("users").Where("id = ?", userID).First(&user).Error

	if errGetUserInfo != nil {
		return u.Message(false, errGetUserInfo.Error()), http.StatusNotFound
	}

	cover := factory.Create(imgName, imgID)

	queryStmt := storage.DB.Model(&user).Updates(map[string]interface{}{"cover": cover.Link})
	err := queryStmt.Error
	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	response := u.Message(true, "Updated Avatar")
	response["data"] = user

	return response, http.StatusOK
}

// CreateTimeLine ...
func (storage *userStorage) CreateTimeLine(timeline models.TimeLine, userID uuid.UUID) (result repo.RepoResponse, status repo.Status) {
	var user models.User
	user.ID = userID

	err := storage.DB.Model(&user).Updates(map[string]interface{}{"update_at": time.Now().Unix()}).Association("TimeLines").Append([]models.TimeLine{timeline})
	if err != nil {
		result = repo.RepoResponse{Status: false}
		if err == gorm.ErrRecordNotFound {
			status = repo.NotFound

			return
		}
		status = repo.GetError
		return
	}

	result = repo.RepoResponse{Status: true, Data: timeline}
	status = repo.Created

	return
}

// UpdateTimeLine ...
func (storage *userStorage) UpdateTimeLine(timeline models.TimeLine, userID uuid.UUID) (result repo.RepoResponse, status repo.Status) {
	var user models.User
	user.ID = userID

	err := storage.DB.Model(&user).Updates(map[string]interface{}{"update_at": time.Now().Unix()}).Association("TimeLines").Append([]models.TimeLine{timeline})

	if err != nil {
		result = repo.RepoResponse{Status: false}
		if err == gorm.ErrRecordNotFound {
			status = repo.NotFound

			return
		}
		status = repo.GetError
		return
	}

	result = repo.RepoResponse{Status: true, Data: timeline}
	status = repo.Success

	return
}

// DeleteTimeLine ...
func (storage *userStorage) DeleteTimeLine(timelineID uuid.UUID, userID uuid.UUID) (result repo.RepoResponse, status repo.Status) {
	var timeline models.TimeLine
	timeline.ID = timelineID

	queryStmt := storage.DB.Model(&timeline).Where("user_id = ? AND delete_at IS NULL", userID)

	var count int64
	queryStmt.Count(&count)
	if count == 0 {
		result = repo.RepoResponse{Status: false}
		status = repo.NotFound

		return
	}

	queryStmt = storage.DB.Model(&timeline).Where("user_id = ? AND delete_at IS NULL", userID).Updates(map[string]interface{}{"delete_at": time.Now()})
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

	result = repo.RepoResponse{Status: true, Data: timeline}
	status = repo.Deleted

	return
}
