package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"new-e-commerce/models"
	"time"
)

var (
	sessionTokenKey = "t:#f#YY_G(TA{Zp!&a^5YHNBK%f4C$c$M("
)

type SessionData struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	TokenExpiry string `json:"tokenExpiry"`
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
	req, err := http.NewRequest("POST", app.cfg.api+"/api/catalogue", nil)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	email, err := c.Cookie("email")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	token, _, err := app.database.RetrieveTokenDataFromTable(email)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	var books []models.Book
	err = json.NewDecoder(res.Body).Decode(&books)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(map[string]interface{})
	data["books"] = books
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

func (app *application) MyProfile(c *gin.Context)      {}
func (app *application) EditProduct(c *gin.Context)    {}
func (app *application) DeleteProduct(c *gin.Context)  {}
func (app *application) ChargeOnce(c *gin.Context)     {}
func (app *application) ResetPassword(c *gin.Context)  {}
func (app *application) ForgetPassword(c *gin.Context) {}
func (app *application) Contact(c *gin.Context)        {}
