package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"new-e-commerce/models"
	"new-e-commerce/utils/cards"
	"new-e-commerce/utils/urlsigner"
	"strconv"
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

type StripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	CardBrand     string `json:"card_brand"`
	ExpiryMonth   int    `json:"exp_month"`
	ExpiryYear    int    `json:"exp_year"`
	LastFour      string `json:"last_four"`
	Plan          string `json:"plan"`
	ProductID     string `json:"product_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

type JsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type ResetPasswordData struct {
	Email string `json:"email"`
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
		Category:    book.Category,
		AddCategory: book.AddCategory,
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
	p := c.Param("page")
	page, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		app.errorLog.Println(err)
		return
	}

	books, err := app.database.GetPaginatedBooks(page)
	if err != nil {
		c.JSON(400, gin.H{"error": "error getting all books from database: " + err.Error()})
		app.errorLog.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"books": books,
	})
}

// GetPaginatedCatalogue is used for showing exactly 4 books at client catalogue
func (app *application) GetPaginatedCatalogue(c *gin.Context) {
	page := c.Request.URL.Query().Get("page")
	p, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		app.errorLog.Println(err)
		return
	}

	books, err := app.database.GetPaginatedBooks(p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		app.errorLog.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"books": books,
	})
}

func (app *application) IsAuthenticated(c *gin.Context) {
	return
}

func (app *application) SendPasswordResetEmail(c *gin.Context) {
	var data ResetPasswordData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		app.errorLog.Println(err)
		return
	}

	_, err := app.database.FindUserByEmail(data.Email)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{"error": "user with given email not found"})
		app.errorLog.Println(err)
		return
	}

	link := fmt.Sprintf("%s/reset-password?email=%s", app.cfg.client, data.Email)

	sign := urlsigner.Signer{Secret: []byte(app.cfg.secretKey)}

	signedLink := sign.GenerateTokenFromString(link)

	var payload struct {
		Link string
	}

	payload.Link = signedLink

	// send email
	err = app.SendEmail("info@qniwwwersss.com", "info@qniwwwersss.com", "Password Reset Request", "password-reset", payload)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		app.errorLog.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Check your inbox!"})
}

// func (app *application) OrderBook(c *gin.Context) {
//	id := c.Param("id")
//	bookID, err := strconv.Atoi(id)
//	if err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		app.errorLog.Println(err)
//		return
//	}
//
//	book, err := app.database.GetBookById(bookID)
//	if err != nil {
//		c.JSON(500, gin.H{"error": fmt.Sprintf("error searching book : %v", err)})
//		app.errorLog.Println(err)
//		return
//	}
//
//}

func (app *application) GetPaymentIntent(c *gin.Context) {
	var payload StripePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{"error": err.Error()})
		app.errorLog.Println(err)
		return
	}
	amount, err := strconv.Atoi(payload.Amount)
	amount = amount * 100
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error " + err.Error()})
		app.errorLog.Println(err)
		return
	}

	card := cards.Card{
		Secret:   app.cfg.stripe.secret,
		Key:      app.cfg.stripe.key,
		Currency: payload.Currency,
	}

	ok := true
	paymentIntent, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		ok = false
	}

	if ok {
		c.JSON(200, gin.H{"data": paymentIntent})
		return
	} else {
		j := JsonResponse{
			OK:      false,
			Message: msg,
			Content: "",
		}

		c.JSON(400, gin.H{"error": j.Message})
	}
}

func (app *application) Edit(c *gin.Context)    {}
func (app *application) Details(c *gin.Context) {}
func (app *application) Remove(c *gin.Context)  {}
func (app *application) Delete(c *gin.Context)  {}
func (app *application) GetAll(c *gin.Context)  {}
