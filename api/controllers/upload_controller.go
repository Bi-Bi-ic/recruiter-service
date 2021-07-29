package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/models"
	"golang.org/x/crypto/bcrypt"
)

// GetMD5Hash generate unquie hash each time
func GetMD5Hash(text string, userID string) string {
	SumHash := text + userID

	hash, err := bcrypt.GenerateFromPassword([]byte(SumHash), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

// UploadImages for multi images
func UploadImages(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file[]"]

	for _, file := range files {
		log.Println("Input:" + file.Filename)

		var myimage models.Image
		image, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}

		ext := path.Ext(file.Filename)
		outfile := file.Filename[0:len(file.Filename)-len(ext)] + ".jpg"
		file.Filename = outfile

		byteContainer, err := ioutil.ReadAll(image)
		myimage.Source = byteContainer
		fmt.Printf("size:%d \n", len(byteContainer))
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))

}

// UploadOneImage ...
/*for clients who want Upload to set Avatar
  Or they just want to update one image */
func UploadProfileImage(c *gin.Context, imgPath string) (string, string) {

	//We get file from request
	file, errFile := c.FormFile("file")

	//When recieve no file or file can not gotten we a handle message here
	if errFile != nil {
		return "", ""
	}
	log.Println("Input:" + file.Filename)

	// //init the Image struct to store to db
	// var myimage models.Image

	image, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}
	ext := path.Ext(file.Filename)
	outfile := file.Filename[0:len(file.Filename)-len(ext)] + ".jpeg"
	file.Filename = outfile

	//hmm... let's get FileName and Extension
	FileNameWithoutExtension := strings.TrimSuffix(file.Filename, path.Ext(file.Filename))
	FileNameExtension := path.Ext(file.Filename)

	//Now we will get hashMD5 filename
	hashFileName := GetMD5Hash(FileNameWithoutExtension, imgPath)

	//then combine them
	filename := hashFileName + FileNameExtension
	id := hashFileName

	//we can see how much bytes of the file
	byteContainer, err := ioutil.ReadAll(image)

	// Push WorkerMessage to UploadPool for Worker.
	UploadPool <- WorkerMessage{FileName: filename, Size: byteContainer, Updated: true}

	return filename, id
}

// WorkerMessage show Task for Workers
type WorkerMessage struct {
	FileName string
	Size     []byte
	Updated  bool
}

// UploadPool where Store Tasks to Workers
var UploadPool chan WorkerMessage

// InitWorker is a Function initing to run n worker in n goroutine
func InitWorker(UploadPool chan WorkerMessage) {
	for i := 0; i < 10; i++ {
		go UploadImageTaskConsumer(UploadPool)
	}
}

// ImageAnalysis show Details Images
func ImageAnalysis(filename string, size []byte) {
	fmt.Println("Output:" + filename)
	fmt.Printf("size:%d \n", len(size))
}

// UploadToDatabase Update Images to Database Storage
func UploadToDatabase(filename string, size []byte, finished chan bool) {

	models.GetDB().Create(&models.Image{ID: filename, Source: size})

	finished <- true
}

// UploadImageTaskConsumer Handle Tasks for Workers
func UploadImageTaskConsumer(UploadPool chan WorkerMessage) {
	for {
		workerMessage := <-UploadPool
		finished := make(chan bool)
		if workerMessage.Updated == true {
			go ImageAnalysis(workerMessage.FileName, workerMessage.Size)
		}
		go UploadToDatabase(workerMessage.FileName, workerMessage.Size, finished)
		<-finished
	}
}
