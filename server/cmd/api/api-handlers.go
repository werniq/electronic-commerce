package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"new-e-commerce/models"
	"strings"
	"time"
)

type TokenData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Author      string    `json:"author"`
	Category    string    `json:"category"`
	AddCategory string    `json:"addCategory"`
	DateOfIssue time.Time `json:"date_of_issue"`
	QuoteFrom   string    `json:"quoteFrom"`
	Language    string    `json:"language"`
}

// Authorize registers user in this website
func (app *application) Authorize(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	u, err := app.database.FindUserByEmail(user.Email)
	if u != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id, err := app.GenerateUserId()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user.UserID = id
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		c.JSON(400, gin.H{"error": "error generating password " + err.Error()})
		return
	}

	user.Password = string(pass)
	user.CreatedAt = time.Now()
	err = app.database.InsertUser(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "successfully created new user"})
}

// GenerateToken generates JWT authentication token for user
func (app *application) GenerateToken(c *gin.Context) {
	var r TokenData
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": "error decoding request body: " + err.Error()})
		return
	}

	user, err := app.database.FindUserByEmail(r.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": "user not found " + err.Error()})
		return
	}

	err = user.CheckPassword(r.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	t, err := models.GenerateToken(user.Email, user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "error generating token for user: " + user.Email + " " + err.Error()})
		return
	}
	//var payload struct {
	//	Email       string    `json:"email"`
	//	Token       string    `json:"token"`
	//	TokenExpiry time.Time `json:"tokenExpiry"`
	//}
	//
	//payload.Email = t.Email
	//payload.Token = t.Token
	//payload.TokenExpiry = t.ExpiresAt
	//
	//body, err := json.Marshal(payload)
	//if err != nil {
	//	c.JSON(400, gin.H{"error": err.Error()})
	//	return
	//}
	//_, err = http.Post("http://localhost:4000/get-user-info", "application/json", bytes.NewBuffer(body))
	//if err != nil {
	//	c.JSON(400, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(200, gin.H{
		"error":       false,
		"token":       t.Token,
		"email":       t.Email,
		"tokenExpiry": t.ExpiresAt,
	})
}

// Create "/create" is used for creating new books
func (app *application) Create(c *gin.Context) {
	book := &Book{}
	err := c.ShouldBindJSON(&book)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	author := strings.Split(book.Author, " ")
	if len(author) <= 1 {
		c.JSON(400, gin.H{"error": "wrong author name"})
		return
	}
	aut, err := app.database.FindAuthorByName(author[0], author[1])
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	modelBook := &models.Book{
		Title:       book.Title,
		Description: book.Description,
		Price:       int(book.Price),
		AuthorID:    aut.AuthorID,
		DateOfIssue: book.DateOfIssue,
		QuoteFrom:   book.QuoteFrom,
		Language:    book.Language,
	}

	err = app.database.InsertBook(modelBook)
	if err != nil {
		c.JSON(400, gin.H{"error": "error inserting book: " + err.Error()})
		return
	}

	err = app.database.UpdateAuthorBook(aut, book.Title)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"msg": "book soon will appear in catalogue :)"})
}

func (app *application) Catalogue(c *gin.Context) {
	books, err := app.database.GetAllBooks()
	fmt.Println(books)
	fmt.Println(books)
	fmt.Println(books)
	if err != nil {
		c.JSON(400, gin.H{"error": "error getting all books from database: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"books": books,
	})
}

func (app *application) IsAuthenticated(c *gin.Context) {
	return
}

func (app *application) Edit(c *gin.Context)    {}
func (app *application) Details(c *gin.Context) {}
func (app *application) Remove(c *gin.Context)  {}
func (app *application) Delete(c *gin.Context)  {}
func (app *application) GetAll(c *gin.Context)  {}
