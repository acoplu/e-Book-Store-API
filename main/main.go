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

	//to run a local server in :8080 port
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
