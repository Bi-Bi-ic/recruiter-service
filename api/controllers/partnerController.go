package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/models/message"
	repository "github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/company"
	"github.com/rgrs-x/service/api/repository/partner"
	"github.com/rgrs-x/service/api/repository/post"
	u "github.com/rgrs-x/service/api/utils"
	"github.com/rgrs-x/service/api/validator"
	uuid "github.com/satori/go.uuid"
	govalidator "gopkg.in/go-playground/validator.v9"
)

// PartnerResponse ...
func PartnerResponse(statusCode int, Data u.ResultRepository, c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	var partnerFactory factory.PartnerInfoFactory
	var avatarFactory factory.AvatarInfoFactory

	switch statusCode {
	case http.StatusOK, http.StatusCreated:
		Data.Result = partnerFactory.Create(Data.Result)
		c.JSON(statusCode, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	case http.StatusAccepted:
		Data.Result = partnerFactory.Create(Data.Result)
		c.JSON(statusCode, u.Response{Status: false, Message: Data.Error.Error(), Data: Data.Result})
		return
	case models.PartnerRequest:
		c.JSON(http.StatusOK, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	case models.PartnerAvatarUpload:
		Data.Result = avatarFactory.PartnerAvatar(Data.Result)
		c.JSON(http.StatusOK, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	case models.PartnerInfo:
		Data.Result = partnerFactory.CreateDetail(Data.Result)
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

// Get Status

/*
	--For Partner Authenciation--------------------------------------------------------------
*/

// CreatePartnerAccount ...
func CreatePartnerAccount(c *gin.Context) {
	var customer models.Partner
	var partnerFactory factory.PartnerInfoFactory
	//@ decode the request body into struct and failed if any error occur
	if ok := BindRequest(models.PartnerNormal, &customer, c); !ok {
		return
	}

	partnerSecure := validator.NewPartnerValidator()
	err := partnerSecure.Valid(customer)
	if err != nil {
		statusCode := partnerSecure.Handle(err.(govalidator.ValidationErrors))

		if statusCode != code.Ok {
			response := u.BTResponse{Status: false, Message: message.DataInvalid, Data: []string{}, Code: statusCode}
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	repo, status := partnerRepository.Create(customer)

	switch status {
	case repository.Existed:
		resp := u.BTResponse{Status: false, Message: message.EmailIsUsed, Data: []string{}, Code: code.EmailIsUsed}
		c.JSON(http.StatusForbidden, resp)

	case repository.Created:
		repo.Data = partnerFactory.Create(repo.Data)
		resp := u.BTResponse{Status: true, Message: message.PartnerCreated, Data: repo.Data, Code: code.Created}

		c.JSON(getStatusCode(status), resp)

	case repository.GetError:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		c.JSON(getStatusCode(status), resp)

	default:
		resp := u.BTResponse{Status: false, Message: message.InternalServerError, Data: []string{}, Code: code.InternalServerError}
		c.JSON(getStatusCode(""), resp)
	}

}

// AuthenticatePartner for Login API
func AuthenticatePartner(c *gin.Context) {
	var customer models.Partner
	var partnerFactory factory.PartnerInfoFactory
	//@ decode the request body into struct and failed if any error occur
	if ok := BindRequest(models.PartnerNormal, &customer, c); !ok {
		return
	}

	partnerSecure := validator.NewPartnerValidator()
	partnerSecure.SetMode(validator.Login)

	err := partnerSecure.Valid(customer)
	if err != nil {
		statusCode := partnerSecure.Handle(err.(govalidator.ValidationErrors))

		if statusCode != code.Ok {
			response := u.BTResponse{Status: false, Message: message.DataInvalid, Data: []string{}, Code: statusCode}
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	repo, status := partnerRepository.Login(customer.Email, customer.Password, customer)

	switch status {
	case repository.Success:
		repo.Data = partnerFactory.Create(repo.Data)
		resp := u.BTResponse{Status: true, Message: message.Login, Data: repo.Data, Code: code.Ok}

		c.JSON(getStatusCode(status), resp)

	case repository.NotFound:
		resp := u.BTResponse{Status: false, Message: message.EmailNotFound, Data: []string{}, Code: code.EmailNotFound}
		c.JSON(getStatusCode(status), resp)

	case repository.GetError:
		resp := u.BTResponse{Status: false, Message: message.QueryError, Data: []string{}, Code: code.QueryError}
		c.JSON(getStatusCode(status), resp)

	case repository.Unauthorized:
		resp := u.BTResponse{Status: false, Message: message.PasswordNotMatch, Data: []string{}, Code: code.PasswordError}
		c.JSON(getStatusCode(status), resp)
	}

}

// GetPartnerInfo return Partner's informations
func GetPartnerInfo(c *gin.Context) {
	partnerID := c.Writer.Header().Get("user")
	_, err := uuid.FromString(partnerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	resp, statusCode := partnerRepository.GetByID(partnerID)

	PartnerResponse(statusCode, resp, c)
}

// PublicPartnerInfo ...
func PublicPartnerInfo(c *gin.Context) {
	partnerCard := c.Param("id")
	partnerID, err := uuid.FromString(partnerCard)
	if err != nil {
		resp := u.BTResponse{Status: false, Message: message.MalFormedID, Data: []string{}, Code: code.MalFormedID}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	repo, status := partnerRepository.PublicInfo(partnerID)

	handlerStatus(repo, status, models.PartnerNormal, c)
}

// UpdateAvatarPartner ...
func UpdateAvatarPartner(c *gin.Context) {
	partnerID := c.Writer.Header().Get("user")
	_, err := uuid.FromString(partnerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	//Then execute Query
	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	resp, statusCode := partnerRepository.GetByID(partnerID)
	if statusCode != models.PartnerInfo {
		PartnerResponse(statusCode, resp, c)
	} else {

		// Then reflect to partner entity for reponse partner-able
		var customer models.Partner
		err := mapstructure.Decode(resp.Result, &customer)
		if err != nil {
			log.Println(err)
			PartnerResponse(http.StatusRequestTimeout, u.ResultRepository{Result: []string{}, Error: errors.New("Something go wrong. Please retry")}, c)
		}

		imageName, _ := UploadProfileImage(c, partnerID)

		// Whenever can not Upload image here, exit the function
		if imageName == "" {
			c.Abort()
			return
		}

		imgPath := "/api/partner/avatar/" + imageName

		/* Here if everything is Ok we can generate Partner-Avatar
		Otherwise, return error message as usual
		*/
		resp, statusCode = partnerRepository.UpdateAvatar(customer, imgPath)
		PartnerResponse(statusCode, resp, c)
	}
}

// UpdatePartnerCover ...
func UpdatePartnerCover(c *gin.Context) {
	partnerCard := c.Writer.Header().Get("user")
	partnerID, err := uuid.FromString(partnerCard)
	if err != nil {
		resp := u.BTResponse{Status: false, Message: message.MalFormedID, Data: []string{}, Code: code.MalFormedID}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	imgName, imgID := UploadProfileImage(c, partnerCard)
	if imgName == "" {
		resp := u.BTResponse{Status: false, Message: message.ImageError, Data: []string{}, Code: code.UploadError}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var partnerRepository = partner.NewPartnerRepository(models.GetDB())
	repo, status := partnerRepository.UpdateCoverImg(imgName, imgID, partnerID)

	handlerStatus(repo, status, models.CoverNormal, c)
}

// UpdatePartnerInfo ...
func UpdatePartnerInfo(c *gin.Context) {
	partnerCard := c.Writer.Header().Get("user")
	partnerID, err := uuid.FromString(partnerCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	var customer models.Partner
	customer.ID = partnerID
	if ok := BindRequest(models.PartnerNormal, &customer, c); !ok {
		return
	}

	if customer.ValidEmpty() == false {
		c.JSON(400, models.BadRequest.SendMessage())
		return
	}

	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	resp, statusCode := partnerRepository.Update(customer)

	PartnerResponse(statusCode, resp, c)
}

// GetRequestList ...
func GetRequestList(c *gin.Context) {
	companyID := c.Query("id")

	adminID := c.Writer.Header().Get("user")
	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	response, statusCode := partnerRepository.RequestCompanyList(adminID, companyID)

	PartnerResponse(statusCode, response, c)
}

// GetPartnerContents ...
func GetPartnerContents(c *gin.Context) {
	partnerCard := c.Param("id")
	partnerID, err := uuid.FromString(partnerCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	postRepository := post.NewPostRepository(models.GetDB())
	response, statusCode := postRepository.GetPartnerContents(partnerID)

	PostHandler(statusCode, response, c)
}

// JoinRequest ...
func JoinRequest(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	partnerID := c.Writer.Header().Get("user")
	_, err := uuid.FromString(partnerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	var workSpace models.Company
	err = c.Bind(&workSpace)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: company.ErrCompanyNotFound.Error(), Data: []string{}})
		return
	}

	companyRepository := company.NewCompanyRepository(models.GetDB())
	if ok := companyRepository.ValidateBlank(workSpace.Name); !ok {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: company.ErrCompanyNotFound.Error(), Data: []string{}})
		return
	}

	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	response, statusCode := partnerRepository.JoinRequest(workSpace.Name, partnerID)

	PartnerResponse(statusCode, response, c)
}

// CancelRequest ...
func CancelRequest(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	partnerID := c.Writer.Header().Get("user")
	_, err := uuid.FromString(partnerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	var company models.Company
	err = c.Bind(&company)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	response, statusCode := partnerRepository.CancelRequest(company.ID.String(), partnerID)

	PartnerResponse(statusCode, response, c)
}

// AcceptMemberRequest ...
func AcceptMemberRequest(c *gin.Context) {
	adminCard := c.Writer.Header().Get("user")
	adminID, err := uuid.FromString(adminCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	var member models.Partner
	err = c.Bind(&member)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	response, statusCode := partnerRepository.AcceptRequest(adminID, member.ID)

	PartnerResponse(statusCode, response, c)
}

// DeclineMemberRequest ...
func DeclineMemberRequest(c *gin.Context) {
	adminCard := c.Writer.Header().Get("user")
	adminID, err := uuid.FromString(adminCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	var member models.Partner
	err = c.Bind(&member)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	response, statusCode := partnerRepository.DeclineRequest(adminID, member.ID)

	PartnerResponse(statusCode, response, c)
}

// LikeMentor ...
func LikeMentor(c *gin.Context) {
	mentorID := c.Param("id")

	var partnerRepository = partner.NewPartnerRepository(models.GetDB())
	repo, status := partnerRepository.UpdateMentorLike(mentorID)

	handlerStatus(repo, status, models.PartnerNormal, c)
}
