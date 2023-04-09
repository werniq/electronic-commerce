package models

import (
	"time"
)

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	DateOfIssue time.Time `json:"date_of_issue"`
	QuoteFrom   string    `json:"quote_from"`
	Language    string    `json:"language"`
	Category    string    `json:"category"`
	AddCategory string    `json:"addcategory"`
	AuthorID    int       `json:"author_id"`
}

type Review struct {
	Rating  float32 `json:"rating"`
	Comment string  `json:"comment"`
	User    *User   `json:"user"`
	Likes   int     `json:"likes"`
}

type Author struct {
	ID          int       `json:"id"`
	AuthorID    int       `json:"author_id"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Biography   string    `json:"biography"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Books       string    `json:"book"`
}
