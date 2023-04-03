package models

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

type DatabaseModel struct {
	DB *sql.DB
}

var errorLogger = log.New(os.Stdout, "ERROR\t", log.Lshortfile|log.Ltime|log.Lshortfile)

// FindUserByEmail retrieves user from database by given email.
func (m *DatabaseModel) FindUserByEmail(email string) (*User, error) {
	stmt := "SELECT * FROM users where email=$1"

	rows, err := m.DB.Query(stmt, email)
	if err != nil {
		errorLogger.Println(err)
		return nil, err
	}
	defer rows.Close()

	u := &User{}
	found := false

	for rows.Next() {
		found = true
		if err := rows.Scan(
			&u.ID,
			&u.UserID,
			&u.Username,
			&u.Email,
			&u.Phone,
			&u.Password,
			&u.CreatedAt,
		); err != nil {
			errorLogger.Println(err)
			continue
		}
	}

	if !found {
		return nil, errors.New("user not found")
	}

	return u, nil
}

// InsertUser creates new user record
func (m *DatabaseModel) InsertUser(u *User) error {
	stmt := `
		INSERT INTO
				users(userID, username, email, phone, password, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
			`

	_, err := m.DB.Exec(stmt, u.UserID, u.Username, u.Email, u.Phone, u.Password, u.CreatedAt)
	if err != nil {
		errorLogger.Println(err)
		return err
	}
	return nil
}
