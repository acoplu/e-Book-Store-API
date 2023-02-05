package main

import (
	"CelikGroupCRUDAPI/models" //user and book struct
	"github.com/gin-gonic/gin" //golang http web framework
	"strconv"
)

// define our User object to keep user's feedbacks and purchases
var user = models.User{Feedbacks: nil, Purchases: nil}

func main() {
	//creating a router to create routes that associate HTTP requests with URL paths and handler functions
	router := gin.Default()
	models.ConnectDatabase()

	//first publishers of books have to update book info frequently so let's define a route for this
	router.POST("/books", addBook)
	//let's define a demo route that show us our database, so we can check if our routes work well
	router.GET("/books", getBooks)
	//now let's define a route for customers. because customers will be able to view the book and see the other books with the same author
	router.GET("/books/:id", checkBook)
	//now let's define a route for providing feedback. we can use it with PATCH HTTP request
	router.PATCH("/feedback/:id", giveFeedback)
	//now let's define a shopping card for user's purchases
	router.GET("/cart", getCard)
	//now let's define a route for customers to make be able to purchase
	router.PATCH("/purchase", buyBook)
	//to run a local server in :8080 port
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}

// this handler function will show all books
func getBooks(c *gin.Context) {
	//first create our Book array object
	var books []models.Book

	//now this Find function will insert all Book instances to Book array. we use pointer to array otherwise the value
	//of the array won't change. because golang is pass-by-value
	models.DB.Find(&books)

	//listing all books database sorted by IDs. because we keep IDs in string data type we have to parse it
	for i := 0; i < len(books); i++ {
		for j := 0; j < len(books)-1; j++ {
			num1, err1 := strconv.Atoi(books[j].ID)
			num2, err2 := strconv.Atoi(books[j+1].ID)
			if err1 == nil && err2 == nil {
				if num1 > num2 {
					temp := books[j]
					books[j] = books[j+1]
					books[j+1] = temp
				}
			}
		}
	}

	c.IndentedJSON(200, books)
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

func checkBook(c *gin.Context) {
	//in this function customers will provide a book name, and they can see the details of book like reviews and else.
	//moreover they can see the other books with the same author
	id := c.Param("id")

	//now let's define an empty book object and call Where function with *orm.DB instance, so we can take the ID of the
	//book
	var newBook models.Book

	//now this chaining method calling will take our Param in URL path as ID and then check if there is a book with this
	//ID. then it will bind the matched book to our newBook variable. if there will be error than err cannot be nil, and
	//we have to give 404
	if err := models.DB.Where("id = ?", id).First(&newBook).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": "No matched book"})
		return
	}

	//let's define a slice to keep the same name authors, and then we can print them easily
	var books, sameAuthor []models.Book
	models.DB.Find(&books)
	for i := 0; i < len(books); i++ {
		if newBook.Author == books[i].Author && newBook.ID != books[i].ID {
			sameAuthor = append(sameAuthor, books[i])
		}
	}
	c.IndentedJSON(200, newBook)
	//if the sameAuthor slice is nil then it will print null on screen, so let's handle it
	if sameAuthor != nil {
		c.IndentedJSON(200, sameAuthor)
	}
}

// now in this function we will take the ID of the book, the rating and the short comment about the book
func giveFeedback(c *gin.Context) {
	//for taking the ID of the book that customer will give feedback lets use GetQuery method
	id := c.Param("id")

	//let's get the request JSON and bind it to exBook. We will take the request through this Book
	var exBook, originalBook models.Book
	if err := c.BindJSON(&exBook); err != nil {
		c.IndentedJSON(404, gin.H{"message": "Invalid data"})
		return
	}

	//now this chaining method calling will take our Param in URL path as ID and then check if there is a book with this
	//ID. then it will bind the matched book to our newBook variable. if there will be error than err cannot be nil, and
	//we have to give 404
	if err := models.DB.Where("id = ?", id).First(&originalBook).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": "No matched book"})
		return
	}

	//now check if user has such feedback on this book
	for i := 0; i < len(user.Feedbacks); i++ {
		if user.Feedbacks[i] == originalBook.ID {
			c.IndentedJSON(404, gin.H{"message": "Already given a feedback"})
			return
		}
	}

	//let's update in database
	models.DB.Model(&originalBook).Update("Rating", (originalBook.Rating+exBook.Rating)/float64(len(originalBook.Comments)+1))
	models.DB.Model(&originalBook).Update("Comments", append(originalBook.Comments, exBook.Comments[0]))
	user.Feedbacks = append(user.Feedbacks, originalBook.ID)
	c.IndentedJSON(200, originalBook)
}

// the bookstore can see the ordered books list and details in shopping card section. so let's define a shopping card handler func
func getCard(c *gin.Context) {
	if user.Purchases == nil {
		c.IndentedJSON(200, gin.H{"message": "No purchases yet"})
		return
	}

	//now we have to print the details of the list and the total cost
	var totalCost float64
	for i := 0; i < len(user.Purchases); i++ {
		totalCost = totalCost + (user.Purchases[i].Price)*(float64(user.Purchases[i].Quantity))
	}
	c.IndentedJSON(200, user.Purchases)
	c.IndentedJSON(200, totalCost)
}

// customers who wish to purchase can select their books into the website's cart
func buyBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	var matchedBook models.Book

	if !ok {
		c.IndentedJSON(404, gin.H{"message": "Invalid ID"})
		return
	}

	//now this chaining method calling will take our Param in URL path as ID and then check if there is a book with this
	//ID. then it will bind the matched book to our newBook variable. if there will be error than err cannot be nil, and
	//we have to give 404
	if err := models.DB.Where("id = ?", id).First(&matchedBook).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": "No matched book"})
		return
	}

	//and now check if there is a 0 book
	if matchedBook.Quantity <= 0 {
		c.IndentedJSON(404, gin.H{"message": "Book is not available"})
		return
	}

	bookForBuy := matchedBook
	//if user buys a book decrement the quantity 1 in database
	models.DB.Model(&matchedBook).Update("Quantity", matchedBook.Quantity-1)
	for i := 0; i < len(user.Purchases); i++ {
		if user.Purchases[i].Name == bookForBuy.Name {
			user.Purchases[i].Quantity++
			c.IndentedJSON(200, gin.H{"message": "Added to shopping cart"})
			return
		}
	}

	bookForBuy.Quantity = 1
	user.Purchases = append(user.Purchases, bookForBuy)
	c.IndentedJSON(200, gin.H{"message": "Added to shopping cart"})
}
