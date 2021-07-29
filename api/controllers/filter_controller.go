package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/repository/post"
	u "github.com/rgrs-x/service/api/utils"
)

// List of query ...
const (
	filter_position = "position"
	filter_jobkind  = "job_kind"
	filter_district = "district"
)

// Filter ...
func Filter(c *gin.Context) {
	statusCode := http.StatusOK

	filter := GenerateFilteringRequest(c)

	response := FilteringQuery(c, filter)

	if !response.Status {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, response)
}

// GenerateFilteringRequest ...
func GenerateFilteringRequest(c *gin.Context) *models.Filter {
	var position, jobKind, district []string
	query := c.Request.URL.Query()
	fmt.Println(query)

	for key, value := range query {
		queryValue := value

		switch key {
		case filter_position:
			position = append(position, queryValue...)
			break
		case filter_jobkind:
			jobKind = append(jobKind, queryValue...)
			break

		case filter_district:
			district = append(district, queryValue...)
			break
		}
	}
	return &models.Filter{Position: position, JobKind: jobKind, District: district}
}

// FilteringQuery ...
func FilteringQuery(c *gin.Context, fitler *models.Filter) u.Response {
	postRepository := post.NewPostRepository(models.GetDB())

	result := postRepository.Filter(fitler)
	if result.Error != nil {
		return u.Response{Status: false, Message: result.Error.Error()}
	}
	var data = result.Result.(*models.Filter)
	return u.Response{Status: true, Message: result.Message, Data: data}
}
