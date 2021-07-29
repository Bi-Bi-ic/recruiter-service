package controllers

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/models"
)

// Render will show image if existed
func Render(c *gin.Context) {
	filename := c.Param("name")

	img := &models.Image{}

	img.ID = filename
	models.GetDB().Table("images").Where("id = ?", img.ID).First(img)

	fmt.Printf("size:%d \n", len(img.Source))

	//imgSource, _, _ := image.Decode(bytes.NewReader(img.Source))

	reader := bytes.NewReader(img.Source)
	c.Header("Content-Type", "image/jpeg") // <-- set the content-type header
	io.CopyBuffer(c.Writer, reader, img.Source)
}
