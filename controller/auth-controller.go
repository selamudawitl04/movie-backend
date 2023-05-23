package controller

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"gilab.com/progrmaticreviwes/golang-gin-poc/utilService"
	"github.com/gin-gonic/gin"
)

type AuthResponse struct{
	ID          string `json:"id"`
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Image string `json:"image"`
	Token string `json:"token"`
}

func sendToken(ctx *gin.Context, role string, response AuthResponse) {
	token, err := utilService.GetToken(response.ID, role)
	
	if err != nil {
		ctx.JSON(400, gin.H{"message": "Something went wrong"})
		return
	}
	response.Token = token
	ctx.JSON(200, response)
}

func Login( ctx *gin.Context){
	// Get user data from request body
	var inputUser struct {
		Input struct{
			Arg1 struct{
				Email string `json:"email"`
				Password string `json:"password"`
			} `json:"arg1"`
		} `json:"input"`
	}
	
	if err := ctx.ShouldBindJSON(&inputUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//  Define the GraphQL query to execute
	var query struct {
		Users[] struct {
			ID          string `json:"id"`
			Email string `json:"email"`
			FirstName string `json:"firstName"`
			LastName string `json:"lastName"`
			Password string `json:"password"` 
			Image string `json:"image"` 

		} `graphql:"users(where: {email: {_eq: $email}})"`
	}
	// variables
	variables := map[string]interface{}{
		"email":  inputUser.Input.Arg1.Email,
	}
	// execute the request 
	err := utilService.Client().Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(400, gin.H{"error Meassagw": "wrong token"})
		return
	}

	//  Check if the user exists and the password is correct
	if len(query.Users) > 0 && query.Users[0].Password == inputUser.Input.Arg1.Password {
		// if the user exists, send the token with user data

		var response AuthResponse
		response.Email = query.Users[0].Email
		response.FirstName = query.Users[0].FirstName
		response.LastName = query.Users[0].LastName
		response.ID = query.Users[0].ID
		response.Image = query.Users[0].Image
		sendToken(ctx, "user", response)
		return
	}
	ctx.JSON(400, gin.H{"message": "Invalid credentials"})

}
func Signup(ctx *gin.Context) {
    // Get the user input from the request body
    type inputUser struct {
		Input struct{
			Arg1 struct{
				FirstName string `json:"firstName"`
				LastName string `json:"lastName"`
				Email string `json:"email"`
				Password string `json:"password"`
			} `json:"arg1"`
		} `json:"input"`
    }

	var newUser inputUser
    if err := ctx.ShouldBindJSON(&newUser); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Define the GraphQL mutation string
	var mutation struct {
		InsertUsers struct {	
			Returning []struct{
				ID string `json:"id"`
				FirstName string `json:"firstName"`
				Email string `json:"Email"`
			} `json:"returning"`
		} `graphql:"insert_users(objects: {firstName: $firstName, lastName: $lastName, email: $email, password: $password})"`
	}

	// set variable
	variables := map[string]interface{}{
		"firstName":  newUser.Input.Arg1.FirstName,
		"lastName":  newUser.Input.Arg1.LastName,
		"email":  newUser.Input.Arg1.Email,
		"password":  newUser.Input.Arg1.Password,
	}
	// execute the request
	err := utilService.Client().Mutate(context.Background(), &mutation, variables)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// send token
	var response AuthResponse
	response.Email = newUser.Input.Arg1.Email
	response.FirstName = newUser.Input.Arg1.FirstName
	response.LastName = newUser.Input.Arg1.LastName
	response.ID = mutation.InsertUsers.Returning[0].ID
	response.Image = ""

	sendToken(ctx, "user", response)
	// sendToken(ctx, mutation.InsertUsers.Returning[0].ID, "admin")
}


func ForgotPassword( ctx *gin.Context){
	// Get user data from request body
	var inputUser struct {
		Input struct{
			Arg1 struct{
				Email string `json:"email"`
			} `json:"arg1"`
		} `json:"input"`
	}
	
	if err1 := ctx.ShouldBindJSON(&inputUser); err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
    
	//  Define the GraphQL query to execute
	var query struct {
		User[] struct {
			Email string `json:"email"` 
		} `graphql:"user(where: {email: {_eq: $email}})"`
	}
	// variables
	variables := map[string]interface{}{
		"email":  inputUser.Input.Arg1.Email,
	}
	// execute the request 
	err2 := utilService.Client().Query(context.Background(), &query, variables)
	if err2 != nil {
		fmt.Println(err2.Error())
		ctx.JSON(400, gin.H{"message": "Something went wrong"})
		return
	}

	//  Check if the user exists
	if len(query.User) == 0  {
		message := fmt.Sprintf("There is no User with email  %s ", inputUser.Input.Arg1.Email)
		ctx.JSON(400, gin.H{"message": message})
		return
	}

	// Generate a random 32-byte token
    tokenBytes := make([]byte, 32)
    _, err3 := rand.Read(tokenBytes)
    if err3 != nil {
		ctx.JSON(400, gin.H{"message": "Something went wrong"})
        return
    }
    token := fmt.Sprintf("%x", tokenBytes)
    // Calculate the expiration date (1 hour from now)
    expiration := time.Now().Add(time.Hour * 1)

	// Define the GraphQL mutation string
	var mutation struct {
		UpdateUser struct {
			AffectedRows int `json:"affected_rows"`
		} `graphql:"update_user(where: {email: {_eq: $email}}, _set: {resetToken: $resetToken, resetTokenExpiry: $resetTokenExpiry})"`
	}
	// set variable
	variables2 := map[string]interface{}{
		"resetToken":  token,
		"resetTokenExpiry":  expiration,
		"email":  inputUser.Input.Arg1.Email,
	}


	// execute the request
	err4 := utilService.Client().Mutate(context.Background(), &mutation, variables2)
	if err4 != nil {
		ctx.JSON(400, gin.H{"error": err4.Error()})
		return
	}

	// Send password reset token user by  email

	message, err5 := utilService.SendEmail(inputUser.Input.Arg1.Email, token)
	if err5 != nil {
		ctx.JSON(400, gin.H{"error": err5.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": message})

}

func ResetPassword(ctx *gin.Context){	
}	
func ChangePassword(ctx *gin.Context){	
}	






























