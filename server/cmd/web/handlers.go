package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"new-e-commerce/models"
	"strconv"
	"time"
)

var (
	sessionTokenKey = "t:#f#YY_G(TA{Zp!&a^5YHNBK%f4C$c$M("
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

type SessionData struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	TokenExpiry string `json:"tokenExpiry"`
}

type PaginationData struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalPages  int `json:"totalPages"`
	TotalCount  int `json:"totalCount"`
}

// HomeHandler renders home page
func (app *application) HomeHandler(c *gin.Context) {
	if err := app.renderTemplate(c, "home", &templateData{}); err != nil {
		app.errorLog.Printf("error rendering home page: %v\n", err)
	}
}

// RegisterHandler renders page for registration
func (app *application) RegisterHandler(c *gin.Context) {
	if err := app.renderTemplate(c, "register", &templateData{}); err != nil {
		app.errorLog.Printf("error rendering home page: %v\n", err)
	}
}

// AuthenticationHandler renders page for authentication
func (app *application) AuthenticationHandler(c *gin.Context) {
	if err := app.renderTemplate(c, "authentication", &templateData{}); err != nil {
		app.errorLog.Printf("error rendering home page: %v\n", err)
	}
}

// Login function is rendering Login page
func (app *application) Login(c *gin.Context) {
	if err := app.renderTemplate(c, "login", &templateData{}); err != nil {
		app.errorLog.Printf("error rendering login page: %v\n", err)
	}
}

// CreateProduct function is rendering CreateProduct page
func (app *application) CreateProduct(c *gin.Context) {
	if err := app.renderTemplate(c, "createProduct", &templateData{}); err != nil {
		app.errorLog.Printf("error rendering createProduct page: %v\n", err)
	}
}

// Catalogue creates request to API, retrieves all books from database, and renders page with data from API response
func (app *application) Catalogue(c *gin.Context) {
	page := c.Param("page")
	var p int
	var err error
	if page == "" {
		p = 0
	}

	p, err = strconv.Atoi(page)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	uri := fmt.Sprintf(app.cfg.api+"/api/catalogue/page/%d", p)
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	email, err := c.Cookie("email")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	_, err = app.database.RetrieveTokenDataFromTable(email)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	//req.Header.Add("Authorization", "Bearer "+sessionData.Token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	var books struct {
		Books []models.Book `json:"books"`
	}

	err = json.NewDecoder(res.Body).Decode(&books)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	b, err := app.database.GetAllBooks()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	pages := int(len(b) / 4)

	var pg []int
	for i := 1; i <= pages; i++ {
		pg = append(pg, i)
	}

	data := make(map[string]interface{})
	data["books"] = books.Books
	data["pages"] = pg
	data["max"] = pages

	if err := app.renderTemplate(
		c,
		"catalogue",
		&templateData{
			Data: data,
		},
	); err != nil {
		app.errorLog.Println(err)
	}

}

// GetUserInfo stores data in database with name of user's email, so that I can receive token and tokenExpiry for my own purposes :D
func (app *application) GetUserInfo(c *gin.Context) {
	sessionData := &SessionData{}

	err := json.NewDecoder(c.Request.Body).Decode(&sessionData)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	c.SetCookie("email", sessionData.Email, 3600*24*7, "/", "localhost", false, true)
	sessionData.TokenExpiry = sessionData.TokenExpiry[:10]
	tokenExp, err := time.Parse("2006-01-02", sessionData.TokenExpiry)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	err = app.database.StoreEmailAsSessionToken(sessionData.Email, sessionData.Token, tokenExp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) ForgotPassword(c *gin.Context) {
	if err := app.renderTemplate(c, "forgot-password", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}

}

func (app *application) MyProfile(c *gin.Context)      {}
func (app *application) EditProduct(c *gin.Context)    {}
func (app *application) DeleteProduct(c *gin.Context)  {}
func (app *application) ChargeOnce(c *gin.Context)     {}
func (app *application) ResetPassword(c *gin.Context)  {}
func (app *application) ForgetPassword(c *gin.Context) {}
func (app *application) Contact(c *gin.Context)        {}
