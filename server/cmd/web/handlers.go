package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"new-e-commerce/models"
)

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

func (app *application) Login(c *gin.Context) {
	if err := app.renderTemplate(c, "login", &templateData{}); err != nil {
		app.errorLog.Printf("error rendering login page: %v\n", err)
	}
}

func (app *application) CreateProduct(c *gin.Context) {
	if err := app.renderTemplate(c, "createProduct", &templateData{}); err != nil {
		app.errorLog.Printf("error rendering createProduct page: %v\n", err)
	}
}

func (app *application) Catalogue(c *gin.Context) {
	req, err := http.NewRequest("POST", app.cfg.api+"/api/catalogue", nil)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	token := app.Session.Get(context.Background(), "token")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", token.(string))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	var books []*models.Book
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
	var payload struct {
		Email       string `json:"email"`
		Token       string `json:"token"`
		TokenExpiry string `json:"tokenExpiry"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.Session.Put(context.Background(), "email", payload.Email)
	app.Session.Put(context.Background(), "token", payload.Token)
	app.Session.Put(context.Background(), "tokenExpiry", payload.TokenExpiry)
}

func (app *application) MyProfile(c *gin.Context)      {}
func (app *application) EditProduct(c *gin.Context)    {}
func (app *application) DeleteProduct(c *gin.Context)  {}
func (app *application) ChargeOnce(c *gin.Context)     {}
func (app *application) ResetPassword(c *gin.Context)  {}
func (app *application) ForgetPassword(c *gin.Context) {}
func (app *application) Contact(c *gin.Context)        {}
