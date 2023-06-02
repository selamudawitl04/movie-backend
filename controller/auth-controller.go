package controller

// imports
import (
	"context"
	"fmt"
	"net/http"

	"gilab.com/progrmaticreviwes/golang-gin-poc/utilService"
	"github.com/gin-gonic/gin"
)

// this respose object for login, signup, update, and resert password handlers
type AuthResponse struct{
	ID          string `json:"id"`
	Role string `json:"role"`
	Token string `json:"token"`
}

func deletePendingUser(email string) {
	//1.  Define the GraphQL query to execute
	var query struct {	
		DeletePendingUsers struct {
			Returning []struct{
				Email string `json:"email"`
			} `json:"returning"`
		} `graphql:"delete_pending_users(where: {email: {_eq: $email}})"`
	}
	//2. construct graphql variables
	variables := map[string]interface{}{
		"email":  email,
	}
	//3. execute the request
	err := utilService.Client().Mutate(context.Background(), &query, variables)
	if err != nil {
		fmt.Println(err.Error())
	}
}







// this function accept payload to create token and then call util service with payload
// finally send token to client
func sendToken(ctx *gin.Context, role string, response AuthResponse) {
	token, err := utilService.GetToken(response.ID, role)
	
	if err != nil {
		ctx.JSON(400, gin.H{"message": "Something went wrong"})
		return
	}
	response.Role = role
	response.Token = token
	ctx.JSON(200, response)
}

// login controller
func Login( ctx *gin.Context){
	//1. Get user data from request body
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
	//2.  Define the GraphQL query to execute
	var query struct {
		Users[] struct {
			ID          string `json:"id"`
			Email string `json:"email"`
			Password string `json:"password"` 
			Role string `json:"role"`
		
		} `graphql:"users(where: {email: {_eq: $email}})"`
	}

	//3. construct graphql variables
	variables := map[string]interface{}{
		"email":  inputUser.Input.Arg1.Email,
	}
	//4. execute the graphql query 
	err := utilService.Client().Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(400, gin.H{"error Meassagw": "wrong token"})
		return
	}

	//5. Check if the user exists and the password is correct
	if len(query.Users) > 0 && utilService.ComparePasswords(query.Users[0].Password, inputUser.Input.Arg1.Password) {
		// if the user exists, send the token with user data
		var response AuthResponse
		response.ID = query.Users[0].ID
		response.Role = query.Users[0].Role
		sendToken(ctx, query.Users[0].Role, response)
		return
	}else{
		// if there is no user the send invalid credential message
		ctx.JSON(400, gin.H{"message": "Invalid credentials"})
	}

}
// signup controller
func Signup(ctx *gin.Context) {
    //1. Get the user input from the request body
    type inputUser struct {
		Input struct{
			Arg1 struct{
				FirstName string `json:"firstName"`
				LastName string `json:"lastName"`
				Email string `json:"email"`
				Password string `json:"password"`
				Token string `json:"token"`
			} `json:"arg1"`
		} `json:"input"`
    }
	var newUser inputUser
    if err := ctx.ShouldBindJSON(&newUser); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	//2. Validate the token
	err1 := utilService.ValidateToken(newUser.Input.Arg1.Token)
	if err1 != nil {
		fmt.Println(err1.Error())
		ctx.JSON(400, gin.H{"error": err1.Error()})
		// delete token from database
		deletePendingUser(newUser.Input.Arg1.Email)
	}

	//3. Define the GraphQL mutation string
	var mutation struct {
		InsertUsers struct {	
			Returning []struct{
				ID string `json:"id"`
			} `json:"returning"`
		} `graphql:"insert_users(objects: {firstName: $firstName, lastName: $lastName, email: $email, password: $password})"`
	}

	// 3. hash password 
	password, err4 := utilService.HashPassword(newUser.Input.Arg1.Password)
	if err4 != nil {
		ctx.JSON(400, gin.H{"error": err4.Error()})
		return
	}

	//4.  construct graphql variable
	variables := map[string]interface{}{
		"firstName":  newUser.Input.Arg1.FirstName,
		"lastName":  newUser.Input.Arg1.LastName,
		"email":  newUser.Input.Arg1.Email,
		"password": password,
	}

	//5. execute the request
	err := utilService.Client().Mutate(context.Background(), &mutation, variables)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	deletePendingUser(newUser.Input.Arg1.Email)
	//6. If data stored successfuly call sendToken function with response object
	var response AuthResponse
	response.ID = mutation.InsertUsers.Returning[0].ID
	sendToken(ctx, "user", response)
}
// forgot password controller
func ForgotPassword( ctx *gin.Context){
	//1. Get user data from request body
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
    
	//2.  Define the GraphQL query to execute
	var query struct {
		User[] struct {
			Email string `json:"email"` 
		} `graphql:"users(where: {email: {_eq: $email}})"`
	}
	//3. construct graphql variables
	variables := map[string]interface{}{
		"email":  inputUser.Input.Arg1.Email,
	}
	//4 execute the request 
	err2 := utilService.Client().Query(context.Background(), &query, variables)
	if err2 != nil {
		fmt.Println(err2.Error())
		ctx.JSON(400, gin.H{"message": "Something went  wrong when quering"})
		return
	}

	//5.  Check if the user exists
	if len(query.User) == 0  {
		message := fmt.Sprintf("There is no User with email  %s ", inputUser.Input.Arg1.Email)
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
	token, err3 := utilService.ResetPasswordAndRegisterToken(inputUser.Input.Arg1.Email)
	if err3 != nil {
		ctx.JSON(400, gin.H{"error": err3.Error()})
		return
	}

	//8. construct graphql variables
	variables2 := map[string]interface{}{
		"resetToken":  token,
		"email":  inputUser.Input.Arg1.Email,
	}
	//9. execute the set password Reset token mutation
	err4 := utilService.Client().Mutate(context.Background(), &mutation, variables2)
	if err4 != nil {
		ctx.JSON(400, gin.H{"error": err4.Error()})
		return
	}
	//10. Send password reset token to user by  email
	message, err5 := utilService.SendEmail(inputUser.Input.Arg1.Email, token, "Reset your password")
	if err5 != nil {
		ctx.JSON(400, gin.H{"error": err5.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": token, "message2": message})
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
		ctx.JSON(400, gin.H{"error": err1.Error()})
		// delete token from database
		return
	}
	// //4. If the token is valid then search for the user with the token
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
			ctx.JSON(400, gin.H{"error": err4.Error()})
			return
		}
		// 1.Define the GraphQL mutation string
		var mutation struct {
			UpdateUsers struct {
				Returning []struct{
					ID string `json:"id"`
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
			ctx.JSON(400, gin.H{"error": err5.Error()})
			return
		}
		//4. send token
		ctx.JSON(200, gin.H{"message": "Password reset successfuly"})
		return
	}
	ctx.JSON(400, gin.H{"message": "Invalid credentials"})
}

func UpdateUser(ctx *gin.Context){	
	//1. Get user data from request body
	var inputUser struct {
		Input struct{
			Arg1 struct{
				FirstName string `json:"firstName"`
				LastName string `json:"lastName"`
				Email string `json:"email"`
				Password string `json:"password"`
				NewPassword string `json:"newPassword"`
			} `json:"arg1"`
		} `json:"input"`
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
		"email":  inputUser.Input.Arg1.Email,
	}
	//4. execute the request 
	err := utilService.Client().Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(400, gin.H{"error Meassag": "Something went wrong"})
		return
	}

	//5.  Check if the user exists and the password is correct
	if len(query.Users) > 0 && utilService.ComparePasswords(query.Users[0].Password, inputUser.Input.Arg1.Password) {
		// if the user exists, send the token with user data
		// change password
		password, err4 := utilService.HashPassword(inputUser.Input.Arg1.NewPassword)
		if err4 != nil {
			ctx.JSON(400, gin.H{"error": err4.Error()})
			return
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
			"password":  password,
			"firstName":  inputUser.Input.Arg1.FirstName,
			"lastName":  inputUser.Input.Arg1.LastName,
			"email":  inputUser.Input.Arg1.Email,
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
func RequestRegister( ctx *gin.Context){
	//1. Get user data from request body
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
    // 2. Create register token	
	token, err2 := utilService.ResetPasswordAndRegisterToken(inputUser.Input.Arg1.Email)
	if err2 != nil {
		ctx.JSON(400, gin.H{"error": err2.Error()})
		return
	}
	//3.  Define the GraphQL mutation string that store password reset token in database
	var mutation struct {
		InsertPendingUsers struct {	
			Returning []struct{
				Email string `json:"id"`
			} `json:"returning"`
		} `graphql:"insert_pending_users(objects: {email: $email, token: $token})"`
	}
	//4. construct graphql variables
	variables2 := map[string]interface{}{
		"token":  token,
		"email":  inputUser.Input.Arg1.Email,

	}
	//5. execute the set password Reset token mutation
	err3 := utilService.Client().Mutate(context.Background(), &mutation, variables2)
	if err3 != nil {
		ctx.JSON(400, gin.H{"error": err3.Error()})
		return
	}
	//6. Send password reset token to user by  email
	message, err5 := utilService.SendEmail(inputUser.Input.Arg1.Email, token, "Verify your email")
	if err5 != nil {
		ctx.JSON(400, gin.H{"error": err5.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": token, "message2": message})
}






























