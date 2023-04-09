package models

import (
	"fmt"
	"strings"
	"time"
)

type SesionData struct {
	Email       string    `json:"email"`
	Token       string    `json:"token"`
	TokenExpiry time.Time `json:"tokenExpiry"`
}

// StoreEmailAsSessionToken is used to create new table, which contains data for token and tokenExpiry
// name of database is user email
func (m *DatabaseModel) StoreEmailAsSessionToken(email, token string, tokenExpiry time.Time) error {
	email = strings.ReplaceAll(email, "@", "")
	email = strings.ReplaceAll(email, ".", "")
	stmt := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, email)

	_, err := m.DB.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = fmt.Sprintf(`CREATE TABLE %s(token varchar(180) not null, tokenExpiry date not null);`, email)

	_, err = m.DB.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = fmt.Sprintf(`INSERT 
				INTO %s(token, tokenExpiry)
			VALUES 
			    ($1, $2)
				`, email)

	_, err = m.DB.Exec(stmt, token, tokenExpiry)
	if err != nil {
		return err
	}
	return nil
}

// RetrieveTokenDataFromTable receives data from table Email, and returns token and token Expiry
// if any errors occurs, returns empty string, null time and error
func (m *DatabaseModel) RetrieveTokenDataFromTable(email string) (*SesionData, error) {
	email = strings.ReplaceAll(email, "@", "")
	email = strings.ReplaceAll(email, ".", "")
	stmt := fmt.Sprintf("SELECT * FROM %s", email)
	row, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	sessionData := &SesionData{
		Email: email,
	}

	if row.Next() {
		err = row.Scan(
			&sessionData.Token,
			&sessionData.TokenExpiry,
		)
	}

	if err != nil {
		return nil, err
	}
	return sessionData, nil
}
