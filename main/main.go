package main

import (
	"CelikGroupCRUDAPI/models" //user and book struct
	//"github.com/gin-gonic/gin" //golang http web framework
)

// define our User object to keep user's feedbacks and purchases
var user = models.User{Feedbacks: nil, Purchases: nil}
