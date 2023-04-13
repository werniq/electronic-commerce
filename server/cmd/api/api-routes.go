package main

import "github.com/gin-gonic/gin"

func (app *application) SetupApiRoutes(router *gin.Engine) {
	router.Use(app.CorsMiddleware())
	router.POST("/api/signup", app.Authorize)
	router.POST("/api/signin", app.GenerateToken)

	//router.Use(app.Auth())
	router.POST("/api/is-authenticated", app.IsAuthenticated)
	router.POST("/api/catalogue/page/:page", app.Catalogue)
	router.POST("/api/forgot-password", app.SendPasswordResetEmail)
	router.POST("/api/create", app.Create)

	router.POST("/api/books/buy/:id", app.GetPaymentIntent)
	router.POST("/api/books/delete/:id", app.Remove)
	router.POST("/api/books/edit/:id", app.Edit)

	//{{.API}}/api/book/order/{{$book.ID}}
	router.POST("/api/book/order/:id", app.GetPaymentIntent)

}
