package main

import "github.com/gin-gonic/gin"

func (app *application) SetupRoutes(r *gin.Engine) {
	r.GET("/", app.HomeHandler)
	r.GET("/register", app.RegisterHandler)
	r.GET("/login", app.Login)
	r.GET("/my-profile", app.AuthenticationHandler)
	r.GET("/catalogue/:page", app.Catalogue)
	r.GET("/create", app.CreateProduct)
	r.GET("/edit/{id}", app.AuthenticationHandler)
	r.GET("/delete/{id}", app.AuthenticationHandler)
	r.GET("/details/{id}", app.AuthenticationHandler)
	r.GET("/charge-once", app.AuthenticationHandler)
	r.GET("/reset-password", app.AuthenticationHandler)
	r.GET("/forgot-password", app.AuthenticationHandler)
	r.GET("/contact", app.AuthenticationHandler)
	r.POST("/get-user-info", app.GetUserInfo)
}
