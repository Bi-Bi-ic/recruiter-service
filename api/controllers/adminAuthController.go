package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/models/message"
	repository "github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/admin"
	u "github.com/rgrs-x/service/api/utils"
)

// AdminSignIn ...
func AdminSignIn(c *gin.Context) {
	var customer models.Admin
	//@ decode the request body into struct and failed if any error occur
	if ok := BindRequest(models.AdminNormal, &customer, c); !ok {
		return
	}
	fmt.Println(customer.Email)
	if valid := customer.IsEmpty(); valid {
		c.JSON(http.StatusBadRequest, u.BTResponse{Status: false, Message: message.DataInvalid, Data: []string{}, Code: code.DataIsEmpty})
		return
	}

	adminRepository := admin.New(models.GetDB())

	repo, status := adminRepository.Create(customer.Email, customer.Password)

	switch status {
	case repository.Success:
		c.JSON(http.StatusAccepted, u.BTResponse{Status: false, Message: message.OK, Data: repo.Data, Code: code.Ok})
	}

}
