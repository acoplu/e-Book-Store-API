package main

import (
	"CelikGroupCRUDAPI/models" //user and book struct
	"github.com/gin-gonic/gin" //golang http web framework
)

// define our User object to keep user's feedbacks and purchases
var user = models.User{Feedbacks: nil, Purchases: nil}

func main() {
	//creating a router to create routes that associate HTTP requests with URL paths and handler functions
	router := gin.Default()
	models.ConnectDatabase()

	//first publishers of books have to update book info frequently so let's define a route for this
	router.POST("/books", addBook)

	//to run a local server in :8080 port
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}

// our all handler functions must take gin.Context pointer as parameter
func addBook(c *gin.Context) {
	//first we have to check the POST JSON info with BindJSON function. it will return err if the JSON POST request invalid
	var newBook models.Book //we will convert JSON request to data of this object

	// we are passing our reference of newBook otherwise function cannot change object's data. because golang is
	//pass-by-value. this function will return error if the binding process cannot be handled
	if err := c.BindJSON(&newBook); err != nil {
		//if the error not null it means JSON request is invalid. so let's make HTTP status 404 and give a JSON
		//message to screen with gin.H{}
		c.IndentedJSON(404, gin.H{"message": "Invalid data"})
		return //finish the execution of func
	}

	//after we bind it to our newBook instance let's create a new ROW in our database. this "Create" function will
	//return an error if it cannot create this new book. for example if there is already a book with same primary key
	//than it will return an error
	if err := models.DB.Create(&newBook).Error; err != nil {
		//let's give an error message
		c.IndentedJSON(404, gin.H{"message": "Cannot add it"})
		return
	}

	//if the error is null it means JSON request is valid and the method bound this JSON to object, and there
	//is no other book with the same ID, so we can add our book to our database
	c.IndentedJSON(200, newBook)
}
