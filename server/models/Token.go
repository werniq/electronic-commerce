package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var (
	JwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

type Token struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// GenerateToken returns signed JWT token
func GenerateToken(email, username string) (*Token, error) {
	expTime := time.Now().Add(time.Hour * 24 * 7)
	claims := &Claims{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecretKey)
	if err != nil {
		return nil, err
	}

	t := &Token{
		Username:  username,
		Email:     email,
		Token:     tokenString,
		ExpiresAt: expTime,
	}

	return t, nil
}

func ValidateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecretKey), nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return errors.New("couldn't parse jwt claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return errors.New("your authentication token has expired")
	}
	return nil
}

func (t *Token) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}
