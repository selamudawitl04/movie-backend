package controller

import (
	"fmt"
	"gilab.com/progrmaticreviwes/golang-gin-poc/utilService"
	"github.com/gin-gonic/gin"
	"bytes"
	"html/template"
)

type EmailDataMessage struct {
    Message string
	Header string
}

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
	// create email body
	var body bytes.Buffer   
	t, _ := template.ParseFiles("replyTemplate.html")

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Reply from Solflx \n%s\n\n", mimeHeaders)))
	var newHeader = fmt.Sprintf("%s to: %s", inputData.Subject, inputData.Email)
	fmt.Println(newHeader ,"from sending email")
		// Define the email data
	emailData := EmailDataMessage{
		Message: inputData.Message,
		Header: newHeader,
	}
	t.Execute(&body, emailData)
	//10. Send message  to user by  email
	message, err5 := utilService.SendEmail(inputData.Email, body)
	if err5 != nil {
		fmt.Println("There is error when sending email", err5)
		ctx.JSON(400, gin.H{"error": err5.Error()})
		return
	}
	fmt.Println(message)
	ctx.JSON(200, gin.H{"message": message})
	// 5. Send the url to the client
}


