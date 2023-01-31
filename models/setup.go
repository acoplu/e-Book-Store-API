package models

import (
	"github.com/lib/pq"
	//"gorm.io/driver/postgres"
	//"gorm.io/gorm"
)

//now firstly we have to create a "Book" struct for our library. library will contain this "Book" objects. we have to
//define our fields with their json format.

type Book struct {
	//to record number of book lets define ID, but I removed it
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Author   string         `json:"author"`
	Price    float64        `json:"price"`
	Rating   float64        `json:"rating"`
	Comments pq.StringArray `gorm:"type:text[]" json:"comments"` //if we use string slice rather than "pq.StringArray" than it will be give error like cannot understand []string
	Quantity int            `json:"quantity"`
}

//now let's define a "User" struct to keep the user's feedback because each user can give only 1 feedback for each book.
//and we can keep the user's purchases too.

type User struct {
	//our feedbacks slice will hold the ID's of the books that user gave feedback
	Feedbacks []string `json:"feedbacks"`
	Purchases []Book   `json:"purchases"`
}
