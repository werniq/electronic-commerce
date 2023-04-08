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
			return nil, err
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

// InsertBook basically inserts book into database
func (m *DatabaseModel) InsertBook(b *Book) error {
	stmt := `
			INSERT INTO
					book(title, description, price, date_of_issue, quote_from, language)
			VALUES
			    	($1, $2, $3, $4, $5, $6)
			`

	_, err := m.DB.Exec(stmt, b.Title, b.Description, b.Price, b.DateOfIssue, b.QuoteFrom, b.Language)
	if err != nil {
		return err
	}

	author, err := m.FoundAuthor(b.AuthorID)
	if err != nil {
		return err
	}
	err = m.UpdateAuthorBook(author, b.Title)

	if err != nil {
		return err
	}

	return nil
}

// FoundAuthor seeks for author by firstname and lastname
func (m *DatabaseModel) FoundAuthor(authorID int) (*Author, error) {
	stmt := `
			SELECT 
			    * 
			FROM author
				WHERE 
				    author_id = $1 
			`
	rows, err := m.DB.Query(stmt, authorID)
	if err != nil {
		return nil, err
	}
	a := &Author{}

	for rows.Next() {
		rows.Scan(
			&a.ID,
			&a.Firstname,
			&a.Lastname,
			&a.Biography,
			&a.DateOfBirth,
			&a.Books,
		)
	}

	return a, nil
}

// UpdateAuthorBook retrieves author from database, and adds book from arguments into author(book) field
func (m *DatabaseModel) UpdateAuthorBook(aut *Author, book string) error {
	stmt := `
			 SELECT 
			     * 
			 FROM 
			     author 
			 where 
			     id = $1
			`
	rows, err := m.DB.Query(stmt, aut.ID)
	if err != nil {
		return err
	}

	a := &Author{}
	if rows.Next() {
		err = rows.Scan(
			&a.ID,
			&a.Firstname,
			&a.Lastname,
			&a.Biography,
			&a.DateOfBirth,
			&a.Books,
		)
	}

	books := a.Books
	books += book
	a.Books = books

	stmt = `
				UPDATE 
				    author
				SET 
				    books = $1
				WHERE 
				    id = $2
			`

	_, err = m.DB.Exec(stmt, a.Books, a.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllBooks retrieves all existing in database books
func (m *DatabaseModel) GetAllBooks() ([]*Book, error) {
	stmt := `
			SELECT * FROM book
		`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	var books []*Book

	for rows.Next() {
		book := &Book{}
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.Price,
			&book.AuthorID,
			&book.QuoteFrom,
			&book.DateOfIssue,
			&book.Language,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (m *DatabaseModel) FindAuthorByName(firstname, lastname string) (*Author, error) {
	stmt := `
			select 
			    * 
			from author 
				where 
				    firstname = $1 
				  and 
				    lastname = $2`

	row, err := m.DB.Query(stmt, firstname, lastname)
	if err != nil {
		return nil, err
	}
	/*
			ID          int       `json:"id"`
		AuthorID    int       `json:"author_id"`
		Firstname   string    `json:"firstname"`
		Lastname    string    `json:"lastname"`
		Biography   string    `json:"biography"`
		DateOfBirth time.Time `json:"date_of_birth"`
		Books       string    `json:"book"`
	*/
	a := &Author{}
	if row.Next() {
		err = row.Scan(
			&a.ID,
			&a.AuthorID,
			&a.Firstname,
			&a.Lastname,
			&a.Biography,
			&a.DateOfBirth,
			&a.Books,
		)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}
