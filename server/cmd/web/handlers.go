package main

import "github.com/gin-gonic/gin"

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

func (app *application) MyProfile(c *gin.Context) {}
func (app *application) Catalogue(c *gin.Context) {}

func (app *application) EditProduct(c *gin.Context)    {}
func (app *application) DeleteProduct(c *gin.Context)  {}
func (app *application) ChargeOnce(c *gin.Context)     {}
func (app *application) ResetPassword(c *gin.Context)  {}
func (app *application) ForgetPassword(c *gin.Context) {}
func (app *application) Contact(c *gin.Context)        {}
