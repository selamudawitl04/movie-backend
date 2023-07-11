package controller

// imports
import (
	"context"
	"fmt"
	"net/http"
	"bytes"
	"html/template"
	"gilab.com/progrmaticreviwes/golang-gin-poc/utilService"
	"github.com/gin-gonic/gin"
)

// this respose object for login, signup, update, and resert password handlers
type AuthResponse struct{
	ID          string `json:"id"`
	Role string `json:"role"`
	Token string `json:"token"`
}

// to send email to user
type EmailDataToken struct {
    Link string
	Header string
}



// this function accept payload to create token and then call util service with payload
// finally send token to client
func sendToken(ctx *gin.Context, role string, response AuthResponse) {
	token, err := utilService.GetToken(response.ID, role)

	if err != nil {
		fmt.Println(err.Error(), "when sending token ")

		ctx.JSON(400, gin.H{"message": "Something went wrong when creating token"})
		return
	}
	response.Role = role
	response.Token = token
	ctx.JSON(200, response)
}

// login controller
func Login( ctx *gin.Context){

	var input struct {
		Email string `json:"email"`
		Password string `json:"password"`
		
	}
	fmt.Println("login controller", input)
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var query struct {
		Users[] struct {
			ID          string `json:"id"`
			Email string `json:"email"`
			Password string `json:"password"`
			Role string `json:"role"`

		} `graphql:"users(where: {email: {_eq: $email}})"`
	}

	variables := map[string]interface{}{
		"email":  input.Email,
	}
	err := utilService.Client().Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println(err.Error(), "when querying user to login")
		ctx.JSON(400, gin.H{"message": "there is no user with this email"})
		return
	}
	if len(query.Users) > 0 && utilService.ComparePasswords(query.Users[0].Password, input.Password) {
		var response AuthResponse
		response.ID = query.Users[0].ID
		response.Role = query.Users[0].Role
		sendToken(ctx, query.Users[0].Role, response)
		return
	}else{
		ctx.JSON(400, gin.H{"message": "Invalid credentials"})
		return
	}

	

}
// signup controller
func Signup(ctx *gin.Context) {
    //1. Get the user input from the request body
    type inputUser struct {
		FirstName string `json:"firstName"`
		LastName string `json:"lastName"`
		Email string `json:"email"`
		Password string `json:"password"`
    }

	var newUser inputUser
    if err := ctx.ShouldBindJSON(&newUser); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	fmt.Printf("%+v\n", newUser, "new user coming")

	//3. Define the GraphQL mutation string
	var mutation struct {
		InsertUsers struct {	
			Returning []struct{
				ID string `json:"id"`
			} `json:"returning"`
		} `graphql:"insert_users(objects: {firstName: $firstName, lastName: $lastName, email: $email, password: $password})"`
	}

	// 3. hash password 
	password, err4 := utilService.HashPassword(newUser.Password)
	if err4 != nil {
		fmt.Println(err4.Error(), "when hashing password")
		ctx.JSON(400, gin.H{"error": err4.Error()})
		return
	}
	//4.  construct graphql variable
	variables := map[string]interface{}{
		"firstName":  newUser.FirstName,
		"lastName":  newUser.LastName,
		"email":  newUser.Email,
		"password": password,
	}
	//5. execute the request
	err := utilService.Client().Mutate(context.Background(), &mutation, variables)
	if err != nil {
		fmt.Println(err.Error(), "when excuting signup  mutation")

		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//6. If data stored successfuly call sendToken function with response object
	var response AuthResponse
	response.ID = mutation.InsertUsers.Returning[0].ID
	sendToken(ctx, "user", response)
}

// check api controller
func CheckAPI( ctx *gin.Context){
	var input struct {
		FirstName string `json:"firstName"`
		LastName string `json:"lastName"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input.FirstName)
	fmt.Println(input.LastName)
	ctx.JSON(200, gin.H{"message": "API is working"})
}

// forgot password controller
func ForgotPassword( ctx *gin.Context){
	
	fmt.Println("forgot password controller")
	// 1. Get the user input from the request body
	var inputUser struct {
		Email string `json:"email"`	
	}
	
	if err1 := ctx.ShouldBindJSON(&inputUser); err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}

	fmt.Println(inputUser.Email, "email comes here")
    
	//2.  Define the GraphQL query to execute
	var query struct {
		User[] struct {
			Email string `json:"email"` 
		} `graphql:"users(where: {email: {_eq: $email}})"`
	}
	//3. construct graphql variables
	variables := map[string]interface{}{
		"email":  inputUser.Email,
	}
	fmt.Println(inputUser.Email, "email comes here")
	//4 execute the request 
	err2 := utilService.Client().Query(context.Background(), &query, variables)
	if err2 != nil {
		fmt.Println(err2.Error(), "when quering user to forgot password")
		ctx.JSON(400, gin.H{"message": "Something went  wrong when quering"})
		return
	}

	//5.  Check if the user exists
	if len(query.User) == 0  {
		message := fmt.Sprintf("There is no User with email  %s ", inputUser.Email)
		fmt.Println(message)
		ctx.JSON(400, gin.H{"message": message})
		return
	}
	//6.  Define the GraphQL mutation string that store password reset token in database
	var mutation struct {
		UpdateUsers struct {
			Returning []struct{
				ID string `json:"id"`	
			} `json:"returning"`
		} `graphql:"update_users(where: {email: {_eq: $email}}, _set: {resetToken: $resetToken})"`
	}
	//7. create password reset token
	token, err3 := utilService.ResetPasswordAndRegisterToken(inputUser.Email)
	if err3 != nil {
		ctx.JSON(400, gin.H{"error": err3.Error()})
		return
	}

	fmt.Println(token)
	//8. construct graphql variables
	variables2 := map[string]interface{}{
		"resetToken":  token,
		"email":  inputUser.Email,
	}
	//9. execute the set password Reset token mutation
	err4 := utilService.Client().Mutate(context.Background(), &mutation, variables2)
	if err4 != nil {
		ctx.JSON(400, gin.H{"error": err4.Error()})
		return
	}
	//10. Send password reset token to user by  email

	
	// Create the reset URL with the token
	resetURL := "http://localhost:3000/auth/resetPassword/" + token
	t, _ := template.ParseFiles("template.html")
	var body bytes.Buffer   
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Reset password from Solflx \n%s\n\n", mimeHeaders)))
	 // Define the email data

	
	emailData := EmailDataToken{
        Link: resetURL,
		Header:"Reset your password with above link" ,
    }
	t.Execute(&body, emailData)
 
	message, err5 := utilService.SendEmail(inputUser.Email, body)
	if err5 != nil {
		fmt.Println("There is error when sending email")
		ctx.JSON(400, gin.H{"error": err5.Error()})
		return
	}
	fmt.Println(message)
	ctx.JSON(200, gin.H{"message": "email is sent"})
}
// Reset Password controller
func ResetPassword(ctx *gin.Context){

	//1. get user new password  data from request body
	var inputUser struct {
		Input struct{
			Arg1 struct{
				Password string `json:"password"`
				Token string `json:"token"`
			} `json:"arg1"`
		} `json:"input"`
	}

	if err := ctx.ShouldBindJSON(&inputUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	//3. Parse the token
	err1 := utilService.ValidateToken(inputUser.Input.Arg1.Token)
	if err1 != nil {
		fmt.Println(err1.Error())
		ctx.JSON(200, gin.H{"message": "invalid token"})
		// delete token from database
		return
	}
	//4. If the token is valid then search for the user with the token
	var query struct {
		Users[] struct {
			ID          string `json:"id"`
			Email string `json:"email"`
			ResetToken string `json:"resetToken"`
		} `graphql:"users(where: {resetToken: {_eq: $resetToken}})"`
	}
	//5. construct graphql variables
	variables := map[string]interface{}{
		"resetToken":  inputUser.Input.Arg1.Token,
	}
	//6. execute the request
	err2 := utilService.Client().Query(context.Background(), &query, variables)
	if err2 != nil {
		fmt.Println(err2.Error())
		ctx.JSON(400, gin.H{"message": "Something went wrong"})
		return
	}
	//7. If the user exists then update the password
	if len(query.Users) > 0 && query.Users[0].ResetToken == inputUser.Input.Arg1.Token {
		// if the user exists, send the token with user data
		// change password
		password, err4 := utilService.HashPassword(inputUser.Input.Arg1.Password)
		if err4 != nil {
			ctx.JSON(400, gin.H{"message": err4.Error()})
			return
		}
		// 1.Define the GraphQL mutation string
		var mutation struct {
			UpdateUsers struct {
				Returning []struct{
					ID string `json:"id"`
					Role string `json:"role"`

				} `json:"returning"`
			} `graphql:"update_users(where: {email: {_eq: $email}}, _set: {password: $password, resetToken: $resetToken})"`
		}
		fmt.Println(query.Users[0].Email)

		//2. set variable
		variables2 := map[string]interface{}{
			"password":  password,
			"resetToken":  "",
			"email":  query.Users[0].Email,
		}
		//3. execute the request
		err5 := utilService.Client().Mutate(context.Background(), &mutation, variables2)
		if err5 != nil {
			ctx.JSON(400, gin.H{"message": err5.Error()})
			return
		}
		var response AuthResponse
		response.ID = mutation.UpdateUsers.Returning[0].ID
		response.Role = mutation.UpdateUsers.Returning[0].Role
		fmt.Println(response.ID, response.Role)
		sendToken(ctx, response.Role, response)

	}else{

		ctx.JSON(400, gin.H{"message": "Invalid credentials"})
	}
}

func UpdateUser(ctx *gin.Context){	
	//1. Get user data from request body
	var inputUser struct {
		
		FirstName string `json:"firstName"`
		LastName string `json:"lastName"`
		Email string `json:"email"`
		Password string `json:"password"`
		NewPassword string `json:"newPassword"`
	
	}
	
	if err := ctx.ShouldBindJSON(&inputUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//2.  Define the GraphQL query to execute
	var query struct {
		Users[] struct {
			ID          string `json:"id"`
			Email string `json:"email"`
			Password string `json:"password"` 
			Role string `json:"role"`
		
		} `graphql:"users(where: {email: {_eq: $email}})"`
	}

	//3. variables
	variables := map[string]interface{}{
		"email":  inputUser.Email,
	}
	//4. execute the request 
	err := utilService.Client().Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(400, gin.H{"error Meassag": "Something went wrong"})
		return
	}

	//5.  Check if the user exists and the password is correct
	if len(query.Users) > 0 && utilService.ComparePasswords(query.Users[0].Password, inputUser.Password) {
		// if the user exists, send the token with user data
		// change password
		var newPassword = query.Users[0].Password
		// check if there is new password then hash the new password
		if(inputUser.NewPassword != ""){
			password, err4 := utilService.HashPassword(inputUser.NewPassword)
			newPassword = password
			if err4 != nil {
				ctx.JSON(400, gin.H{"error": err4.Error()})
				return
			}
		}
		// 1.Define the GraphQL mutation string
		var mutation struct {
			UpdateUsers struct {
				Returning []struct{
					ID string `json:"id"`
					Role string `json:"role"`
					Email string `json:"email"`
				} `json:"returning"`
			} `graphql:"update_users(where: {email: {_eq: $email}}, _set: {password: $password, firstName: $firstName, lastName: $lastName})"`
		}
		//2. set variable
		
		variables2 := map[string]interface{}{
			"password":  newPassword,
			"firstName":  inputUser.FirstName,
			"lastName":  inputUser.LastName,
			"email":  inputUser.Email,
		}
		//3. execute the request
		err5 := utilService.Client().Mutate(context.Background(), &mutation, variables2)
		if err5 != nil {
			ctx.JSON(400, gin.H{"error": err5.Error()})
			return
		}
		//4. send token
		var response AuthResponse
		response.Role = mutation.UpdateUsers.Returning[0].Role
		response.ID = mutation.UpdateUsers.Returning[0].ID
		sendToken(ctx, query.Users[0].Role, response)
		return
	}
	ctx.JSON(400, gin.H{"message": "Invalid credentials"})
}	































