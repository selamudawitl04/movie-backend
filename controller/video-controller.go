package controller

import (
	"gilab.com/progrmaticreviwes/golang-gin-poc/entity"
	"gilab.com/progrmaticreviwes/golang-gin-poc/service"
	"github.com/gin-gonic/gin"
)

type VideoController interface {
	Save(ctx *gin.Context) entity.Video
	FindAll() []entity.Video
}
type controller struct{
	service service.VideoService
}

func New(service service.VideoService) VideoController {
	return controller{
		service: service,
	}
}
// implement save
func(c controller) Save(ctx *gin.Context) entity.Video{
	var video entity.Video
	ctx.BindJSON(&video)
	c.service.Save(video)
	return video
}

func(c controller) FindAll() []entity.Video{
	return c.service.FindAll()
}


 
