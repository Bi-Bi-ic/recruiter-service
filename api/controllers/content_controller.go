package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/post"
	uuid "github.com/satori/go.uuid"

	u "github.com/rgrs-x/service/api/utils"
)

// GetContents is return all contents did popular
func GetContents(c *gin.Context) {

	var factory factory.PostInfoFactoty

	popular := c.Query("popular")
	println(popular)
	if popular == "1" {

		postReps := GetListPostWithPopular(c)

		if postReps == nil {
			resps := u.Message(false, "Can not get All Posts")
			resps["data"] = []models.Post{}
			c.JSON(http.StatusOK, resps)
			c.Abort()
			return
		}

		//Then assign to make a List of Posts
		reps := u.Message(true, "Found all Posts")
		reps["data"] = postReps

		reps["data"] = factory.NewListPost(reps["data"])

		c.JSON(http.StatusOK, reps)
		c.Abort()
		return
	}
	c.JSON(http.StatusAccepted, u.Message(false, "Not found with popular"))
}

// GetCompanyContents ...
func GetCompanyContents(c *gin.Context) {
	companyCard := c.Param("id")
	companyID, err := uuid.FromString(companyCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	postRepository := post.NewPostRepository(models.GetDB())
	response, statusCode := postRepository.GetCompanyContents(companyID)

	PostHandler(statusCode, response, c)

}

// GetAllTags in development ...
func GetAllTags(c *gin.Context) {
	popular := c.Query("popular")

	//check ok
	println(popular)
	if popular == "1" {
		postRepository := post.NewPostRepository(models.GetDB())
		response, statusCode := postRepository.GetAllTags()

		c.JSON(statusCode, response)
		c.Abort()
		return
	}

	c.JSON(http.StatusAccepted, u.Message(false, "Not found with popular"))
}

// GetListPostWithPopular ...
func GetListPostWithPopular(c *gin.Context) []models.Post {
	postStorage := post.NewPostRepository(models.GetDB())
	contents := postStorage.GetAllPosts()
	if len(contents) < 1 {
		return nil
	}

	return contents
}
