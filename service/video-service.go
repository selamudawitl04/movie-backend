package service

import "gilab.com/progrmaticreviwes/golang-gin-poc/entity"
type VideoService interface{
	Save(entity.Video) entity.Video
	FindAll() []entity.Video
}

type videoService struct{
	videos []entity.Video
}


// the videoService object
func New() VideoService{
	return &videoService{}
}
// save video service 

func(service *videoService) Save(video entity.Video) entity.Video{
	service.videos = append(service.videos, video)
	return video
}

// get video service
func(service *videoService) FindAll() []entity.Video{
	return service.videos
}



// func handleSignup(ctx *gin.Context) {
//     // Get the user input from the request body
// 	// Create graphql client
// 	// Set up the HTTP client with the request headers
// 	// An HTTP transport that adds headers to requests

//     type inputUser struct {
// 		Input struct{
// 			Arg1 struct{
// 				FirstName string `json:"firstName"`
// 				LastName string `json:"lastName"`
// 				Email string `json:"email"`
// 				Password string `json:"password"`
// 			} `json:"arg1"`
// 		} `json:"input"`
//     }

// 	var newUser inputUser
//     if err := ctx.ShouldBindJSON(&newUser); err != nil {
//         ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

// 	fmt.Println(newUser.Input.Arg1.FirstName)
// 	fmt.Println(newUser.Input.Arg1.LastName)
// 	fmt.Println(newUser.Input.Arg1.Email)
// 	fmt.Println(newUser.Input.Arg1.Password)

	
//     // Set up the GraphQL client
// 	// Set up the HTTP client with the request headers
// 	// An HTTP transport that adds headers to requests
// 	headers := http.Header{}
// 	headers.Add("X-Hasura-Admin-Secret", "ym6arlrrdMol6MfV156smTMo8L72B6QBLxiyZtWUZl0w0YxctdVN9YTppWkYB5Gn")
// 	httpClient := &http.Client{Transport: &headersTransport{headers, http.DefaultTransport}}

// // Set up the GraphQL client
// 	// client := graphql.NewClient("https://vue-shopping.hasura.app/v1/graphql", httpClient)


	
//     // Set up the request headers
    
// 	// Define the GraphQL query you want to execute
// 	// query := `
// 	// 	query {
// 	// 		user {
// 	// 			id
// 	// 			firstName
// 	// 			lastName
// 	// 			email
// 	// 			password
// 	// 		}
// 	// 	}
// 	// `

// 	// type User struct {
// 	// 	ID        int    `graphql:"id"`
// 	// 	FirstName string `graphql:"firstName"`
// 	// }
	
// 	// var query struct {
// 	// 	Users []User `graphql:"user"`
// 	// }
// 	// var query struct {
// 	// 	User struct {
// 	// 		ID          int
// 	// 		FirstName string
// 	// 	} `graphql:"user_by_pk(id:2)"`
// 	// }
	

// 	// err := client.Query(context.Background(), &query, nil)
// 	// if err != nil {
// 	// 	// Handle error.
// 	// 	fmt.Println((err.Error()))
// 	// }
// 	// fmt.Println(query.Users[0].FirstName)

	


// 	// Define a struct to unmarshal the response into
	 

	

// 	// Execute the GraphQL query and unmarshal the response into the Response struct
// 	// response := &Response{}
// 	// if err := client.Query(context.Background(), query, nil, response); err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Printf("%+v\n", query)
// 	// fmt.Printf("%s", "Selamu Dawit Se")

// 	// return
// 	// client.Mutate()


//     // Define the GraphQL mutation string
//     // mutation := `
//     //     mutation ($firstName: String!, $lastName: String!, $email: String!, $password: String!) {
//     //         insert_users(objects: {firstName: $firstName, lastName: $lastName, email: $email, password: $password}) {
//     //             returning {
// 	// 				id
// 	// 				firstName
// 	// 			  }
//     //         }
//     //     }
//     // `

// 	// var m struct {
// 	// 	CreateUser struct {
// 	// 		ID      int
// 	// 		Email string
// 	// 	} `graphql:"createUser(episode: $ep, review: $review)"`
// 	// }
	
	
	

// 	var mut struct {
// 		CreateUser struct {
// 			ID      int
// 			Email string
			
// 		} `graphql:"createUser(firstName: $firstName, lastName: $lastName, email: $email, password: $password)"`
// 	}

//     // Set up the variables for the GraphQL mutation
//     variables := map[string]interface{}{
//         "firstName":     newUser.Input.Arg1.FirstName,
//         "lastName":     newUser.Input.Arg1.FirstName,
//         "email":    newUser.Input.Arg1.Email,
//         "password": newUser.Input.Arg1.Password,
//     }

//     // Send the GraphQL mutation to Hasura Cloud
//     // var response struct {
//     //     InsertUsers struct {
//     //         AffectedRows int `json:"affected_rows"`
//     //     } `json:"insert_users"`
//     // }
//     if err := client.Mutate(context.Background(), mut, variables, ); err != nil {
// 		fmt.Println((err.Error()))
//         ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

// 	fmt.Println(mut.CreateUser.Email)
// 	fmt.Println(mut.CreateUser.ID)

	
	

//     // Check if the user was successfully created
//     // if response.InsertUsers.AffectedRows == 0 {
//     //     ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
//     //     return
//     // }

//     // Return a success response
// 	// var newResponse response
// 	// newResponse.accessToken = "asdhfjvamsdghk"
// 	ctx.JSON(200, gin.H{"accessToken": "asvhkfyeekruow3asfdsvBZCM"})


// }