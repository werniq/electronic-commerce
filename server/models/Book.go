package models

import "time"

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	DateOfIssue time.Time `json:"date_of_issue"`
	QuoteFrom   string    `json:"quote_from"`
	Language    string    `json:"language"`
	Review      []Review  `json:"review"`
}

type Review struct {
	Rating  float32 `json:"rating"`
	Comment string  `json:"comment"`
	User    *User   `json:"user"`
	Likes   int     `json:"likes"`
}
