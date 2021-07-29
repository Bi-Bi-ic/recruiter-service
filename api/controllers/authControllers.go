package controllers

import (
	"fmt"
	"net/http"

	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/models/message"
	"github.com/rgrs-x/service/api/repository/user"
	u "github.com/rgrs-x/service/api/utils"
	"github.com/rgrs-x/service/api/validator"
	uuid "github.com/satori/go.uuid"
	govalidator "gopkg.in/go-playground/validator.v9"

	"github.com/gin-gonic/gin"
)

/*
	--For General Logic----------------------------------------------------------
*/

// BindRequest ...
func BindRequest(mode models.UserMode, customer interface{}, c *gin.Context) bool {
	// Switch User mode
	switch mode {
	case models.UserNormal:
		err := c.ShouldBindJSON(&customer)
		if err != nil {
			response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
			c.JSON(http.StatusBadRequest, response)
			return false
		}
		return true

	case models.PartnerNormal:
		err := c.ShouldBindJSON(&customer)
		if err != nil {
			response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
			c.JSON(http.StatusBadRequest, response)
			return false
		}
		return true

	case models.AdminNormal:
		err := c.ShouldBindJSON(&customer)
		if err != nil {
			statusCode := models.BadRequest
			c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: statusCode.SendMessage().Message, Data: []string{}, Code: &statusCode})
			return false
		}
		return true

	case models.TimeLineNormal:
		err := c.ShouldBindJSON(&customer)
		if err != nil {
			response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
			c.JSON(http.StatusBadRequest, response)
			return false
		}
		return true

	default:
		fmt.Println("Can not find User Mode")
		return false
	}
}

/*
	--For User Authenciation----------------------------------------------------------
*/

// AuthenticateUser ...
func AuthenticateUser(c *gin.Context) {
	var customer models.User
	//@ decode the request body into struct and failed if any error occur
	if ok := BindRequest(models.UserNormal, &customer, c); !ok {
		return
	}
	userSecure := validator.NewUserValidator()
	err := userSecure.Valid(customer)
	if err != nil {
		statusCode := userSecure.Handle(err.(govalidator.ValidationErrors))

		if statusCode != code.Ok {
			response := u.BTResponse{Status: false, Message: message.DataInvalid, Data: []string{}, Code: statusCode}
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	userRepository := user.NewUserRepository(models.GetDB())
	resp, statusCode := userRepository.Login(customer.Email, customer.Password, customer)

	UserResponse(statusCode, resp, c)
}

// GetAuthUserInfo to return a info of user
func GetAuthUserInfo(c *gin.Context) {
	userID := c.Writer.Header().Get("user")
	customer := models.User{}

	userRepository := user.NewUserRepository(models.GetDB())

	/* Here if everything is Ok we can generate User-Detail
	Otherwise, return error message as usual
	*/
	resp, statusCode := userRepository.GetByID(userID, customer)
	if statusCode == http.StatusOK {
		var userFactory factory.UserInfoFactory
		resp["data"] = userFactory.CreateDetail(resp["data"])

		c.JSON(statusCode, resp)
	} else {
		c.JSON(statusCode, resp)
	}
}

// PublicUserInfo ...
func PublicUserInfo(c *gin.Context) {
	userCard := c.Param("id")
	getUserById(c, userCard)
}

func getUserById(c *gin.Context,id string) {
	customer := models.User{}

	userRepository := user.NewUserRepository(models.GetDB())
	resp, statusCode := userRepository.GetByID(id, customer)
	if statusCode == http.StatusOK {
		var userFactory factory.UserInfoFactory
		resp["data"] = userFactory.CreateDetail(resp["data"])

		c.JSON(statusCode, resp)
	} else {
		c.JSON(statusCode, resp)
	}
}

// CreateUserTimeLine ...
func CreateUserTimeLine(c *gin.Context) {
	userCard := c.Writer.Header().Get("user")
	userID, err := uuid.FromString(userCard)
	if err != nil {
		resp := u.BTResponse{Status: false, Message: message.MalFormedID, Data: []string{}, Code: code.MalFormedID}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var timeline models.TimeLine
	if ok := BindRequest(models.UserNormal, &timeline, c); !ok {
		return
	}

	var userRepository = user.NewUserRepository(models.GetDB())
	repo, status := userRepository.CreateTimeLine(timeline, userID)
	if repo.Status == false {
		handlerStatus(repo, status, models.TimeLineNormal, c)
		return
	}

	repo, _ = userRepository.GetInfo(userID)
	handlerStatus(repo, status, models.TimeLineNormal, c)
}

// UpdateUserTimeLine ...
func UpdateUserTimeLine(c *gin.Context) {
	userCard := c.Writer.Header().Get("user")
	userID, err := uuid.FromString(userCard)
	if err != nil {
		resp := u.BTResponse{Status: false, Message: message.MalFormedID, Data: []string{}, Code: code.MalFormedID}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	timelineCard := c.Param("id")
	timelineID, err := uuid.FromString(timelineCard)
	if err != nil {
		resp := u.BTResponse{Status: false, Message: message.MalFormedID, Data: []string{}, Code: code.MalFormedID}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var timeline models.TimeLine
	if ok := BindRequest(models.UserNormal, &timeline, c); !ok {
		return
	}
	timeline.ID = timelineID

	var userRepository = user.NewUserRepository(models.GetDB())
	repo, status := userRepository.UpdateTimeLine(timeline, userID)
	if repo.Status == false {
		handlerStatus(repo, status, models.TimeLineNormal, c)
		return
	}

	repo, _ = userRepository.GetInfo(userID)
	handlerStatus(repo, status, models.TimeLineNormal, c)
}

// DeleteUserTimeLine ...
func DeleteUserTimeLine(c *gin.Context) {
	userCard := c.Writer.Header().Get("user")
	userID, err := uuid.FromString(userCard)
	if err != nil {
		resp := u.BTResponse{Status: false, Message: message.MalFormedID, Data: []string{}, Code: code.MalFormedID}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	timelineCard := c.Param("id")
	timelineID, err := uuid.FromString(timelineCard)
	if err != nil {
		resp := u.BTResponse{Status: false, Message: message.MalFormedID, Data: []string{}, Code: code.MalFormedID}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var userRepository = user.NewUserRepository(models.GetDB())
	repo, status := userRepository.DeleteTimeLine(timelineID, userID)

	handlerStatus(repo, status, models.TimeLineNormal, c)
}
