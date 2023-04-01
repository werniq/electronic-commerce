package models

import "database/sql"

type DatabaseModel struct {
	DB *sql.DB
}
