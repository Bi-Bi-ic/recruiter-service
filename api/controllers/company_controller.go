package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/company"
	u "github.com/rgrs-x/service/api/utils"
	uuid "github.com/satori/go.uuid"
)

// CreateCompany ...
func CreateCompany(c *gin.Context) {
	var workSpace models.Company
	err := c.Bind(&workSpace)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: "Invalid Request"})
	}

	companyRepository := company.NewCompanyRepository(models.GetDB())
	response, statusCode := companyRepository.Create(&workSpace)

	if statusCode == http.StatusCreated {
		c.JSON(http.StatusCreated, u.Response{Status: true, Message: "Company has been Created", Data: response.Result})
		return
	}

	c.JSON(statusCode, u.Response{Status: false, Message: response.Error.Error(), Data: response.Result})
}

// SwitchGetCompany ...
func SwitchGetCompany(c *gin.Context) {
	params := c.Request.URL.Query()
	// fmt.Println(len(params))
	if len(params) == 0 {
		GetCompanyList(c)
		return
	}

	for key, value := range params {
		queryValue := value[len(value)-1]

		switch key {
		case "company_id":
			GetMembers(c, queryValue)
			return
		}
	}

	c.JSON(http.StatusForbidden, u.Response{Status: false, Message: "Url Request is invalid", Data: []string{}})
}

// GetCompanyList ...
func GetCompanyList(c *gin.Context) {
	companyRepository := company.NewCompanyRepository(models.GetDB())
	response, statusCode := companyRepository.GetCompanyList()

	if statusCode == http.StatusOK {
		c.JSON(statusCode, u.Response{Status: true, Message: "Got All Companies", Data: response.Result})
		return
	}

	c.JSON(statusCode, u.Response{Status: false, Message: response.Error.Error(), Data: response.Result})
}

// GetMembers ...
func GetMembers(c *gin.Context, companyID string) {
	var workSpace models.Company
	var err error

	workSpace.ID, err = uuid.FromString(companyID)
	if err != nil {
		c.JSON(http.StatusForbidden, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	companyRepository := company.NewCompanyRepository(models.GetDB())
	response, statusCode := companyRepository.GetMembers(workSpace)

	if statusCode != http.StatusOK {
		c.JSON(statusCode, u.Response{Status: false, Message: response.Error.Error(), Data: response.Result})
		return
	}

	c.JSON(statusCode, u.Response{Status: true, Message: response.Message, Data: response.Result})
}
