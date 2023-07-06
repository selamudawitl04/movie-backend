package controller

import (
	"fmt"
	"gilab.com/progrmaticreviwes/golang-gin-poc/utilService"
	"github.com/gin-gonic/gin"
)

// image upload controller
func SendMessage(ctx *gin.Context){
	//1. Get the image data from the request body
	var inputData struct{
		
		Message string `json:"message"`	
		Email string `json:"email"`
		Subject string `json:"subject"`
	}

	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//10. Send message  to user by  email
	message, err5 := utilService.SendMessageEmail(inputData.Email, inputData.Message,  inputData.Subject)
	if err5 != nil {
		fmt.Println("There is error when sending email", err5)
		ctx.JSON(400, gin.H{"error": err5.Error()})
		return
	}
	fmt.Println(message)
	ctx.JSON(200, gin.H{"message": message})
	// 5. Send the url to the client
}
