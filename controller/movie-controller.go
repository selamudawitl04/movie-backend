package controller

import (
	"context"
	"net/http"

	"gilab.com/progrmaticreviwes/golang-gin-poc/utilService"
	"github.com/gin-gonic/gin"
)

// this is function will called from hasura events when new tickets is added
// if seats are finshed the movie status is changed to closed
func ChangeStatus( ctx *gin.Context){
	// Get user data from request body
	var inputData struct {
		New struct{
			Seat_number int `json:"seat_number"`
			Movie_id string `json:"movie_id"`
			User_id string `json:"user_id"`
			Price float64 `json:"price"`
		} `json:"new"`
	}
	
	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Define the GraphQL mutation string
	var mutation struct {
		UpdateUsers struct {
			Returning []struct {
				Title          string `json:"title"`

			} `json:"returning"`

		} `graphql:"update_movies(where: {tickets_aggregate: {count: {predicate: {_gte: $count}}}}, _set: {status: $status})"`
	}

	variables := map[string]interface{}{
		"status": "closed",
		"count": 20,
		
	}
	// execute the request
	err5 := utilService.Client().Mutate(context.Background(), &mutation, variables)
	if err5 != nil {
		ctx.JSON(400, gin.H{"error": err5.Error()})
		return
	}
	
	ctx.JSON(200, mutation.UpdateUsers.Returning)

}
