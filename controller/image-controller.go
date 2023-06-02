package controller

import (
	"fmt"
	"net/http"
	"os"

	"gilab.com/progrmaticreviwes/golang-gin-poc/utilService"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

// image upload controller
func UploadImage(ctx *gin.Context){
	//1. Get the image data from the request body
	var inputData struct{
		Input struct{
			Arg1 struct{
				Image string `json:"image"`	
				Images [] string `json:"images"`
			} `json:"arg1"`
		} `json:"input"`
	}

	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//2. Set up the Cloudinary configuration
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_SECRET"))
	var images = inputData.Input.Arg1.Images

	//3. the cover image url is last url in the array
	if (inputData.Input.Arg1.Image != ""){
		images = append(images, inputData.Input.Arg1.Image)
	}
	//4.upload images to cloudinary and store the urls in an array
	var urls []string
	for index := range  images{
		// Upload the image to Cloudinary
		response , err := cld.Upload.Upload(ctx.Request.Context(), images[index], uploader.UploadParams{
			PublicID: utilService.PublicID(),
			Folder:   "persons",
		})
		if err != nil {
			fmt.Println("Error:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(response.SecureURL)
		urls = append(urls, response.SecureURL)
	}
	// 5. Send the url to the client
	ctx.JSON(200, gin.H{"urls": urls})
}
