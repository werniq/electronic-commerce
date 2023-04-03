package main

import "github.com/gin-gonic/gin"

func (app *application) SetupApiRoutes(router *gin.Engine) {
	router.Use(app.CorsMiddleware())
	router.POST("/api/signup", app.Authorize)

	router.Use(app.Auth())
	router.POST("/api/signin", app.GenerateToken)
	router.POST("/api/get-all", app.GetAll)
	router.POST("/api/create", app.Create)
	router.POST("/api/get/{id}", app.Details)
	router.POST("/api/edit/{id}", app.Edit)
	router.POST("/api/delete/{id}", app.Delete)
}
