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

/*
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

*/

// UpdateExistingBook is used to update information about particular book
func (b *Book) UpdateExistingBook(newBook Book) {
	if newBook.Title != "" {
		b.Title = newBook.Title
	}

	if newBook.Description != "" {
		b.Description = newBook.Description
	}

	if newBook.Price != 0 {
		b.Price = newBook.Price
	}

	if newBook.QuoteFrom != "" {
		b.QuoteFrom = newBook.QuoteFrom
	}

	if newBook.Language != "" {
		b.Language = newBook.Language
	}

	if newBook.Category != "" {
		b.Category = newBook.Category
	}

	if newBook.AddCategory != "" {
		b.AddCategory = newBook.AddCategory
	}

	if newBook.AuthorID != 0 {
		b.AuthorID = newBook.AuthorID
	}
}
