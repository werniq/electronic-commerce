package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"new-e-commerce/models"
)

func (app *application) CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Accept", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		/*
			Access-Control-Allow-Credentials:
				Indicates whether the response to the request can be exposed when the credentials flag is true.
				In your case, you are allowing credentials by setting this to "true".
		*/
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "300")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func (app *application) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "no authentication header received"})
			return
		}

		err := models.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		c.Next()
	}
}
