package main

import (
	"crypto/rand"
)

func (app *application) GenerateUserId() (string, error) {
	p, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		return "", err
	}
	return p.String(), nil
}
