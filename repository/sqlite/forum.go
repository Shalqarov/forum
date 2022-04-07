package sqlite

import (
	"database/sql"
)

type Forum struct {
	DB *sql.DB
}
