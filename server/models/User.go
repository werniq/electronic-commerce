package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	UserID    string    `json:"userID"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" binding="required"`
	PasswordConfirm string `json:"password" binding="required"`
}

// CheckPassword compares user actual password, with one, from input
func (u *User) CheckPassword(password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		fmt.Println(u.Password, password)
		return err
	}
	return
}

// ChangeUsername function used for changing user's username
func (u *User) ChangeUsername(us string) {
	u.Username = us
}

// ChangePassword function used for changing user's password
func (u *User) ChangePassword(newPassword string) (err error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return
	}
	u.Password = string(pass)
	return nil
}
