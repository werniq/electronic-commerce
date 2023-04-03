package main

import (
	"github.com/gin-gonic/gin"
	"new-e-commerce/models"
)

type TokenData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	var u models.User
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": "error decoding request body: " + err.Error()})
		c.Abort()
		return
	}

	user, err := app.database.FindUserByEmail(u.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": "user not found " + err.Error()})
		c.Abort()
		return
	}

	err = user.CheckPassword(r.Password)
	if err != nil {
		c.JSON(401, gin.H{"invalid credentials ": err.Error()})
		c.Abort()
		return
	}

	t, err := models.GenerateToken(user.Email, user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "error generating token for user: " + user.Email + " " + err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"error":       false,
		"token":       t.Token,
		"tokenExpiry": t.ExpiresAt,
	})
}

func (app *application) Create(c *gin.Context)  {}
func (app *application) Edit(c *gin.Context)    {}
func (app *application) Details(c *gin.Context) {}
func (app *application) Remove(c *gin.Context)  {}
func (app *application) Delete(c *gin.Context)  {}
func (app *application) GetAll(c *gin.Context)  {}
