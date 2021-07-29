package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/factory"
	u "github.com/rgrs-x/service/api/utils"
)

func UploadFile(c *gin.Context) {
	//Get Id from header
	tempID := c.Writer.Header().Get("user")
	imageName, id := UploadProfileImage(c, tempID)

	// Whenever can not Upload image here, exit the function
	resp := make(map[string]interface{})

	var factory factory.FileInfoFactoty

	if imageName == "" {
		u := u.Message(false, "Mail Contact address is invalid")
		resp = u
		resp["data"] = nil
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	u := u.Message(true, "Uploaded")

	resp = u

	resp["data"] = factory.Create(imageName, id)

	c.JSON(200, resp)

}
