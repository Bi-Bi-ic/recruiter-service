package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/repository/post"
	u "github.com/rgrs-x/service/api/utils"
)

//Pagination ...
func Pagination(c *gin.Context) {
	code := http.StatusOK

	pagination := GeneratePaginationRequest(c)

	response := PaginationQuery(c, pagination)

	if !response.Status {
		code = http.StatusBadRequest
	}

	c.JSON(code, response)
}

//GeneratePaginationRequest ...
func GeneratePaginationRequest(c *gin.Context) *models.Pagination {
	// default limit, page & sort parameter
	limit := 10
	offset := 0
	sort := models.Lastest

	query := c.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "offset":
			offset, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = models.Sort(queryValue)
			break
		}
	}

	return &models.Pagination{Limit: limit, Offset: offset, Sort: sort}
}

//PaginationQuery ...
func PaginationQuery(c *gin.Context, pagination *models.Pagination) u.Response {
	postRepository := post.NewPostRepository(models.GetDB())

	err := postRepository.CountContents(pagination)
	if err != nil {
		return u.Response{Status: false}
	}

	if pagination.Offset >= pagination.TotalContents || pagination.Offset < 0 {
		return u.Response{Status: false, Message: "Out of Range Contents !!!", Data: []string{}}
	}
	operationResult := postRepository.Pagination(pagination)

	if operationResult.Error != nil {
		return u.Response{Status: false, Message: operationResult.Error.Error(), Data: operationResult.Result}
	}

	var data = operationResult.Result.(*models.Pagination)
	return u.Response{Status: true, Data: data}
}
