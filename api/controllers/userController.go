package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/models/message"
	"github.com/rgrs-x/service/api/repository/user"
	user_f "github.com/rgrs-x/service/api/repository/user"
	u "github.com/rgrs-x/service/api/utils"
	"github.com/rgrs-x/service/api/validator"
	uuid "github.com/satori/go.uuid"
	govalidator "gopkg.in/go-playground/validator.v9"
)

// UserResponse ...
func UserResponse(statusCode int, Data u.ResultRepository, c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	var userFactory factory.UserInfoFactory
	var avatarFactory factory.AvatarInfoFactory

	switch statusCode {
	case http.StatusOK, http.StatusCreated:
		Data.Result = userFactory.Create(Data.Result)
		c.JSON(statusCode, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	case http.StatusAccepted:
		Data.Result = userFactory.Create(Data.Result)
		c.JSON(statusCode, u.Response{Status: false, Message: Data.Error.Error(), Data: Data.Result})
		return
	case models.PartnerRequest:
		c.JSON(http.StatusOK, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	case models.PartnerAvatarUpload:
		Data.Result = avatarFactory.UserAvatar(Data.Result)
		c.JSON(http.StatusOK, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	case models.PartnerInfo:
		Data.Result = userFactory.CreateDetail(Data.Result)
		c.JSON(http.StatusOK, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	case models.ErrValidate:
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: Data.Error.Error(), Data: Data.Result, Code: &Data.Code})
		return
	default:
		c.JSON(statusCode, u.Response{Status: false, Message: Data.Error.Error(), Data: Data.Result})
		return
	}
}

/*
	-------------------------------------------------------------------------
*/

// CreateUserAccount ...
func CreateUserAccount(c *gin.Context) {
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
	repo, status := userRepository.Create(customer)

	handlerStatus(repo, status, models.UserNormal, c)
	return
}

// UpdateAvatarUser ...
func UpdateAvatarUser(c *gin.Context) {
	//Get Id from header
	tempID := c.Writer.Header().Get("user")

	//init new User struct
	customer := models.User{}

	//We convert id string to id uuid for store in UserID
	customer.ID, _ = uuid.FromString(tempID)

	//Then execute Query
	userRepository := user_f.NewUserRepository(models.GetDB())
	resp, statusCode := userRepository.GetByID(tempID, customer)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, resp)
	} else {

		// Then reflect to partner entity for reponse partner-able
		err := mapstructure.Decode(resp["data"], &customer)
		if err != nil {
			log.Println(err)
			response := u.Message(false, "Something go wrong. Please retry")
			c.JSON(http.StatusRequestTimeout, response)
			c.Abort()
		}

		imageName, _ := UploadProfileImage(c, tempID)

		// Whenever can not Upload image here, exit the function
		if imageName == "" {
			c.Abort()
			return
		}

		imgPath := "/api/user/avatar/" + imageName

		/* Here if everything is Ok we can generate User-Avatar
		Otherwise, return error message as usual
		*/
		resp, statusCode = userRepository.UpdateAvatar(customer, imgPath)
		if statusCode != http.StatusOK {
			c.JSON(statusCode, resp)
			c.Abort()
		}
		//Create avatarable's data
		var userFactory factory.UserInfoFactory
		resp["data"] = userFactory.CreateDetail(resp["data"])

		c.JSON(statusCode, resp)
	}
}

// UpdateUserCover ...
func UpdateUserCover(c *gin.Context) {
	userCard := c.Writer.Header().Get("user")
	userID, err := uuid.FromString(userCard)
	if err != nil {
		resp := u.BTResponse{Status: false, Message: message.MalFormedID, Data: []string{}, Code: code.MalFormedID}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	imgName, imgID := UploadProfileImage(c, userCard)
	if imgName == "" {
		resp := u.BTResponse{Status: false, Message: message.ImageError, Data: []string{}, Code: code.UploadError}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var userRepository = user.NewUserRepository(models.GetDB())
	resp, statusCode := userRepository.UpdateCoverImg(imgName, imgID, userID)

	if statusCode != http.StatusOK {
		c.JSON(statusCode, resp)
		c.Abort()
	}

	//Create avatarable's data
	var userFactory factory.UserInfoFactory
	resp["data"] = userFactory.CreateDetail(resp["data"])

	c.JSON(statusCode, resp)
}

// UpdateUserInfo ...
func UpdateUserInfo(c *gin.Context) {
	userID := c.Writer.Header().Get("user")

	customer := models.User{}

	// fmt.Println(tempID)
	customer.ID, _ = uuid.FromString(userID)

	if ok := BindRequest(models.UserNormal, &customer, c); !ok {
		return
	}

	if customer.ValidEmpty() == false {
		c.JSON(400, models.BadRequest.SendMessage())
		return
	}

	userRepository := user_f.NewUserRepository(models.GetDB())
	// Check if user exists
	var user models.User
	resp, statusCode := userRepository.GetByID(userID, user)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, resp)
		c.Abort()
	}

	// Then reflect to partner entity for reponse user-able
	err := mapstructure.Decode(resp["data"], &user)
	if err != nil {
		log.Println(err)
		response := u.Message(false, "Something go wrong. Please retry")
		c.JSON(http.StatusRequestTimeout, response)
		c.Abort()
	}

	// Evrything is OK ... Going to fetch user-infor and update new user-information
	customer.CreateAt = user.CreateAt
	customer.Avatar = user.Avatar
	customer.Email = user.Email

	/* Here if everything is Ok we can generate User-Detail
	Otherwise, return error message as usual
	*/
	resp, statusCode = userRepository.Update(customer)
	if statusCode == http.StatusOK {

		var userFactory factory.UserInfoFactory
		resp["data"] = userFactory.CreateDetail(resp["data"])

		c.JSON(statusCode, resp)
	} else {
		c.JSON(statusCode, resp)
	}
}
