package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type DatabaseModel struct {
	DB *sql.DB
}

type DbBook struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Price       int            `json:"price"`
	DateOfIssue time.Time      `json:"date_of_issue"`
	QuoteFrom   sql.NullString `json:"quote_from"`
	Language    sql.NullString `json:"language"`
	Category    sql.NullString `json:"category"`
	AddCategory sql.NullString `json:"addcategory"`
	AuthorID    int            `json:"author_id"`
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

func (m *DatabaseModel) GetAllBooks() ([]Book, error) {
	stmt := `
			SELECT
			    id, title, description,
			    price, date_of_issue, quote_from,
			    language, category, addcategory, author_id
			FROM book;
			`

	var books []Book
	rows, err := m.DB.Query(stmt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		book := DbBook{}
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.Price,
			&book.DateOfIssue,
			&book.QuoteFrom,
			&book.Language,
			&book.Category,
			&book.AddCategory,
			&book.AuthorID,
		)

		modelBook := Book{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Price:       book.Price,
			DateOfIssue: book.DateOfIssue,
			QuoteFrom:   "",
			Language:    "",
			Category:    "",
			AddCategory: "",
			AuthorID:    book.AuthorID,
		}

		if book.QuoteFrom.Valid {
			modelBook.QuoteFrom = book.QuoteFrom.String
		}

		if book.Language.Valid {
			modelBook.Language = book.Language.String
		}

		if book.Category.Valid {
			modelBook.Category = book.Category.String
		}

		if book.AddCategory.Valid {
			modelBook.AddCategory = book.AddCategory.String
		}

		if err != nil {
			return nil, err
		}
		books = append(books, modelBook)
	}
	return books, nil
}

// FindAuthorByName searches for author by given first and lastname. In case err == sql.ErrNoRows (author does not exist
// in database), creates new author record
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
	if err != sql.ErrNoRows {
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

// CreateNewAuthor func creates new author record in database
func (m *DatabaseModel) CreateNewAuthor(a *Author) error {
	stmt := `
		INSERT INTO 
		    author
		    (firstname, lastname, biography	, date_of_birth, books, author_id) 
		VALUES
		    ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.Exec(stmt, a.Firstname, a.Lastname, a.Biography, a.DateOfBirth, a.Books, a.AuthorID)
	if err != nil {
		return err
	}

	return nil
}

// GetPaginatedBooks is used for returning books in paginated form
func (m *DatabaseModel) GetPaginatedBooks(page int) ([]Book, error) {
	offset := page * 4
	//var book DbBook

	stmt := `
			SELECT
			    id, title, description,
			    price, date_of_issue, quote_from,
			    language, category, addcategory, author_id
			FROM book
			    LIMIT $1 OFFSET $2;
			`

	var books []Book
	rows, err := m.DB.Query(stmt, 4, offset)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		book := DbBook{}
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.Price,
			&book.DateOfIssue,
			&book.QuoteFrom,
			&book.Language,
			&book.Category,
			&book.AddCategory,
			&book.AuthorID,
		)

		modelBook := Book{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Price:       book.Price,
			DateOfIssue: book.DateOfIssue,
			QuoteFrom:   "",
			Language:    "",
			Category:    "",
			AddCategory: "",
			AuthorID:    book.AuthorID,
		}

		if book.QuoteFrom.Valid {
			modelBook.QuoteFrom = book.QuoteFrom.String
		}

		if book.Language.Valid {
			modelBook.Language = book.Language.String
		}

		if book.Category.Valid {
			modelBook.Category = book.Category.String
		}

		if book.AddCategory.Valid {
			modelBook.AddCategory = book.AddCategory.String
		}

		if err != nil {
			return nil, err
		}
		books = append(books, modelBook)
	}
	return books, nil
}

func (m *DatabaseModel) GetBookById(id int) (Book, error) {
	var b Book
	stmt := `
		SELECT 
		    id, title, description,
			price, date_of_issue, quote_from,
			language, category, addcategory, author_id
		FROM book 
			WHERE 
			    id = $1
		`
	res, err := m.DB.Query(stmt, id)
	if err != nil {
		return b, err
	}

	book := &Book{}

	for res.Next() {
		book := DbBook{}
		err = res.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.Price,
			&book.DateOfIssue,
			&book.QuoteFrom,
			&book.Language,
			&book.Category,
			&book.AddCategory,
			&book.AuthorID,
		)

		modelBook := Book{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Price:       book.Price,
			DateOfIssue: book.DateOfIssue,
			QuoteFrom:   "",
			Language:    "",
			Category:    "",
			AddCategory: "",
			AuthorID:    book.AuthorID,
		}

		if book.QuoteFrom.Valid {
			modelBook.QuoteFrom = book.QuoteFrom.String
		}

		if book.Language.Valid {
			modelBook.Language = book.Language.String
		}

		if book.Category.Valid {
			modelBook.Category = book.Category.String
		}

		if book.AddCategory.Valid {
			modelBook.AddCategory = book.AddCategory.String
		}

		if err == nil {
			return modelBook, err
		}
	}
	fmt.Println(book)
	return b, nil
}

func (m *DatabaseModel) RemoveBookFromDatabase(bookID int) error {
	stmt := `
			DELETE 
				* 
			FROM 
				book 
			WHERE 
			    id = $1
		`

	_, err := m.DB.Exec(stmt, bookID)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveAllBooksByEmail returns all books, which are linked to a user with email from arguments
func (m *DatabaseModel) RetrieveAllBooksByEmail(email string) ([]Book, error) {
	stmt := `
			SELECT
			    id, title, description,
			    price, date_of_issue, quote_from,
			    language, category, addcategory, author_id
			FROM book
				WHERE owner_email = $1;
			`

	var books []Book
	rows, err := m.DB.Query(stmt, email)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		book := DbBook{}
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.Price,
			&book.DateOfIssue,
			&book.QuoteFrom,
			&book.Language,
			&book.Category,
			&book.AddCategory,
			&book.AuthorID,
		)

		modelBook := Book{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Price:       book.Price,
			DateOfIssue: book.DateOfIssue,
			QuoteFrom:   "",
			Language:    "",
			Category:    "",
			AddCategory: "",
			AuthorID:    book.AuthorID,
		}

		if book.QuoteFrom.Valid {
			modelBook.QuoteFrom = book.QuoteFrom.String
		}

		if book.Language.Valid {
			modelBook.Language = book.Language.String
		}

		if book.Category.Valid {
			modelBook.Category = book.Category.String
		}

		if book.AddCategory.Valid {
			modelBook.AddCategory = book.AddCategory.String
		}

		if err != nil {
			return nil, err
		}
		books = append(books, modelBook)
	}
	return books, nil

}
