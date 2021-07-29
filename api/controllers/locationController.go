package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/models/message"
	"github.com/rgrs-x/service/api/repository/location"
	u "github.com/rgrs-x/service/api/utils"
)

// FindLocation ...
func FindLocation(c *gin.Context) {
	locationID := c.Param("id")
	_, err := strconv.Atoi(locationID)
	if err != nil {
		resp := u.BTResponse{Status: false, Message: message.MalFormedID, Data: []string{}, Code: code.MalFormedID}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	locationRepository := location.NewLocationRepository(models.GetDB())
	repo, status := locationRepository.FindAddress(locationID)

	handlerStatus(repo, status, models.LocationNormal, c)
}
