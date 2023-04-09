package main

import "github.com/gin-gonic/gin"

func (app *application) SetupApiRoutes(router *gin.Engine) {
	router.Use(app.CorsMiddleware())
	router.POST("/api/signup", app.Authorize)
	router.POST("/api/signin", app.GenerateToken)

	//router.Use(app.Auth())
	router.POST("/api/is-authenticated", app.IsAuthenticated)
	router.POST("/api/catalogue", app.Catalogue)
	router.POST("/api/create", app.Create)
	router.POST("/api/get/{id}", app.Details)
	router.POST("/api/edit/{id}", app.Edit)
	router.POST("/api/delete/{id}", app.Delete)
}
