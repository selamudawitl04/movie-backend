package controller

import (
	"fmt"
	"net/http"

	"gilab.com/progrmaticreviwes/golang-gin-poc/utilService"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadImage(ctx *gin.Context){
	// Get the image data from the request body
	var inputData struct{
		Input struct{
			Arg1 struct{	
				Images [] string `json:"images"`
			} `json:"arg1"`
		} `json:"input"`
	}

	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var urls []string
	// Set up the Cloudinary configuration
	cld, _ := cloudinary.NewFromParams("selamu-dawit", "349479795881572", "_PCvEGOJlFmS6wMA6z1JZ93b53o")

	var images = inputData.Input.Arg1.Images

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

	ctx.JSON(200, gin.H{"urls": urls})
}
