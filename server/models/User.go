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

// CheckPassword compares user actual password, with one, from input
func (u *User) CheckPassword(password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		fmt.Println(u.Password, password)
		return err
	}
	return
}
