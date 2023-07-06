package main
import (
	"log"

	"gilab.com/progrmaticreviwes/golang-gin-poc/controller"
	"gilab.com/progrmaticreviwes/golang-gin-poc/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
func init() {
    // Load environmental variables from .env file
    err := godotenv.Load("go_server.env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}
func main() {
	// create server
	server := gin.New()
	// middlewares
	server.Use(middlewares.Logger())
	// server.Use(middlewares.CorsMiddleware())
	// routes
	server.POST("/login", controller.Login)
	server.POST("/signup", controller.Signup)
	server.POST("/forgotPassword", controller.ForgotPassword)
	server.POST("/resetPassword", controller.ResetPassword)
	server.POST("/updateUser", controller.UpdateUser)
	server.POST("/uploadImage", controller.UploadImage)
	server.POST("/changeStatus", controller.ChangeStatus)
	server.POST("/sendMessage", controller.SendMessage)
	server.POST("CheckAPI", controller.CheckAPI)
	
		
	server.Run(":7000")
}
